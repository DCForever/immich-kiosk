// Package immich provides an adapter from *Asset to source.ProviderOps so that
// routes can use the source abstraction while still using the Immich API.

package immich

import (
	"context"
	"strings"

	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/source"
)

// Adapter wraps *Asset and implements source.ProviderOps by delegating to the asset.
type Adapter struct {
	asset *Asset
}

// Ensure Adapter implements source.ProviderOps.
var _ source.ProviderOps = (*Adapter)(nil)

func (*Adapter) _provider() {}

// NewAdapter returns a source.ProviderOps that uses the Immich API via the given config.
func NewAdapter(ctx context.Context, cfg config.Config) source.ProviderOps {
	a := New(ctx, cfg)
	return &Adapter{asset: a}
}

// assetOrder maps string order to Immich AssetOrder.
func assetOrder(order string) AssetOrder {
	switch strings.ToLower(order) {
	case "asc", "ascending", "oldest":
		return Asc
	case "desc", "descending", "newest":
		return Desc
	default:
		return Rand
	}
}

func immichTagFromSource(t source.Tag) Tag {
	return Tag{ID: t.ID, Name: t.Name, Value: t.Value, Color: t.Color}
}

func displayAssetFromImmich(a *Asset, requestID, deviceID string) source.DisplayAsset {
	people := make([]source.Person, len(a.People))
	for i, p := range a.People {
		people[i] = source.Person{
			ID:        p.ID,
			Name:      p.Name,
			BirthDate: source.BirthDate(p.BirthDate),
		}
	}
	tags := make(source.Tags, len(a.Tags))
	for i, t := range a.Tags {
		tags[i] = source.Tag{ID: t.ID, Name: t.Name, Value: t.Value, Color: t.Color}
	}
	albums := make(source.Albums, len(a.AppearsIn))
	for i, al := range a.AppearsIn {
		albums[i] = source.Album{ID: al.ID, AlbumName: al.AlbumName}
	}
	exif := source.ExifInfo{
		City:             a.ExifInfo.City,
		Country:          a.ExifInfo.Country,
		DateTimeOriginal: a.ExifInfo.DateTimeOriginal,
		Description:      a.ExifInfo.Description,
		ExifImageHeight:  a.ExifInfo.ExifImageHeight,
		ExifImageWidth:   a.ExifInfo.ExifImageWidth,
		ExposureTime:     a.ExifInfo.ExposureTime,
		FileSizeInByte:   a.ExifInfo.FileSizeInByte,
		FNumber:          a.ExifInfo.FNumber,
		FocalLength:      a.ExifInfo.FocalLength,
		Iso:              a.ExifInfo.Iso,
		Latitude:         a.ExifInfo.Latitude,
		LensModel:        a.ExifInfo.LensModel,
		Longitude:        a.ExifInfo.Longitude,
		Make:             a.ExifInfo.Make,
		Model:            a.ExifInfo.Model,
		ModifyDate:       a.ExifInfo.ModifyDate,
		Orientation:      a.ExifInfo.Orientation,
		Rating:           a.ExifInfo.Rating,
		State:            a.ExifInfo.State,
		TimeZone:         a.ExifInfo.TimeZone,
	}
	return source.DisplayAsset{
		ID:               a.ID,
		Type:             source.AssetType(a.Type),
		OriginalMimeType: a.OriginalMimeType,
		ServedMimeType:   a.ServedMimeType,
		LocalDateTime:    a.LocalDateTime,
		LivePhotoVideoID: a.LivePhotoVideoID,
		MemoryTitle:      a.MemoryTitle,
		People:           people,
		Tags:             tags,
		AppearsIn:        albums,
		ExifInfo:         exif,
		Owner:            source.Owner{ID: a.Owner.ID, Email: a.Owner.Email, Name: a.Owner.Name},
		IsFavorite:       a.IsFavorite,
		IsPortrait:        a.IsPortrait,
		IsLandscape:      a.IsLandscape,
		Bucket:        a.Bucket,
		BucketID:      a.BucketID,
		SelectedUser:  a.requestConfig.SelectedUser,
		UserOwnsAsset: a.UserOwnsAsset(requestID, deviceID),
	}
}

// DisplayAsset returns the current asset as a source.DisplayAsset.
func (ad *Adapter) DisplayAsset(requestID, deviceID string) source.DisplayAsset {
	return displayAssetFromImmich(ad.asset, requestID, deviceID)
}

func (ad *Adapter) RandomAsset(requestID, deviceID string, isPrefetch bool) error {
	return ad.asset.RandomAsset(requestID, deviceID, isPrefetch)
}

func (ad *Adapter) RandomAssetFromFavourites(requestID, deviceID string, isPrefetch bool) error {
	return ad.asset.RandomAssetFromFavourites(requestID, deviceID, isPrefetch)
}

func (ad *Adapter) AssetFromAlbum(albumID, order string, requestID, deviceID string) error {
	return ad.asset.AssetFromAlbum(albumID, assetOrder(order), requestID, deviceID)
}

func (ad *Adapter) RandomAssetInDateRange(dateRange, requestID, deviceID string, isPrefetch bool) error {
	return ad.asset.RandomAssetInDateRange(dateRange, requestID, deviceID, isPrefetch)
}

func (ad *Adapter) RandomAssetOfPerson(personID, requestID, deviceID string, isPrefetch bool) error {
	return ad.asset.RandomAssetOfPerson(personID, requestID, deviceID, isPrefetch)
}

func (ad *Adapter) RandomMemoryAsset(requestID, deviceID string) error {
	return ad.asset.RandomMemoryAsset(requestID, deviceID)
}

func (ad *Adapter) RandomAssetWithTag(tagID, requestID, deviceID string, isPrefetch bool) error {
	return ad.asset.RandomAssetWithTag(tagID, requestID, deviceID, isPrefetch)
}

func (ad *Adapter) RandomAssetWithRating(ratingID, requestID, deviceID string, isPrefetch bool) error {
	return ad.asset.RandomAssetWithRating(ratingID, requestID, deviceID, isPrefetch)
}

func (ad *Adapter) AssetInfo(assetID, requestID, deviceID string) error {
	ad.asset.ID = assetID
	return ad.asset.AssetInfo(requestID, deviceID)
}

func (ad *Adapter) ApplyUserFromAssetID(assetID string) (string, string) {
	return ad.asset.ApplyUserFromAssetID(assetID)
}

func (ad *Adapter) ApplyDefaultUser() {
	ad.asset.ApplyDefaultUser()
}

func (ad *Adapter) SelectedUser() string {
	return ad.asset.SelectedUser()
}

func (ad *Adapter) PersonAssetCount(personID, requestID, deviceID string) (int, error) {
	return ad.asset.PersonAssetCount(personID, requestID, deviceID)
}

func (ad *Adapter) AlbumImageCount(albumID, requestID, deviceID string) (int, error) {
	return ad.asset.AlbumImageCount(albumID, requestID, deviceID)
}

func (ad *Adapter) AssetsWithTagCount(tagID, requestID, deviceID string) (int, error) {
	return ad.asset.AssetsWithTagCount(tagID, requestID, deviceID)
}

func (ad *Adapter) MemoriesAssetsCount(requestID, deviceID string) int {
	return ad.asset.MemoriesAssetsCount(requestID, deviceID)
}

func (ad *Adapter) AssetsWithRatingCount(rating float32, requestID, deviceID string) (int, error) {
	return ad.asset.AssetsWithRatingCount(rating, requestID, deviceID)
}

func (ad *Adapter) RandomAlbumFromAllAlbums(requestID, deviceID string, excludedAlbums []string) (string, error) {
	return ad.asset.RandomAlbumFromAllAlbums(requestID, deviceID, excludedAlbums)
}

func (ad *Adapter) RandomAlbumFromOwnedAlbums(requestID, deviceID string, excludedAlbums []string) (string, error) {
	return ad.asset.RandomAlbumFromOwnedAlbums(requestID, deviceID, excludedAlbums)
}

func (ad *Adapter) RandomAlbumFromSharedAlbums(requestID, deviceID string, excludedAlbums []string) (string, error) {
	return ad.asset.RandomAlbumFromSharedAlbums(requestID, deviceID, excludedAlbums)
}

func (ad *Adapter) RandomPersonFromAllPeople(requestID, deviceID string, knowPeopleOnly bool) (string, error) {
	return ad.asset.RandomPersonFromAllPeople(requestID, deviceID, knowPeopleOnly)
}

func (ad *Adapter) AllTags(requestID, deviceID string) (source.Tags, string, error) {
	tags, nextPage, err := ad.asset.AllTags(requestID, deviceID)
	if err != nil {
		return nil, "", err
	}
	out := make(source.Tags, len(tags))
	for i, t := range tags {
		out[i] = source.Tag{ID: t.ID, Name: t.Name, Value: t.Value, Color: t.Color}
	}
	return out, nextPage, nil
}

func (ad *Adapter) ExpandTagPatterns(tags []string, requestID, deviceID string) []string {
	return ad.asset.ExpandTagPatterns(tags, requestID, deviceID)
}

func (ad *Adapter) ImagePreview() ([]byte, string, error) {
	return ad.asset.ImagePreview()
}

func (ad *Adapter) Video() ([]byte, string, error) {
	return ad.asset.Video()
}

func (ad *Adapter) AddTag(tag source.Tag) error {
	return ad.asset.AddTag(immichTagFromSource(tag))
}

func (ad *Adapter) RemoveTag(tag source.Tag) error {
	return ad.asset.RemoveTag(immichTagFromSource(tag))
}

func (ad *Adapter) FavouriteStatus(deviceID string, favourite bool) error {
	return ad.asset.FavouriteStatus(deviceID, favourite)
}

func (ad *Adapter) AddToKioskLikedAlbum(requestID, deviceID string) error {
	return ad.asset.AddToKioskLikedAlbum(requestID, deviceID)
}

func (ad *Adapter) RemoveFromKioskLikedAlbum(requestID, deviceID string) error {
	return ad.asset.RemoveFromKioskLikedAlbum(requestID, deviceID)
}

func (ad *Adapter) ArchiveStatus(deviceID string, archive bool) error {
	return ad.asset.ArchiveStatus(deviceID, archive)
}

func (ad *Adapter) RemoveAssetCache(deviceID string) error {
	return ad.asset.RemoveAssetCache(deviceID)
}
