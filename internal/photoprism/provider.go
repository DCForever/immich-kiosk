// Package photoprism provider implements source.ProviderOps for PhotoPrism.
// People (subjects) are supported via search filters person:"Name" and people:"A & B";
// AllNamedPeople and RandomPersonFromAllPeople use GET /api/v1/subjects when available.
// Memories and mutations (AddTag, FavouriteStatus, etc.) are no-ops or best-effort where the API supports them.
package photoprism

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/damongolding/immich-kiosk/internal/cache"
	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/kiosk"
	"github.com/damongolding/immich-kiosk/internal/source"
)

const cacheBatchSize = 30

// buildPersonSearchQuery returns PhotoPrism search q value for person/people filter.
// For RequireAllPeople with multiple config people uses people:"A & B"; else person:"name".
// Names with double quotes are escaped for the search query.
func buildPersonSearchQuery(cfg config.Config, singlePersonID string) string {
	if cfg.RequireAllPeople && len(cfg.People) > 1 {
		escaped := make([]string, 0, len(cfg.People))
		for _, name := range cfg.People {
			escaped = append(escaped, strings.ReplaceAll(name, `"`, `\"`))
		}
		return `people:"` + strings.Join(escaped, " & ") + `"`
	}
	return `person:"` + strings.ReplaceAll(singlePersonID, `"`, `\"`) + `"`
}

// cachedPhotoBatch is stored in cache (JSON) to reuse API responses.
type cachedPhotoBatch struct {
	List          []Photo `json:"list"`
	PreviewToken  string  `json:"preview_token"`
	DownloadToken string  `json:"download_token"`
}

// Provider implements source.ProviderOps for PhotoPrism.
type Provider struct {
	client         *Client
	cfg            config.Config
	ctx            context.Context
	current        []Photo
	previewToken   string
	downloadToken  string
	ratioWanted    string
	currentBucket  kiosk.Source
	currentBucketID string
}

// Ensure Provider implements source.ProviderOps.
var _ source.ProviderOps = (*Provider)(nil)

func (*Provider) Provider() {}

// NewProvider returns a PhotoPrism-backed source.ProviderOps.
func NewProvider(ctx context.Context, cfg config.Config) source.ProviderOps {
	return &Provider{
		client: NewClient(cfg.PhotoprismURL, cfg.PhotoprismToken),
		cfg:    cfg,
		ctx:    ctx,
	}
}

func (p *Provider) photosQuery(count int, order string, albumUID string, searchQ string) url.Values {
	q := url.Values{}
	q.Set("count", strconv.Itoa(count))
	q.Set("order", order)
	q.Set("merged", "true")
	q.Set("primary", "true")
	if albumUID != "" {
		q.Set("s", albumUID)
	}
	if searchQ != "" {
		q.Set("q", searchQ)
	}
	return q
}

func (p *Provider) fetchPhotos(albumUID string, order string, count int, searchQ string, label string) ([]Photo, string, string, error) {
	q := p.photosQuery(count, order, albumUID, searchQ)
	if label != "" {
		q.Set("label", label)
	}
	var list []Photo
	headers, err := p.client.getJSON(p.ctx, apiPrefix+"/photos", q, &list)
	if err != nil {
		return nil, "", "", err
	}
	preview := headers["x-preview-token"]
	download := headers["x-download-token"]
	if p.cfg.Kiosk.Debug && len(list) > 0 {
		log.Debug("PhotoPrism: fetched photos", "count", len(list), "preview_token_set", preview != "", "download_token_set", download != "")
		for i, ph := range list {
			if i > 0 && !p.cfg.Kiosk.DebugVerbose {
				log.Debug("PhotoPrism: ... and more photos", "total", len(list))
				break
			}
			p.logPhotoRetrieved(&ph, i+1, len(list))
		}
	}
	return list, preview, download, nil
}

// fetchPhotosWithCache returns one photo by either popping from cache or fetching a batch from the API.
// When cache is enabled, it stores batches (cachedPhotoBatch) keyed by apiURL+deviceID+date and reuses them.
// label is optional (for tag/label filter); dateFilter is optional (e.g. "after:2020-01-01"); personFilter is optional (e.g. person:"Name" or people:"A & B").
func (p *Provider) fetchPhotosWithCache(albumUID, order, label, dateFilter, personFilter, requestID, deviceID string) error {
	q := p.photosQuery(cacheBatchSize, order, albumUID, "")
	if label != "" {
		q.Set("label", label)
	}
	var searchParts []string
	if personFilter != "" {
		searchParts = append(searchParts, personFilter)
	}
	if dateFilter != "" {
		searchParts = append(searchParts, dateFilter)
	}
	if len(searchParts) > 0 {
		q.Set("q", strings.Join(searchParts, " "))
	}
	apiURL := p.client.BaseURL + apiPrefix + "/photos?" + q.Encode()
	cacheKey := cache.APICacheKey(apiURL, deviceID, "")

	if p.cfg.Kiosk.Cache {
		if data, found := cache.Get(cacheKey); found {
			if b, ok := data.([]byte); ok {
				var batch cachedPhotoBatch
				if err := json.Unmarshal(b, &batch); err != nil {
					cache.Delete(cacheKey)
				} else if len(batch.List) > 0 {
					first := batch.List[0]
					p.current = []Photo{first}
					p.previewToken = batch.PreviewToken
					p.downloadToken = batch.DownloadToken
					if p.cfg.Kiosk.Debug {
						log.Debug("PhotoPrism: using photo from cache", "uid", first.UID, "remaining", len(batch.List)-1)
						p.logPhotoRetrieved(&first, 1, len(batch.List))
					}
					if len(batch.List) > 1 {
						rest := cachedPhotoBatch{
							List:          batch.List[1:],
							PreviewToken:  batch.PreviewToken,
							DownloadToken: batch.DownloadToken,
						}
						jsonBytes, _ := json.Marshal(rest)
						cache.Set(cacheKey, jsonBytes, p.cfg.Duration)
					} else {
						cache.Delete(cacheKey)
					}
					p.enrichCurrentPhotoWithMarkers()
					return nil
				} else {
					cache.Delete(cacheKey)
				}
			}
		}
	}

	searchQ := strings.Join(searchParts, " ")
	list, preview, download, err := p.fetchPhotos(albumUID, order, cacheBatchSize, searchQ, label)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("photoprism: no photos")
	}
	p.current = []Photo{list[0]}
	p.previewToken = preview
	p.downloadToken = download
	if p.cfg.Kiosk.Cache && len(list) > 1 {
		rest := cachedPhotoBatch{
			List:          list[1:],
			PreviewToken:  preview,
			DownloadToken: download,
		}
		jsonBytes, _ := json.Marshal(rest)
		cache.Set(cacheKey, jsonBytes, p.cfg.Duration)
	}
	p.enrichCurrentPhotoWithMarkers()
	return nil
}

// enrichCurrentPhotoWithMarkers fetches full photo details (GET /api/v1/photos/:uid) when
// ShowPersonName is true, so we get Files with Markers (faces/subjects) that the list endpoint often omits.
func (p *Provider) enrichCurrentPhotoWithMarkers() {
	if !p.cfg.ShowPersonName || len(p.current) == 0 {
		return
	}
	uid := p.current[0].UID
	detailed, err := p.client.GetPhoto(p.ctx, uid)
	if err != nil {
		if p.cfg.Kiosk.Debug {
			log.Debug("PhotoPrism: GetPhoto failed (markers may be missing)", "uid", uid, "err", err)
		}
		return
	}
	p.current[0] = *detailed
	if p.cfg.Kiosk.Debug {
		markers := 0
		for _, f := range detailed.Files {
			markers += len(f.Markers)
		}
		log.Debug("PhotoPrism: enriched photo with details", "uid", uid, "markers", markers)
	}
}

// currentPhoto returns the first current photo or nil.
func (p *Provider) currentPhoto() *Photo {
	if len(p.current) == 0 {
		return nil
	}
	return &p.current[0]
}

// bucketForDisplay returns the bucket to show in DisplayAsset (defaults to SourceAlbum when unset).
func (p *Provider) bucketForDisplay() kiosk.Source {
	if p.currentBucket != "" {
		return p.currentBucket
	}
	return kiosk.SourceAlbum
}

// logPhotoRetrieved logs metadata and markers for a photo when debug is enabled.
func (p *Provider) logPhotoRetrieved(ph *Photo, index, total int) {
	if ph == nil || !p.cfg.Kiosk.Debug {
		return
	}
	args := []any{"index", index, "total", total, "uid", ph.UID, "type", ph.Type, "taken_at", ph.TakenAt, "file_name", ph.FileName}
	if ph.Title != "" {
		args = append(args, "title", ph.Title)
	}
	if ph.Description != "" {
		args = append(args, "description", ph.Description)
	}
	args = append(args, "favorite", ph.Favorite, "files", len(ph.Files))
	markerCount := 0
	var markerNames []string
	for _, f := range ph.Files {
		for _, m := range f.Markers {
			markerCount++
			if m.Name != "" {
				markerNames = append(markerNames, m.Name)
			} else if m.SubjUID != "" {
				markerNames = append(markerNames, m.SubjUID)
			}
		}
	}
	args = append(args, "markers", markerCount)
	if markerCount > 0 {
		args = append(args, "people", markerNames)
	}
	if p.cfg.Kiosk.DebugVerbose {
		args = append(args, "width", ph.Width, "height", ph.Height, "camera_make", ph.CameraMake, "camera_model", ph.CameraModel)
		if len(ph.Files) > 0 {
			for fi, f := range ph.Files {
				args = append(args, fmt.Sprintf("file_%d_uid", fi), f.UID, fmt.Sprintf("file_%d_markers", fi), len(f.Markers))
				for mi, m := range f.Markers {
					args = append(args, fmt.Sprintf("file_%d_marker_%d", fi, mi), fmt.Sprintf("%s (SubjUID:%s)", m.Name, m.SubjUID))
				}
			}
		}
	}
	log.Debug("PhotoPrism: photo retrieved", args...)
}

// photoMarkersToPeople collects unique people from a photo's file markers (face/subject data).
// Deduplicates by SubjUID so each person appears once; uses Name for display.
func photoMarkersToPeople(ph *Photo) []source.Person {
	if ph == nil {
		return nil
	}
	seen := make(map[string]string) // SubjUID -> Name
	for _, f := range ph.Files {
		for _, m := range f.Markers {
			uid := strings.TrimSpace(m.SubjUID)
			name := strings.TrimSpace(m.Name)
			if uid == "" && name == "" {
				continue
			}
			if uid == "" {
				uid = name
			}
			if _, ok := seen[uid]; !ok {
				seen[uid] = name
			}
		}
	}
	if len(seen) == 0 {
		return nil
	}
	out := make([]source.Person, 0, len(seen))
	for uid, name := range seen {
		if name == "" {
			name = uid
		}
		out = append(out, source.Person{ID: uid, Name: name})
	}
	return out
}

func (p *Provider) RandomAsset(requestID, deviceID string, isPrefetch bool) error {
	p.currentBucket = kiosk.SourceAlbum
	p.currentBucketID = ""
	return p.fetchPhotosWithCache("", "random", "", "", "", requestID, deviceID)
}

func (p *Provider) RandomAssetFromFavourites(requestID, deviceID string, isPrefetch bool) error {
	q := p.photosQuery(1, "random", "", "")
	q.Set("favorite", "true")
	var list []Photo
	headers, err := p.client.getJSON(p.ctx, apiPrefix+"/photos", q, &list)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("photoprism: no favourite photos")
	}
	p.current = list
	p.previewToken = headers["x-preview-token"]
	p.downloadToken = headers["x-download-token"]
	p.currentBucket = kiosk.SourceAlbum
	p.currentBucketID = ""
	return nil
}

func (p *Provider) AssetFromAlbum(albumID, order string, requestID, deviceID string) error {
	ord := "newest"
	switch strings.ToLower(order) {
	case "asc", "ascending", "oldest":
		ord = "oldest"
	case "desc", "descending", "newest":
		ord = "newest"
	default:
		ord = "random"
	}
	p.currentBucket = kiosk.SourceAlbum
	p.currentBucketID = albumID
	if err := p.fetchPhotosWithCache(albumID, ord, "", "", "", requestID, deviceID); err != nil {
		return err
	}
	if len(p.current) == 0 {
		return fmt.Errorf("photoprism: no photos in album %s", albumID)
	}
	return nil
}

func (p *Provider) RandomAssetInDateRange(dateRange, requestID, deviceID string, isPrefetch bool) error {
	dateFilter := ""
	if dateRange != "" {
		dateFilter = "after:" + dateRange
	}
	p.currentBucket = kiosk.SourceDateRange
	p.currentBucketID = dateRange
	return p.fetchPhotosWithCache("", "random", "", dateFilter, "", requestID, deviceID)
}

// RandomAssetOfPerson fetches a random photo containing the given person (subject name or face UID).
// PhotoPrism search: person:"Name" or people:"A & B" for RequireAllPeople; face:<uid> for face UID.
func (p *Provider) RandomAssetOfPerson(personID, requestID, deviceID string, isPrefetch bool) error {
	p.currentBucket = kiosk.SourcePerson
	p.currentBucketID = personID
	personFilter := buildPersonSearchQuery(p.cfg, personID)
	if err := p.fetchPhotosWithCache("", "random", "", "", personFilter, requestID, deviceID); err != nil {
		return err
	}
	if len(p.current) == 0 {
		return fmt.Errorf("photoprism: no photos for person %s", personID)
	}
	return nil
}

// RandomMemoryAsset is not supported by PhotoPrism.
func (p *Provider) RandomMemoryAsset(requestID, deviceID string) error {
	return fmt.Errorf("photoprism: memories not supported")
}

func (p *Provider) RandomAssetWithTag(tagID, requestID, deviceID string, isPrefetch bool) error {
	p.currentBucket = kiosk.SourceTag
	p.currentBucketID = tagID
	if err := p.fetchPhotosWithCache("", "random", tagID, "", "", requestID, deviceID); err != nil {
		return err
	}
	if len(p.current) == 0 {
		return fmt.Errorf("photoprism: no photos with label %s", tagID)
	}
	return nil
}

// RandomAssetWithRating is not supported by PhotoPrism (no star rating in API).
func (p *Provider) RandomAssetWithRating(ratingID, requestID, deviceID string, isPrefetch bool) error {
	return fmt.Errorf("photoprism: rating filter not supported")
}

func (p *Provider) AssetInfo(assetID, requestID, deviceID string) error {
	q := url.Values{}
	q.Set("count", "1")
	q.Set("merged", "true")
	q.Set("primary", "true")
	q.Set("q", "uid:"+assetID)
	var list []Photo
	headers, err := p.client.getJSON(p.ctx, apiPrefix+"/photos", q, &list)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("photoprism: photo not found: %s", assetID)
	}
	p.current = list
	p.previewToken = headers["x-preview-token"]
	p.downloadToken = headers["x-download-token"]
	p.currentBucket = kiosk.SourceAlbum
	p.currentBucketID = ""
	return nil
}

func (p *Provider) DisplayAsset(requestID, deviceID string) source.DisplayAsset {
	ph := p.currentPhoto()
	if ph == nil {
		return source.DisplayAsset{}
	}
	taken := ph.TakenAtLocal
	if taken.IsZero() {
		taken = ph.TakenAt
	}
	assetType := source.TypeImage
	if strings.EqualFold(ph.Type, "video") {
		assetType = source.TypeVideo
	}
	mime := "image/jpeg"
	if len(ph.Files) > 0 {
		for _, f := range ph.Files {
			if f.Primary && f.Mime != "" {
				mime = f.Mime
				break
			}
		}
	}
	exif := source.ExifInfo{
		Description:      ph.Description,
		DateTimeOriginal:  ph.TakenAt,
		Make:              ph.CameraMake,
		Model:             ph.CameraModel,
		LensModel:         ph.LensModel,
		Iso:               ph.Iso,
		FocalLength:       float64(ph.FocalLength),
		FNumber:           ph.FNumber,
		ExposureTime:      ph.Exposure,
		ExifImageWidth:    ph.Width,
		ExifImageHeight:   ph.Height,
		Latitude:          ph.Lat,
		Longitude:         ph.Lng,
		City:              ph.PlaceCity,
		State:             ph.PlaceState,
		Country:           ph.PlaceCountry,
		TimeZone:          ph.TimeZone,
	}
	people := photoMarkersToPeople(ph)
	if p.cfg.Kiosk.Debug && ph != nil {
		names := make([]string, 0, len(people))
		for _, pe := range people {
			names = append(names, pe.Name)
		}
		log.Debug("PhotoPrism: DisplayAsset", "uid", ph.UID, "type", assetType, "people", names, "file_name", ph.FileName)
	}
	return source.DisplayAsset{
		ID:               ph.UID,
		Type:             assetType,
		OriginalMimeType: mime,
		ServedMimeType:   mime,
		LocalDateTime:    taken,
		MemoryTitle:      "",
		OriginalFileName: ph.FileName,
		People:           people,
		Tags:             nil,
		AppearsIn:        nil,
		ExifInfo:         exif,
		IsFavorite:       ph.Favorite,
		IsPortrait:       ph.Portrait,
		IsLandscape:      !ph.Portrait && ph.Width >= ph.Height,
		Bucket:           p.bucketForDisplay(),
		BucketID:         p.currentBucketID,
		SelectedUser:     p.cfg.SelectedUser,
		UserOwnsAsset:    true,
	}
}

func (p *Provider) ApplyUserFromAssetID(assetID string) (string, string) {
	return assetID, ""
}
func (p *Provider) ApplyDefaultUser() {}
func (p *Provider) SelectedUser() string {
	return p.cfg.SelectedUser
}
func (p *Provider) SetRatioWanted(orientation string) {
	p.ratioWanted = orientation
}

func (p *Provider) PersonAssetCount(personID, requestID, deviceID string) (int, error) {
	if personID == kiosk.PersonKeywordAll {
		// Sum over all named people would require a subject list; without it return 0.
		return 0, nil
	}
	personFilter := buildPersonSearchQuery(p.cfg, personID)
	q := p.photosQuery(1000, "newest", "", personFilter)
	var list []Photo
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/photos", q, &list)
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (p *Provider) AlbumImageCount(albumID, requestID, deviceID string) (int, error) {
	// GET /api/v1/albums?count=1&uid=albumID or list and find
	q := url.Values{}
	q.Set("count", "500")
	var list []Album
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/albums", q, &list)
	if err != nil {
		return 0, err
	}
	for _, a := range list {
		if a.UID == albumID || a.Title == albumID {
			return a.PhotoCount, nil
		}
	}
	return 0, nil
}

func (p *Provider) AssetsWithTagCount(tagID, requestID, deviceID string) (int, error) {
	q := url.Values{}
	q.Set("count", "1000")
	q.Set("label", tagID)
	var list []Photo
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/photos", q, &list)
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (p *Provider) MemoriesAssetsCount(requestID, deviceID string) int {
	return 0
}

func (p *Provider) AssetsWithRatingCount(rating float32, requestID, deviceID string) (int, error) {
	return 0, nil
}

func (p *Provider) RandomAlbumFromAllAlbums(requestID, deviceID string, excludedAlbums []string) (string, error) {
	q := url.Values{}
	q.Set("count", "500")
	var list []Album
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/albums", q, &list)
	if err != nil {
		return "", err
	}
	excl := make(map[string]struct{})
	for _, id := range excludedAlbums {
		excl[id] = struct{}{}
	}
	var candidates []Album
	for _, a := range list {
		if a.PhotoCount <= 0 {
			continue
		}
		if _, ok := excl[a.UID]; ok {
			continue
		}
		if _, ok := excl[a.Title]; ok {
			continue
		}
		candidates = append(candidates, a)
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("photoprism: no albums")
	}
	a := candidates[rand.IntN(len(candidates))]
	return a.UID, nil
}

func (p *Provider) RandomAlbumFromOwnedAlbums(requestID, deviceID string, excludedAlbums []string) (string, error) {
	return p.RandomAlbumFromAllAlbums(requestID, deviceID, excludedAlbums)
}
func (p *Provider) RandomAlbumFromSharedAlbums(requestID, deviceID string, excludedAlbums []string) (string, error) {
	return p.RandomAlbumFromAllAlbums(requestID, deviceID, excludedAlbums)
}

func (p *Provider) RandomPersonFromAllPeople(requestID, deviceID string, knowPeopleOnly bool) (string, error) {
	people, err := p.AllNamedPeople(requestID, deviceID)
	if err != nil || len(people) == 0 {
		return "", fmt.Errorf("photoprism: no subject list available or no named people")
	}
	excl := make(map[string]struct{})
	for _, id := range p.cfg.ExcludedPeople {
		excl[id] = struct{}{}
	}
	var candidates []source.Person
	for _, pe := range people {
		if _, ok := excl[pe.ID]; ok {
			continue
		}
		if _, ok := excl[pe.Name]; ok {
			continue
		}
		candidates = append(candidates, pe)
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("photoprism: no people after exclusions")
	}
	picked := candidates[rand.IntN(len(candidates))]
	return picked.ID, nil
}

func (p *Provider) AllTags(requestID, deviceID string) (source.Tags, string, error) {
	// PhotoPrism labels: GET /api/v1/labels
	q := url.Values{}
	q.Set("count", "500")
	var list []struct {
		UID  string `json:"UID"`
		Name string `json:"Name"`
		Slug string `json:"Slug"`
	}
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/labels", q, &list)
	if err != nil {
		return nil, "", err
	}
	out := make(source.Tags, 0, len(list))
	for _, l := range list {
		slug := l.Slug
		if slug == "" {
			slug = strings.ToLower(strings.ReplaceAll(l.Name, " ", "-"))
		}
		out = append(out, source.Tag{ID: l.UID, Name: l.Name, Value: slug})
	}
	return out, "", nil
}

func (p *Provider) ExpandTagPatterns(tags []string, requestID, deviceID string) []string {
	return tags
}

func (p *Provider) AllNamedPeople(requestID, deviceID string) ([]source.Person, error) {
	// Try GET /api/v1/subjects; if the endpoint does not exist or fails, return empty (URL builder people dropdown stays empty).
	q := url.Values{}
	q.Set("count", "500")
	var list []Subject
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/subjects", q, &list)
	if err != nil {
		return nil, nil
	}
	out := make([]source.Person, 0, len(list))
	for _, s := range list {
		name := strings.TrimSpace(s.Name)
		if name == "" {
			continue
		}
		out = append(out, source.Person{
			ID:        s.UID,
			Name:      name,
			BirthDate: source.BirthDate(s.BirthDate),
		})
	}
	return out, nil
}

func (p *Provider) AllAlbums(requestID, deviceID string) (source.Albums, error) {
	q := url.Values{}
	q.Set("count", "500")
	var list []Album
	_, err := p.client.getJSON(p.ctx, apiPrefix+"/albums", q, &list)
	if err != nil {
		return nil, err
	}
	out := make(source.Albums, len(list))
	for i, a := range list {
		out[i] = source.Album{ID: a.UID, AlbumName: a.Title}
	}
	return out, nil
}

func (p *Provider) ImagePreview() ([]byte, string, error) {
	ph := p.currentPhoto()
	if ph == nil {
		return nil, "", fmt.Errorf("photoprism: no current photo")
	}
	token := p.previewToken
	if token == "" {
		token = "public"
	}
	path := fmt.Sprintf("%s/t/%s/%s/fit_720", apiPrefix, ph.PrimaryHash(), token)
	return p.client.getBytes(p.ctx, path)
}

func (p *Provider) Video() ([]byte, string, error) {
	ph := p.currentPhoto()
	if ph == nil || !strings.EqualFold(ph.Type, "video") {
		return nil, "", fmt.Errorf("photoprism: no current video")
	}
	token := p.previewToken
	if token == "" {
		token = "public"
	}
	path := fmt.Sprintf("%s/videos/%s/%s/avc", apiPrefix, ph.PrimaryHash(), token)
	return p.client.getBytes(p.ctx, path)
}

func (p *Provider) AddTag(tag source.Tag) error       { return nil }
func (p *Provider) RemoveTag(tag source.Tag) error   { return nil }
func (p *Provider) FavouriteStatus(deviceID string, favourite bool) error { return nil }
func (p *Provider) AddToKioskLikedAlbum(requestID, deviceID string) error { return nil }
func (p *Provider) RemoveFromKioskLikedAlbum(requestID, deviceID string) error { return nil }
func (p *Provider) ArchiveStatus(deviceID string, archive bool) error { return nil }
func (p *Provider) RemoveAssetCache(deviceID string) error { return nil }
