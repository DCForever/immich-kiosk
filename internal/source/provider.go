// Package source defines the media-source provider interface used by routes.
// Both Immich and PhotoPrism implement this interface.

package source

// Provider is a marker interface for the media-source abstraction.
type Provider interface {
	_provider()
}

// ProviderOps is the set of operations routes perform on a provider.
// Implementations (Immich adapter, PhotoPrism) implement these methods.
// The "*Asset" methods set the provider's current asset; callers then use
// DisplayAsset(), ImagePreview(), or Video().
type ProviderOps interface {
	Provider

	// Retrieval: set current asset. Order is "asc"|"desc"|"rand" for AssetFromAlbum.
	RandomAsset(requestID, deviceID string, isPrefetch bool) error
	RandomAssetFromFavourites(requestID, deviceID string, isPrefetch bool) error
	AssetFromAlbum(albumID, order string, requestID, deviceID string) error
	RandomAssetInDateRange(dateRange, requestID, deviceID string, isPrefetch bool) error
	RandomAssetOfPerson(personID, requestID, deviceID string, isPrefetch bool) error
	RandomMemoryAsset(requestID, deviceID string) error
	RandomAssetWithTag(tagID, requestID, deviceID string, isPrefetch bool) error
	RandomAssetWithRating(ratingID, requestID, deviceID string, isPrefetch bool) error
	AssetInfo(assetID, requestID, deviceID string) error

	// Current asset as display type (after a retrieval call).
	// requestID and deviceID are used to compute UserOwnsAsset for the view.
	DisplayAsset(requestID, deviceID string) DisplayAsset

	// User/context helpers.
	ApplyUserFromAssetID(assetID string) (string, string)
	ApplyDefaultUser()
	SelectedUser() string

	// Counts for weighting (gatherAssetBuckets).
	PersonAssetCount(personID, requestID, deviceID string) (int, error)
	AlbumImageCount(albumID, requestID, deviceID string) (int, error)
	AssetsWithTagCount(tagID, requestID, deviceID string) (int, error)
	MemoriesAssetsCount(requestID, deviceID string) int
	AssetsWithRatingCount(rating float32, requestID, deviceID string) (int, error)

	// Pickers for weighted random (retrieveImage).
	RandomAlbumFromAllAlbums(requestID, deviceID string, excludedAlbums []string) (string, error)
	RandomAlbumFromOwnedAlbums(requestID, deviceID string, excludedAlbums []string) (string, error)
	RandomAlbumFromSharedAlbums(requestID, deviceID string, excludedAlbums []string) (string, error)
	RandomPersonFromAllPeople(requestID, deviceID string, knowPeopleOnly bool) (string, error)

	// Tags.
	AllTags(requestID, deviceID string) (Tags, string, error)
	ExpandTagPatterns(tags []string, requestID, deviceID string) []string

	// Image/video bytes for the current asset.
	ImagePreview() ([]byte, string, error)
	Video() ([]byte, string, error)

	// Mutations (like/hide/tag); may no-op on backends that do not support them.
	AddTag(tag Tag) error
	RemoveTag(tag Tag) error
	FavouriteStatus(deviceID string, favourite bool) error
	AddToKioskLikedAlbum(requestID, deviceID string) error
	RemoveFromKioskLikedAlbum(requestID, deviceID string) error
	ArchiveStatus(deviceID string, archive bool) error
	RemoveAssetCache(deviceID string) error
}
