package photoprism

import (
	"encoding/json"
	"time"
)

// flexString unmarshals from either a JSON string or number (PhotoPrism list returns string ID, detail returns number id).
type flexString string

func (s *flexString) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && (b[0] == '"' || b[0] == '\'') {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		*s = flexString(str)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	*s = flexString(n.String())
	return nil
}

// Photo represents a photo from GET /api/v1/photos (single item or from list).
// API reference: https://docs.photoprism.dev/
type Photo struct {
	ID            flexString `json:"ID"`            // list may send string, detail sends number
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
// Markers contain face/subject data when present (see https://docs.photoprism.dev/).
type File struct {
	UID       string    `json:"UID"`
	Hash      string    `json:"Hash"`
	Width     int       `json:"Width"`
	Height    int       `json:"Height"`
	Primary   bool      `json:"Primary"`
	FileType  string    `json:"FileType"`
	MediaType string    `json:"MediaType"`
	Mime      string    `json:"Mime"`
	Size      int64     `json:"Size"`
	Markers   []Marker  `json:"Markers,omitempty"`
}

// Marker represents a face/subject on a file (entity.Marker in PhotoPrism API).
type Marker struct {
	UID    string `json:"UID"`
	Name   string `json:"Name"`
	SubjUID string `json:"SubjUID"`
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
// API reference: https://docs.photoprism.dev/
type Subject struct {
	UID       string `json:"UID"`
	Name      string `json:"Name"`
	BirthDate string `json:"BirthDate,omitempty"`
}
