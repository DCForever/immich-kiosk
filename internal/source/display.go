// Package source defines the shared media-source abstraction used by the kiosk.
// Types here are backend-agnostic so that both Immich and PhotoPrism can provide assets.
package source

import (
	"errors"
	"time"

	"github.com/damongolding/immich-kiosk/internal/kiosk"
)

// BirthDate is a date string (YYYY-MM-DD) with a Time() method for age calculation.
type BirthDate string

// Time parses the birth date and returns it as time.Time.
func (b BirthDate) Time() (time.Time, error) {
	if string(b) == "" {
		return time.Time{}, errors.New("empty birth date")
	}
	return time.Parse("2006-01-02", string(b))
}

// Person is a person/subject that may appear in an asset.
type Person struct {
	ID        string
	Name      string
	BirthDate BirthDate
}

// Tag is a label/tag on an asset.
type Tag struct {
	ID    string
	Name  string
	Value string
	Color string
}

// Tags is a slice of Tag (used for display and HasTag checks).
type Tags []Tag

// Album is an album that can contain assets.
type Album struct {
	ID        string
	AlbumName string
}

// Albums is a slice of Album (e.g. AppearsIn).
type Albums []Album

// ExifInfo holds EXIF-like metadata for display (location, camera, rating, etc.).
type ExifInfo struct {
	City             string
	Country          string
	DateTimeOriginal time.Time
	Description      string
	ExifImageHeight  int
	ExifImageWidth   int
	ExposureTime     string
	FileSizeInByte   int
	FNumber          float64
	FocalLength      float64
	Iso              int
	Latitude         float64
	LensModel        string
	Longitude        float64
	Make             string
	Model            string
	ModifyDate       time.Time
	Orientation      string
	Rating           float32
	State            string
	TimeZone         string
}

// Owner is the owner of an asset (for display).
type Owner struct {
	ID    string
	Email string
	Name  string
}

// AssetType is the type of media (image, video, etc.).
type AssetType string

const (
	TypeImage AssetType = "IMAGE"
	TypeVideo AssetType = "VIDEO"
)

// DisplayAsset is the source-agnostic representation of an asset for display.
// Routes and templates use this type only. Backends (Immich, PhotoPrism) map their
// API responses into DisplayAsset.
//
// Fields used by templates: ID, Type, People, Tags, AppearsIn, ExifInfo, Owner,
// IsFavorite, IsPortrait, LocalDateTime, MemoryTitle, LivePhotoVideoID,
// OriginalMimeType, ServedMimeType. SelectedUser and UserOwnsAsset are set when
// building the asset for the view. HasTag(tagValue) checks Tags.
type DisplayAsset struct {
	ID               string
	Type             AssetType
	OriginalMimeType string
	ServedMimeType   string
	LocalDateTime    time.Time
	LivePhotoVideoID string
	MemoryTitle      string

	People    []Person
	Tags      Tags
	AppearsIn Albums
	ExifInfo  ExifInfo
	Owner     Owner

	IsFavorite bool
	IsPortrait bool
	IsLandscape bool

	Bucket   kiosk.Source
	BucketID string

	// Set when building view data from config (for template logic).
	SelectedUser  string
	UserOwnsAsset bool
}

// HasTag returns true if the asset has a tag with the given value (e.g. kiosk.TagSkip).
func (a *DisplayAsset) HasTag(tagValue string) bool {
	for _, t := range a.Tags {
		if t.Value == tagValue || t.Name == tagValue {
			return true
		}
	}
	return false
}
