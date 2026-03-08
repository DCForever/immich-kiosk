package photoprism

import "time"

// Photo represents a photo from GET /api/v1/photos (single item or from list).
// See https://docs.photoprism.app/developer-guide/api/search/
type Photo struct {
	ID            string    `json:"ID"`
	UID           string    `json:"UID"`
	Type          string    `json:"Type"` // "image" or "video"
	Hash          string    `json:"Hash"`
	TakenAt       time.Time `json:"TakenAt"`
	TakenAtLocal  time.Time `json:"TakenAtLocal"`
	TimeZone      string    `json:"TimeZone"`
	Title         string    `json:"Title"`
	Description   string    `json:"Description"`
	Favorite      bool      `json:"Favorite"`
	Year          int       `json:"Year"`
	Month         int       `json:"Month"`
	Day           int       `json:"Day"`
	Iso           int       `json:"Iso"`
	FocalLength   int       `json:"FocalLength"`
	FNumber       float64   `json:"FNumber"`
	Exposure      string    `json:"Exposure"`
	Width         int       `json:"Width"`
	Height        int       `json:"Height"`
	Portrait      bool      `json:"Portrait"`
	CameraMake    string    `json:"CameraMake"`
	CameraModel   string    `json:"CameraModel"`
	LensModel     string    `json:"LensModel"`
	Lat           float64   `json:"Lat"`
	Lng           float64   `json:"Lng"`
	PlaceCity     string    `json:"PlaceCity"`
	PlaceState    string    `json:"PlaceState"`
	PlaceCountry  string   `json:"PlaceCountry"`
	PlaceLabel    string   `json:"PlaceLabel"`
	FileUID       string   `json:"FileUID"`
	FileName      string   `json:"FileName"`
	Files         []File   `json:"Files"`
}

// File is a file belonging to a photo (e.g. primary image or video part of live photo).
type File struct {
	UID       string `json:"UID"`
	Hash      string `json:"Hash"`
	Width     int    `json:"Width"`
	Height    int    `json:"Height"`
	Primary   bool   `json:"Primary"`
	FileType  string `json:"FileType"`
	MediaType string `json:"MediaType"`
	Mime      string `json:"Mime"`
	Size      int64  `json:"Size"`
}

// PrimaryHash returns the hash of the primary file for thumbnails/originals, or the first file's hash.
func (p *Photo) PrimaryHash() string {
	for _, f := range p.Files {
		if f.Primary {
			return f.Hash
		}
	}
	if len(p.Files) > 0 {
		return p.Files[0].Hash
	}
	return p.Hash
}

// Album from GET /api/v1/albums.
type Album struct {
	UID         string `json:"UID"`
	Title       string `json:"Title"`
	PhotoCount  int    `json:"PhotoCount"`
}

// Label from PhotoPrism (used for tags).
type Label struct {
	UID  string `json:"UID"`
	Name string `json:"Name"`
	Slug string `json:"Slug"`
}

// Subject from PhotoPrism GET /api/v1/subjects (list). Used for AllNamedPeople when the endpoint exists.
// See https://docs.photoprism.app/developer-guide/api/ (subjects).
type Subject struct {
	UID       string `json:"UID"`
	Name      string `json:"Name"`
	BirthDate string `json:"BirthDate,omitempty"`
}
