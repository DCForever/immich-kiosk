package routes

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/damongolding/immich-kiosk/internal/cache"
	"github.com/damongolding/immich-kiosk/internal/common"
	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/immich"
	"github.com/damongolding/immich-kiosk/internal/photoprism"
	"github.com/damongolding/immich-kiosk/internal/kiosk"
	"github.com/damongolding/immich-kiosk/internal/source"
	imageComponent "github.com/damongolding/immich-kiosk/internal/templates/components/image"
	videoComponent "github.com/damongolding/immich-kiosk/internal/templates/components/video"
	"github.com/damongolding/immich-kiosk/internal/utils"
	"github.com/damongolding/immich-kiosk/internal/webhooks"
	"github.com/fogleman/gg"
	"github.com/labstack/echo/v5"
)

var errVideoNotReady = errors.New("video not ready")

// getProvider returns a media source provider for the given config.
// When config.Source is "photoprism" returns a PhotoPrism provider; otherwise returns the Immich adapter.
func getProvider(ctx context.Context, cfg config.Config) source.ProviderOps {
	if cfg.Source == config.SourcePhotoPrism {
		return photoprism.NewProvider(ctx, cfg)
	}
	return immich.NewAdapter(ctx, cfg)
}

// gatherAssetBuckets collects asset weightings for people, albums and date ranges.
// For each person, it gets the count of images containing that person.
// For each album, it gets the total count of images in the album.
// For date ranges, it currently assigns a fixed weighting of 1000.
// These weightings are used to determine the probability of selecting images from each source.
//
// Parameters:
//   - provider: The media source provider used to query image counts
//   - requestConfig: Configuration containing people, albums and dates to gather assets for
//   - requestID: Identifier for the current request for logging
//
// Returns:
//   - A slice of AssetWithWeighting containing the weightings for each asset source
//   - An error if any database queries fail
func gatherAssetBuckets(provider source.ProviderOps, requestConfig config.Config, requestID, deviceID string) ([]utils.AssetWithWeighting, error) {

	assets := []utils.AssetWithWeighting{}

	// People bucket
	for _, person := range requestConfig.People {
		if person == "" || strings.EqualFold(person, "none") {
			continue
		}

		personTmp, _ := provider.ApplyUserFromAssetID(person)

		personAssetCount, personCountErr := provider.PersonAssetCount(personTmp, requestID, deviceID)
		if personCountErr != nil {
			if provider.SelectedUser() != "" {
				return nil, fmt.Errorf("user '<b>%s</b>' has no Person '%s'. error='%w'", provider.SelectedUser(), personTmp, personCountErr)
			}
			return nil, fmt.Errorf("getting person image count: %w", personCountErr)
		}

		if personAssetCount == 0 {
			log.Error("No assets found for", "person", personTmp)
			continue
		}

		assets = append(assets, utils.AssetWithWeighting{
			Asset:  utils.WeightedAsset{Type: kiosk.SourcePerson, ID: person},
			Weight: personAssetCount,
		})
	}

	// Albums bucket
	for _, album := range requestConfig.Albums {
		if album == "" || strings.EqualFold(album, "none") {
			continue
		}

		albumTmp, _ := provider.ApplyUserFromAssetID(album)

		albumAssetCount, albumCountErr := provider.AlbumImageCount(albumTmp, requestID, deviceID)
		if albumCountErr != nil {
			if provider.SelectedUser() != "" {
				return nil, fmt.Errorf("user '<b>%s</b>' has no Album '%s'. error='%w'", provider.SelectedUser(), albumTmp, albumCountErr)
			}
			return nil, fmt.Errorf("getting album asset count: %w", albumCountErr)
		}

		if albumAssetCount == 0 {
			log.Error("No assets found for", "album", albumTmp)
			continue
		}

		assets = append(assets, utils.AssetWithWeighting{
			Asset:  utils.WeightedAsset{Type: kiosk.SourceAlbum, ID: album},
			Weight: albumAssetCount,
		})
	}

	// Use the default user for the rest of the request (tags, dates, memories)
	provider.ApplyDefaultUser()

	// Tags bucket
	requestConfig.Tags = provider.ExpandTagPatterns(requestConfig.Tags, requestID, deviceID)

	for _, tag := range requestConfig.Tags {
		if tag == "" || strings.EqualFold(tag, "none") {
			continue
		}

		if strings.Contains(tag, "@") {
			log.Warn("Tags with multi user information are not currently supported")
			tag, _, _ = strings.Cut(tag, "@")
		}

		tags, _, tagsErr := provider.AllTags(requestID, deviceID)
		if tagsErr != nil {
			log.Error("getting tags", "err", tagsErr)
			continue
		}

		tagData, tagErr := tags.Get(tag)
		if tagErr != nil {
			log.Error("getting tag from tags", "tag", tag, "err", tagErr)
			continue
		}

		taggedAssetsCount, tagCountErr := provider.AssetsWithTagCount(tagData.ID, requestID, deviceID)
		if tagCountErr != nil {
			if requestConfig.SelectedUser != "" {
				return nil, fmt.Errorf("user '<b>%s</b>' has no assets with tag '%s'. error='%w'", requestConfig.SelectedUser, tagData.Value, tagCountErr)
			}
			return nil, fmt.Errorf("getting tagged asset count: %w", tagCountErr)
		}

		if taggedAssetsCount == 0 {
			log.Error("No assets found with", "tag", tagData.Value)
			continue
		}

		assets = append(assets, utils.AssetWithWeighting{
			Asset:  utils.WeightedAsset{Type: kiosk.SourceTag, ID: tagData.ID},
			Weight: taggedAssetsCount,
		})
	}

	// Dates bucket
	for _, date := range requestConfig.Dates {
		if date == "" || strings.EqualFold(date, "none") {
			continue
		}

		if strings.Contains(date, "@") {
			log.Warn("Dates with multi user information are not currently supported")
			date, _, _ = strings.Cut(date, "@")
		}

		// use FetchedAssetsSize as a weighting for date ranges
		assets = append(assets, utils.AssetWithWeighting{
			Asset:  utils.WeightedAsset{Type: kiosk.SourceDateRange, ID: date},
			Weight: requestConfig.Kiosk.FetchedAssetsSize,
		})
	}

	// Rating bucket
	if requestConfig.Rating > -1 {
		ratedErr := gatherRatedAssets(provider, requestConfig, requestID, deviceID, &assets)
		if ratedErr != nil {
			log.Error(ratedErr)
		}
	}

	// Memories bucket
	if requestConfig.Memories {
		memories := provider.MemoriesAssetsCount(requestID, deviceID)
		if memories == 0 {
			log.Warn("No assets found for memories")
		} else {
			assets = append(assets, utils.AssetWithWeighting{
				Asset:   utils.WeightedAsset{Type: kiosk.SourceMemories, ID: "memories"},
				Weight:  memories,
				Penalty: requestConfig.MemoryWeight,
			})
		}
	}

	return assets, nil
}

func gatherRatedAssets(provider source.ProviderOps, requestConfig config.Config, requestID, deviceID string, assets *[]utils.AssetWithWeighting) error {
	wantedRating := requestConfig.Rating

	ratedAssetsCount, ratedCountErr := provider.AssetsWithRatingCount(wantedRating, requestID, deviceID)
	if ratedCountErr != nil {
		if requestConfig.SelectedUser != "" {
			return fmt.Errorf("user '<b>%s</b>' has no assets with rating '%f'. error='%w'", requestConfig.SelectedUser, wantedRating, ratedCountErr)
		}
		return fmt.Errorf("getting rated asset count: %w", ratedCountErr)
	}

	if ratedAssetsCount > 0 {
		*assets = append(*assets, utils.AssetWithWeighting{
			Asset:  utils.WeightedAsset{Type: kiosk.SourceRating, ID: fmt.Sprintf("rating-%.2f", requestConfig.Rating)},
			Weight: ratedAssetsCount,
		})
	} else {
		log.Error("No assets found with", "rating", wantedRating)
	}

	return nil
}

// isSleepMode checks if the kiosk should currently be in sleep mode based on configured sleep times
func isSleepMode(requestConfig config.Config) bool {
	if requestConfig.SleepStart == "" || requestConfig.SleepEnd == "" {
		return false
	}

	if isSleepTime, _ := utils.IsSleepTime(requestConfig.SleepStart, requestConfig.SleepEnd, time.Now()); isSleepTime {
		return isSleepTime
	}

	return false
}

// retrieveImage fetches a random image based on the picked image type.
// It returns an error if the image retrieval fails.
func retrieveImage(provider source.ProviderOps, pickedAsset utils.WeightedAsset, albumOrder string, excludedAlbums []string, requestID, deviceID string, isPrefetch bool) error {

	switch pickedAsset.Type {
	case kiosk.SourceAlbum:
		switch pickedAsset.ID {
		case kiosk.AlbumKeywordAll:
			pickedAlbumID, err := provider.RandomAlbumFromAllAlbums(requestID, deviceID, excludedAlbums)
			if err != nil {
				return err
			}
			pickedAsset.ID = pickedAlbumID
		case kiosk.AlbumKeywordOwned:
			pickedAlbumID, err := provider.RandomAlbumFromOwnedAlbums(requestID, deviceID, excludedAlbums)
			if err != nil {
				return err
			}
			pickedAsset.ID = pickedAlbumID
		case kiosk.AlbumKeywordShared:
			pickedAlbumID, err := provider.RandomAlbumFromSharedAlbums(requestID, deviceID, excludedAlbums)
			if err != nil {
				return err
			}
			pickedAsset.ID = pickedAlbumID
		case kiosk.AlbumKeywordFavourites, kiosk.AlbumKeywordFavorites:
			return provider.RandomAssetFromFavourites(requestID, deviceID, isPrefetch)
		}

		switch strings.ToLower(albumOrder) {
		case config.AlbumOrderDescending, config.AlbumOrderDesc, config.AlbumOrderNewest:
			return provider.AssetFromAlbum(pickedAsset.ID, "desc", requestID, deviceID)
		case config.AlbumOrderAscending, config.AlbumOrderAsc, config.AlbumOrderOldest:
			return provider.AssetFromAlbum(pickedAsset.ID, "asc", requestID, deviceID)
		default:
			return provider.AssetFromAlbum(pickedAsset.ID, "rand", requestID, deviceID)
		}

	case kiosk.SourceDateRange:
		return provider.RandomAssetInDateRange(pickedAsset.ID, requestID, deviceID, isPrefetch)

	case kiosk.SourcePerson:
		if pickedAsset.ID == kiosk.PersonKeywordAll {
			pickedPersonID, err := provider.RandomPersonFromAllPeople(requestID, deviceID, true)
			if err != nil {
				return err
			}
			pickedAsset.ID = pickedPersonID
		}

		return provider.RandomAssetOfPerson(pickedAsset.ID, requestID, deviceID, isPrefetch)

	case kiosk.SourceMemories:
		return provider.RandomMemoryAsset(requestID, deviceID)

	case kiosk.SourceTag:
		return provider.RandomAssetWithTag(pickedAsset.ID, requestID, deviceID, isPrefetch)

	case kiosk.SourceRating:
		return provider.RandomAssetWithRating(pickedAsset.ID, requestID, deviceID, isPrefetch)

	case kiosk.SourceRandom:
		fallthrough

	default:
		return provider.RandomAsset(requestID, deviceID, isPrefetch)
	}

}

// fetchImagePreview retrieves and decodes an image preview from the current provider asset.
// Returns the processed image or an error if retrieval or decoding fails.
func fetchImagePreview(provider source.ProviderOps, isOriginal bool, requestID, deviceID string, isPrefetch bool) (image.Image, error) {
	imageGet := time.Now()

	imgBytes, _, err := provider.ImagePreview()
	if err != nil {
		return nil, fmt.Errorf("getting image preview: %w", err)
	}

	img, err := utils.BytesToImage(imgBytes, isOriginal)
	if err != nil {
		return nil, err
	}

	if isPrefetch {
		log.Debug(requestID, "PREFETCH", deviceID, "Got image in", time.Since(imageGet).Seconds())
	} else {
		log.Debug(requestID, "Got image in", time.Since(imageGet).Seconds())
	}

	return img, nil
}

// processAsset handles the entire process of selecting and retrieving an image.
// It returns the image bytes and an error if any step fails.
func processAsset(provider source.ProviderOps, requestConfig config.Config, requestID string, deviceID string, requestURL string, isPrefetch bool) (image.Image, error) {

	var err error

	assets, assetsErr := gatherAssetBuckets(provider, requestConfig, requestID, deviceID)
	if assetsErr != nil {
		return nil, assetsErr
	}

	for range maxProcessAssetRetries {

		pickedAsset := utils.PickRandomImageType(requestConfig.Kiosk.AssetWeighting, assets)

		pickedAsset.ID, _ = provider.ApplyUserFromAssetID(pickedAsset.ID)

		err = retrieveImage(provider, pickedAsset, requestConfig.AlbumOrder, requestConfig.ExcludedAlbums, requestID, deviceID, isPrefetch)
		if err != nil {
			continue
		}

		displayAsset := provider.DisplayAsset(requestID, deviceID)
		if requestConfig.ShowVideos && displayAsset.Type == source.TypeVideo {
			var img image.Image
			img, err = processVideo(provider, requestConfig, requestID, deviceID, requestURL, isPrefetch)
			if err == nil {
				return img, nil
			}
			if errors.Is(err, errVideoNotReady) {
				log.Debug(requestID+" Video not ready, trying another asset", "video", displayAsset.ID)
				continue
			}
			return nil, err
		}

		return processImage(provider, requestConfig, requestID, deviceID, isPrefetch)
	}

	return nil, fmt.Errorf("%w: max retries exceeded", err)
}

// processVideo handles retrieving and processing video assets.
// For Immich it uses VideoManager to download; for other providers it fetches preview from the provider.
func processVideo(provider source.ProviderOps, requestConfig config.Config, requestID string, deviceID string, requestURL string, isPrefetch bool) (image.Image, error) {
	displayAsset := provider.DisplayAsset(requestID, deviceID)
	if ad, ok := provider.(*immich.Adapter); ok {
		immichAsset := ad.Asset()
		if VideoManager.IsDownloaded(immichAsset.ID) {
			return fetchImagePreview(provider, requestConfig.UseOriginalImage, requestID, deviceID, isPrefetch)
		}
		if !VideoManager.IsDownloading(immichAsset.ID) {
			go VideoManager.DownloadVideo(*immichAsset, displayAsset, requestConfig, deviceID, requestURL)
		}
		return nil, errVideoNotReady
	}
	return fetchImagePreview(provider, requestConfig.UseOriginalImage, requestID, deviceID, isPrefetch)
}

// processImage prepares an image asset for display and retrieves a preview.
// For Immich with LivePhotos, triggers background download of the video.
func processImage(provider source.ProviderOps, requestConfig config.Config, requestID string, deviceID string, isPrefetch bool) (image.Image, error) {
	displayAsset := provider.DisplayAsset(requestID, deviceID)
	if requestConfig.LivePhotos && displayAsset.LivePhotoVideoID != "" {
		if ad, ok := provider.(*immich.Adapter); ok {
			immichAsset := ad.Asset()
			isDownloaded := VideoManager.IsDownloaded(immichAsset.LivePhotoVideoID)
			isDownloading := VideoManager.IsDownloading(immichAsset.LivePhotoVideoID)
			if !isDownloaded && !isDownloading {
				livePhoto := immich.New(context.TODO(), requestConfig)
				livePhoto.ID = displayAsset.LivePhotoVideoID
				if err := livePhoto.AssetInfo(requestID, deviceID); err != nil {
					return nil, err
				}
				livePhotoDisplay := source.DisplayAsset{ID: displayAsset.LivePhotoVideoID}
				go VideoManager.DownloadVideo(livePhoto, livePhotoDisplay, requestConfig, deviceID, "")
			}
		}
	}
	return fetchImagePreview(provider, requestConfig.UseOriginalImage, requestID, deviceID, isPrefetch)
}

// imageToBase64 converts image bytes to a base64 string and logs the processing time.
// It returns the base64 string and an error if conversion fails.
func imageToBase64(img image.Image, config config.Config, requestID, deviceID string, action string, isPrefetch bool) (string, error) {
	startTime := time.Now()

	imgBytes, err := utils.ImageToBase64(img)
	if err != nil {
		return "", fmt.Errorf("converting image to base64: %w", err)
	}

	logImageProcessing(config, requestID, deviceID, isPrefetch, action, startTime)
	return imgBytes, nil
}

// shouldSkipBlur determines whether background blur should be skipped.
// - Blur is skipped when BackgroundBlur is disabled.
// - Blur is skipped when ImageFit is "cover" and LivePhotos is disabled.
func shouldSkipBlur(config config.Config) bool {
	if !config.BackgroundBlur {
		return true
	}

	usingImageCover := strings.EqualFold(config.ImageFit, "cover")

	// Skip if using image cover with live photos off
	if usingImageCover && !config.LivePhotos {
		return true
	}

	return false
}

// processBlurredImage applies a blur effect to the image if required by the configuration.
// It returns the blurred image as a base64 string and an error if any occurs.
func processBlurredImage(img image.Image, assetType immich.AssetType, config config.Config, requestID, deviceID string, isPrefetch bool) (string, error) {
	isImage := assetType == immich.ImageType
	skipBlur := shouldSkipBlur(config)

	if isImage && skipBlur {
		return "", nil
	}

	startTime := time.Now()
	imgBlur, err := utils.BlurImage(img, config.BackgroundBlurAmount, config.OptimizeImages, config.ClientData.Width, config.ClientData.Height)
	if err != nil {
		return "", fmt.Errorf("blurring image: %w", err)
	}

	logImageProcessing(config, requestID, deviceID, isPrefetch, "Blurred", startTime)

	return imageToBase64(imgBlur, config, requestID, deviceID, "Converted blurred", isPrefetch)
}

// logImageProcessing logs the time taken for image processing if debug verbose is enabled.
func logImageProcessing(config config.Config, requestID, deviceID string, isPrefetch bool, action string, startTime time.Time) {
	if !config.Kiosk.DebugVerbose {
		return
	}

	duration := time.Since(startTime).Seconds()
	if isPrefetch {
		log.Debug(requestID, "PREFETCH", deviceID, action+" image in", duration)
	} else {
		log.Debug(requestID, action+" image in", duration)
	}
}

// DrawFaceOnImage draws bounding boxes around detected faces in an image
func DrawFaceOnImage(img image.Image, i *immich.Asset) image.Image {

	if len(i.People) == 0 && len(i.UnassignedFaces) == 0 {
		log.Debug("no people found")
		return img
	}

	dc := gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy())

	dc.DrawImage(img, 0, 0)

	for _, person := range i.People {
		for _, face := range person.Faces {
			width := face.BoundingBoxX2 - face.BoundingBoxX1
			height := face.BoundingBoxY2 - face.BoundingBoxY1

			dc.DrawRectangle(float64(face.BoundingBoxX1), float64(face.BoundingBoxY1), float64(width), float64(height))
			dc.SetHexColor("#990000")
			dc.Fill()
		}
	}

	for _, face := range i.UnassignedFaces {
		width := face.BoundingBoxX2 - face.BoundingBoxX1
		height := face.BoundingBoxY2 - face.BoundingBoxY1

		dc.DrawRectangle(float64(face.BoundingBoxX1), float64(face.BoundingBoxY1), float64(width), float64(height))
		dc.SetHexColor("#000099")
		dc.Fill()
	}

	facesBoundX, facesBoundY := i.FacesCenterPointPX()
	dc.DrawRectangle(facesBoundX-10, facesBoundY-10, 20, 20)
	dc.SetHexColor("#889900")
	dc.Fill()

	return dc.Image()

}

// processViewImageData processes an image request and returns view data for display.
// It handles the complete workflow from selecting an image to preparing it for display,
// including face detection, optimization, and format conversion.
//
// Parameters:
//   - requestConfig: Configuration settings for the request
//   - c: Copy of the request context
//   - isPrefetch: Whether this is a prefetch request
//   - options: Additional options for image processing
//
// Returns:
//   - ViewImageData containing the processed image and metadata
//   - Error if any step fails
func processViewImageData(requestConfig config.Config, c common.ContextCopy, isPrefetch bool, options common.ViewImageDataOptions) (common.ViewImageData, error) {
	// Initialize request metadata
	metadata := requestMetadata{
		requestID: utils.ColorizeRequestID(c.ResponseHeader.Get(echo.HeaderXRequestID)),
		deviceID:  c.RequestHeader.Get("kiosk-device-id"),
		urlString: c.URL.String(),
	}

	setupRequestConfig(&requestConfig)
	provider := getProvider(context.Background(), requestConfig)
	if options.ImageOrientation == "PORTRAIT" || options.ImageOrientation == "LANDSCAPE" {
		provider.SetRatioWanted(options.ImageOrientation)
	}

	if options.RelativeAssetWanted {
		handleRelativeAssetConfig(&requestConfig, options)
	}

	img, err := processAsset(provider, requestConfig, metadata.requestID, metadata.deviceID, metadata.urlString, isPrefetch)
	if err != nil {
		return common.ViewImageData{}, fmt.Errorf("selecting asset: %w", err)
	}

	img = handleFaceProcessing(img, provider, requestConfig, metadata)

	if requestConfig.OptimizeImages {
		img, err = utils.OptimizeImage(img, requestConfig.ClientData.Width, requestConfig.ClientData.Height)
		if err != nil {
			return common.ViewImageData{}, err
		}
	}

	displayAsset := provider.DisplayAsset(metadata.requestID, metadata.deviceID)
	imgString, imgBlurString, dominantColor, err := convertImages(img, immich.AssetType(displayAsset.Type), requestConfig, metadata, isPrefetch)
	if err != nil {
		return common.ViewImageData{}, err
	}

	return common.ViewImageData{
		Asset:               displayAsset,
		ImageData:           imgString,
		ImageBlurData:       imgBlurString,
		ImageDominantColor:  dominantColor,
		User:                provider.SelectedUser(),
	}, nil
}

// setupRequestConfig configures the selected user for the request by picking a random
// user from the config if multiple users are provided, otherwise sets to empty string
func setupRequestConfig(config *config.Config) {
	if len(config.User) > 0 {
		randomIndex := rand.IntN(len(config.User))
		config.SelectedUser = config.User[randomIndex]
	} else {
		config.SelectedUser = ""
	}
}

// handleRelativeAssetConfig updates the config buckets based on the relative asset options.
// Resets existing buckets and configures the appropriate bucket based on the asset source type.
func handleRelativeAssetConfig(config *config.Config, options common.ViewImageDataOptions) {
	config.ResetBuckets()
	config.Memories = false

	switch options.RelativeAssetBucket {
	case kiosk.SourceAlbum:
		config.Albums = append(config.Albums, options.RelativeAssetBucketID)
	case kiosk.SourcePerson:
		config.People = append(config.People, options.RelativeAssetBucketID)
	case kiosk.SourceDateRange:
		config.Dates = append(config.Dates, options.RelativeAssetBucketID)
	case kiosk.SourceTag:
		config.Tags = append(config.Tags, options.RelativeAssetBucketID)
	case kiosk.SourceMemories:
		config.Memories = true
	case kiosk.SourceRandom:
	}
}

// handleFaceProcessing processes face detection and drawing for an image.
// For Immich adapter, checks for faces if smart-zoom is enabled and draws faces if configured.
// Returns the processed image.
func handleFaceProcessing(img image.Image, provider source.ProviderOps, config config.Config, metadata requestMetadata) image.Image {
	if ad, ok := provider.(*immich.Adapter); ok {
		asset := ad.Asset()
		if strings.EqualFold(config.ImageEffect, "smart-zoom") && len(asset.People)+len(asset.UnassignedFaces) == 0 {
			asset.CheckForFaces(metadata.requestID, metadata.deviceID)
		}
		if ShouldDrawFacesOnImages() {
			log.Debug("Drawing faces")
			return DrawFaceOnImage(img, asset)
		}
	}
	return img
}

// convertImages converts the provided image to base64 strings for both normal and blurred versions.
// Returns the base64 encoded normal image, blurred image, and any error that occurred.
func convertImages(img image.Image, assetType immich.AssetType, config config.Config, metadata requestMetadata, isPrefetch bool) (string, string, color.RGBA, error) {

	var dominantColor color.RGBA

	imgString, err := imageToBase64(img, config, metadata.requestID, metadata.deviceID, "Converted", isPrefetch)
	if err != nil {
		return "", "", dominantColor, err
	}

	imgBlurString, err := processBlurredImage(img, assetType, config, metadata.requestID, metadata.deviceID, isPrefetch)
	if err != nil {
		return "", "", dominantColor, err
	}

	if config.Theme == kiosk.ThemeBubble || config.UseOfflineMode {
		dominantColor, err = utils.ExtractDominantColor(img)
		if err != nil {
			return "", "", dominantColor, err
		}
	}

	return imgString, imgBlurString, dominantColor, nil
}

// ProcessViewImageData processes view data for an image without orientation constraints
func ProcessViewImageData(requestConfig config.Config, c common.ContextCopy, isPrefetch bool) (common.ViewImageData, error) {
	return processViewImageData(requestConfig, c, isPrefetch, common.ViewImageDataOptions{})
}

func ProcessViewImageDataWithOptions(requestConfig config.Config, c common.ContextCopy, isPrefetch bool, options common.ViewImageDataOptions) (common.ViewImageData, error) {
	return processViewImageData(requestConfig, c, isPrefetch, options)
}

// assetToCache stores view data in the cache and triggers prefetch webhooks
func assetToCache(ctx context.Context, viewDataToAdd common.ViewData, requestConfig *config.Config, deviceID string, requestData *common.RouteRequestData, c common.ContextCopy) {

	cache.AssetToCache(viewDataToAdd, requestConfig, deviceID, c.URL.String())

	go webhooks.Trigger(ctx, requestData, KioskVersion, webhooks.PrefetchAsset, viewDataToAdd)
}

// assetPreFetch handles prefetching assets for the current request
func assetPreFetch(common *common.Common, requestData *common.RouteRequestData, c common.ContextCopy) {

	requestConfig := requestData.RequestConfig
	requestID := requestData.RequestID
	deviceID := requestData.DeviceID

	viewDataToAdd, err := generateViewData(requestConfig, c, requestID, deviceID, true)
	if err != nil {
		log.Error("generateViewData", "prefetch", true, "err", err)
		return
	}

	assetToCache(common.Context(), viewDataToAdd, &requestConfig, deviceID, requestData, c)
}

// fromCache retrieves cached page data for a given request and device ID.
func fromCache(urlString string, deviceID string) []common.ViewData {
	cacheKey := cache.ViewCacheKey(urlString, deviceID)
	if data, found := cache.Get(cacheKey); found {
		cachedPageData, ok := data.([]common.ViewData)
		if !ok {
			log.Error("cache: invalid data type", "type", fmt.Sprintf("%T", data), "key", cacheKey)
			cache.Delete(cacheKey)
			return nil
		}

		if len(cachedPageData) > 0 {
			return cachedPageData
		}

		cache.Delete(cacheKey)
	}
	return nil
}

// renderCachedViewData renders cached page data and updates the cache.
func renderCachedViewData(c *echo.Context, cachedViewData []common.ViewData, requestConfig *config.Config, requestID, deviceID, secret string) error {

	log.Debug(requestID, "deviceID", deviceID, "cache hit for new image", true)

	cacheKey := cache.ViewCacheKey(c.Request().URL.String(), deviceID)

	viewDataToRender := cachedViewData[0]
	cache.Set(cacheKey, cachedViewData[1:], requestConfig.Duration)

	// Update history which will be outdated in cache
	utils.TrimHistory(&requestConfig.History, kiosk.HistoryLimit)
	viewDataToRender.History = requestConfig.History

	if requestConfig.ShowVideos && viewDataToRender.Assets[0].Asset.Type == source.TypeVideo {
		return Render(c, http.StatusOK, videoComponent.Video(viewDataToRender, secret))
	}

	return Render(c, http.StatusOK, imageComponent.Image(viewDataToRender, secret))
}

// fetchSecondSplitViewAsset attempts to retrieve a second asset for split view layouts that is different from the first asset.
// It tries up to three times to obtain a unique asset and appends it to the provided ViewData if successful.
// Returns an error if asset retrieval fails.
func fetchSecondSplitViewAsset(viewData *common.ViewData, viewDataSplitView common.ViewImageData, requestConfig config.Config, c common.ContextCopy, isPrefetch bool, options common.ViewImageDataOptions) error {
	const maxImageRetrievalAttempts = 3

	for range maxImageRetrievalAttempts {
		viewDataSplitViewSecond, err := ProcessViewImageDataWithOptions(requestConfig, c, isPrefetch, options)
		if err != nil {
			return err
		}

		if viewDataSplitView.Asset.ID != viewDataSplitViewSecond.Asset.ID {
			viewData.Assets = append(viewData.Assets, viewDataSplitViewSecond)
			return nil
		}
	}
	return nil
}

// determineLayoutMode returns the appropriate layout mode based on the requested layout and client display dimensions.
// If the requested layout is split view and the client height exceeds the width, it switches to landscape split view mode.
func determineLayoutMode(layout string, clientHeight, clientWidth int) string {
	if layout == kiosk.LayoutSplitview && clientHeight > clientWidth {
		return kiosk.LayoutSplitviewLandscape
	}
	return layout
}

// generateViewData prepares view data for a kiosk page request based on the specified layout and client display dimensions.
// It selects and processes one or two assets as needed for the layout, handling orientation and split view logic, and returns the resulting ViewData or an error.
func generateViewData(requestConfig config.Config, c common.ContextCopy, requestID, deviceID string, isPrefetch bool) (common.ViewData, error) {

	viewData := common.ViewData{
		RequestID: requestID,
		DeviceID:  deviceID,
		Config:    requestConfig,
	}

	requestConfig.Layout = determineLayoutMode(requestConfig.Layout, requestConfig.ClientData.Height, requestConfig.ClientData.Width)

	switch requestConfig.Layout {
	case kiosk.LayoutLandscape, kiosk.LayoutPortrait:
		options := common.ViewImageDataOptions{
			ImageOrientation: "LANDSCAPE",
		}
		if requestConfig.Layout == kiosk.LayoutPortrait {
			options.ImageOrientation = "PORTRAIT"
		}
		viewDataSingle, err := ProcessViewImageDataWithOptions(requestConfig, c, isPrefetch, options)
		if err != nil {
			return viewData, err
		}
		viewData.Assets = append(viewData.Assets, viewDataSingle)

	case kiosk.LayoutSplitview:
		viewDataSplitView, err := ProcessViewImageData(requestConfig, c, isPrefetch)
		if err != nil {
			return viewData, err
		}
		viewData.Assets = append(viewData.Assets, viewDataSplitView)

		if viewDataSplitView.Asset.Type == source.TypeVideo || viewDataSplitView.Asset.IsLandscape {
			return viewData, nil
		}

		options := common.ViewImageDataOptions{
			RelativeAssetWanted:   true,
			RelativeAssetBucket:   viewDataSplitView.Asset.Bucket,
			RelativeAssetBucketID: viewDataSplitView.Asset.BucketID,
			ImageOrientation:      "PORTRAIT",
		}

		// Second image
		if secondAssetErr := fetchSecondSplitViewAsset(&viewData, viewDataSplitView, requestConfig, c, isPrefetch, options); secondAssetErr != nil {
			return viewData, secondAssetErr
		}

	case kiosk.LayoutSplitviewLandscape:
		viewDataSplitView, err := ProcessViewImageData(requestConfig, c, isPrefetch)
		if err != nil {
			return viewData, err
		}
		viewData.Assets = append(viewData.Assets, viewDataSplitView)

		if viewDataSplitView.Asset.IsPortrait {
			return viewData, nil
		}

		options := common.ViewImageDataOptions{
			RelativeAssetWanted:   true,
			RelativeAssetBucket:   viewDataSplitView.Asset.Bucket,
			RelativeAssetBucketID: viewDataSplitView.Asset.BucketID,
			ImageOrientation:      "LANDSCAPE",
		}

		// Second image
		if secondAssetErr := fetchSecondSplitViewAsset(&viewData, viewDataSplitView, requestConfig, c, isPrefetch, options); secondAssetErr != nil {
			return viewData, secondAssetErr
		}

	default:
		viewDataSingle, err := ProcessViewImageData(requestConfig, c, isPrefetch)
		if err != nil {
			return viewData, err
		}
		viewData.Assets = append(viewData.Assets, viewDataSingle)
	}

	return viewData, nil
}
