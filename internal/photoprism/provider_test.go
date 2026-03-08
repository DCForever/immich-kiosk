package photoprism

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/source"
)

func TestProvider_DisplayAsset_Mapping(t *testing.T) {
	// Fixed time for deterministic test
	takenAt := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/photos" {
			t.Errorf("path = %s", r.URL.Path)
			return
		}
		w.Header().Set("X-Preview-Token", "preview-tok")
		w.Header().Set("X-Download-Token", "dl-tok")
		w.Header().Set("Content-Type", "application/json")
		photos := []Photo{{
			UID:          "photo-uid-123",
			Type:         "image",
			Title:        "Test Photo",
			Description:  "A test",
			Favorite:     true,
			TakenAt:      takenAt,
			TakenAtLocal: takenAt,
			FileName:     "IMG_001.jpg",
			CameraMake:   "Canon",
			CameraModel:  "EOS R5",
			Iso:          400,
			FocalLength:  50,
			FNumber:      2.8,
			Exposure:     "1/125",
			Width:        1920,
			Height:       1080,
			Portrait:     false,
			Lat:          52.52,
			Lng:          13.405,
			PlaceCity:    "Berlin",
			PlaceCountry: "Germany",
			Files:        []File{{Hash: "abc123", Primary: true, Mime: "image/jpeg"}},
		}}
		_ = json.NewEncoder(w).Encode(photos)
	}))
	defer server.Close()

	cfg := config.Config{}
	cfg.PhotoprismURL = server.URL
	cfg.PhotoprismToken = "token"
	cfg.Kiosk.Cache = false
	cfg.Duration = 60

	client := NewClient(server.URL, "token")
	client.HTTPClient = server.Client()

	ctx := context.Background()
	p := &Provider{
		client: client,
		cfg:    cfg,
		ctx:    ctx,
	}

	err := p.RandomAsset("req1", "dev1", false)
	if err != nil {
		t.Fatal(err)
	}

	got := p.DisplayAsset("req1", "dev1")

	if got.ID != "photo-uid-123" {
		t.Errorf("ID = %q, want photo-uid-123", got.ID)
	}
	if got.Type != source.TypeImage {
		t.Errorf("Type = %v, want TypeImage", got.Type)
	}
	if got.OriginalFileName != "IMG_001.jpg" {
		t.Errorf("OriginalFileName = %q", got.OriginalFileName)
	}
	if !got.IsFavorite {
		t.Error("IsFavorite = false, want true")
	}
	if got.ExifInfo.Make != "Canon" || got.ExifInfo.Model != "EOS R5" {
		t.Errorf("ExifInfo Make/Model = %q / %q", got.ExifInfo.Make, got.ExifInfo.Model)
	}
	if got.ExifInfo.Iso != 400 || got.ExifInfo.FocalLength != 50 {
		t.Errorf("ExifInfo Iso/FocalLength = %d / %f", got.ExifInfo.Iso, got.ExifInfo.FocalLength)
	}
	if got.ExifInfo.City != "Berlin" || got.ExifInfo.Country != "Germany" {
		t.Errorf("ExifInfo City/Country = %q / %q", got.ExifInfo.City, got.ExifInfo.Country)
	}
	if got.ExifInfo.Latitude != 52.52 || got.ExifInfo.Longitude != 13.405 {
		t.Errorf("ExifInfo Lat/Lng = %f / %f", got.ExifInfo.Latitude, got.ExifInfo.Longitude)
	}
}
