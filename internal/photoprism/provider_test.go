package photoprism

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/kiosk"
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

func TestProvider_DisplayAsset_PeopleFromMarkers(t *testing.T) {
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
			UID:          "photo-with-faces",
			Type:         "image",
			TakenAt:      takenAt,
			TakenAtLocal: takenAt,
			FileName:     "IMG_002.jpg",
			Width:        1920,
			Height:       1080,
			Portrait:     false,
			Files: []File{{
				Hash:    "def456",
				Primary: true,
				Mime:    "image/jpeg",
				Markers: []Marker{
					{UID: "m1", SubjUID: "subj-alice", Name: "Alice"},
					{UID: "m2", SubjUID: "subj-bob", Name: "Bob"},
				},
			}},
		}}
		_ = json.NewEncoder(w).Encode(photos)
	}))
	defer server.Close()

	cfg := config.Config{PhotoprismURL: server.URL, PhotoprismToken: "token", Kiosk: config.KioskSettings{Cache: false}}
	client := NewClient(server.URL, "token")
	client.HTTPClient = server.Client()
	p := &Provider{client: client, cfg: cfg, ctx: context.Background()}

	err := p.RandomAsset("req1", "dev1", false)
	if err != nil {
		t.Fatal(err)
	}
	got := p.DisplayAsset("req1", "dev1")
	if got.ID != "photo-with-faces" {
		t.Errorf("ID = %q", got.ID)
	}
	if len(got.People) != 2 {
		t.Fatalf("len(People) = %d, want 2", len(got.People))
	}
	names := make(map[string]string)
	for _, p := range got.People {
		names[p.ID] = p.Name
	}
	if names["subj-alice"] != "Alice" || names["subj-bob"] != "Bob" {
		t.Errorf("People = %+v", got.People)
	}
}

func TestBuildPersonSearchQuery(t *testing.T) {
	t.Run("single person", func(t *testing.T) {
		cfg := config.Config{People: []string{"Jane Doe"}}
		got := buildPersonSearchQuery(cfg, "Jane Doe")
		if got != `person:"Jane Doe"` {
			t.Errorf("got %q", got)
		}
	})
	t.Run("require all people", func(t *testing.T) {
		cfg := config.Config{RequireAllPeople: true, People: []string{"Alice", "Bob"}}
		got := buildPersonSearchQuery(cfg, "Alice")
		if got != `people:"Alice & Bob"` {
			t.Errorf("got %q", got)
		}
	})
}

func TestProvider_PersonAssetCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/photos" {
			t.Errorf("path = %s", r.URL.Path)
			return
		}
		q := r.URL.Query().Get("q")
		if q != `person:"Jane"` {
			t.Errorf("q = %q", q)
		}
		w.Header().Set("Content-Type", "application/json")
		photos := []Photo{{UID: "p1", Type: "image"}, {UID: "p2", Type: "image"}}
		_ = json.NewEncoder(w).Encode(photos)
	}))
	defer server.Close()

	cfg := config.Config{PhotoprismURL: server.URL, PhotoprismToken: "t", People: []string{"Jane"}}
	client := NewClient(server.URL, "t")
	client.HTTPClient = server.Client()
	p := &Provider{client: client, cfg: cfg, ctx: context.Background()}

	count, err := p.PersonAssetCount("Jane", "req", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestProvider_RandomAssetOfPerson_setsBucket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/photos" {
			t.Errorf("path = %s", r.URL.Path)
			return
		}
		w.Header().Set("X-Preview-Token", "pt")
		w.Header().Set("X-Download-Token", "dt")
		w.Header().Set("Content-Type", "application/json")
		photos := []Photo{{UID: "p1", Type: "image", Files: []File{{Primary: true, Mime: "image/jpeg"}}}}
		_ = json.NewEncoder(w).Encode(photos)
	}))
	defer server.Close()

	cfg := config.Config{PhotoprismURL: server.URL, PhotoprismToken: "t", Kiosk: config.KioskSettings{Cache: false}}
	client := NewClient(server.URL, "t")
	client.HTTPClient = server.Client()
	p := &Provider{client: client, cfg: cfg, ctx: context.Background()}

	err := p.RandomAssetOfPerson("Jane", "req", "dev", false)
	if err != nil {
		t.Fatal(err)
	}
	got := p.DisplayAsset("req", "dev")
	if got.Bucket != kiosk.SourcePerson {
		t.Errorf("Bucket = %v, want SourcePerson", got.Bucket)
	}
	if got.BucketID != "Jane" {
		t.Errorf("BucketID = %q, want Jane", got.BucketID)
	}
}

func TestProvider_AllNamedPeople_returnsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/subjects" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		subjects := []Subject{{UID: "s1", Name: "Alice"}, {UID: "s2", Name: "Bob"}}
		_ = json.NewEncoder(w).Encode(subjects)
	}))
	defer server.Close()

	cfg := config.Config{PhotoprismURL: server.URL, PhotoprismToken: "t"}
	client := NewClient(server.URL, "t")
	client.HTTPClient = server.Client()
	p := &Provider{client: client, cfg: cfg, ctx: context.Background()}

	people, err := p.AllNamedPeople("req", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if len(people) != 2 {
		t.Fatalf("len(people) = %d, want 2", len(people))
	}
	if people[0].ID != "s1" || people[0].Name != "Alice" {
		t.Errorf("people[0] = %+v", people[0])
	}
	if people[1].ID != "s2" || people[1].Name != "Bob" {
		t.Errorf("people[1] = %+v", people[1])
	}
}

func TestProvider_AllNamedPeople_emptyOnFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.Config{PhotoprismURL: server.URL, PhotoprismToken: "t"}
	client := NewClient(server.URL, "t")
	client.HTTPClient = server.Client()
	p := &Provider{client: client, cfg: cfg, ctx: context.Background()}

	people, err := p.AllNamedPeople("req", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if people != nil {
		t.Errorf("people = %v, want nil on 404", people)
	}
}

func TestProvider_RandomPersonFromAllPeople_errorWhenNoList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.Config{PhotoprismURL: server.URL, PhotoprismToken: "t"}
	client := NewClient(server.URL, "t")
	client.HTTPClient = server.Client()
	p := &Provider{client: client, cfg: cfg, ctx: context.Background()}

	_, err := p.RandomPersonFromAllPeople("req", "dev", true)
	if err == nil {
		t.Fatal("expected error when no subject list")
	}
}

func TestPhotoMarkersToPeople(t *testing.T) {
	t.Run("nil photo", func(t *testing.T) {
		got := photoMarkersToPeople(nil)
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})
	t.Run("no markers", func(t *testing.T) {
		ph := &Photo{Files: []File{{Primary: true}}}
		got := photoMarkersToPeople(ph)
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})
	t.Run("one marker", func(t *testing.T) {
		ph := &Photo{Files: []File{{
			Primary: true,
			Markers: []Marker{{UID: "m1", SubjUID: "subj-1", Name: "Alice"}},
		}}}
		got := photoMarkersToPeople(ph)
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		if got[0].ID != "subj-1" || got[0].Name != "Alice" {
			t.Errorf("got[0] = %+v", got[0])
		}
	})
	t.Run("dedupe by SubjUID", func(t *testing.T) {
		ph := &Photo{Files: []File{
			{Markers: []Marker{{SubjUID: "s1", Name: "Bob"}}},
			{Markers: []Marker{{SubjUID: "s1", Name: "Bob"}}},
		}}
		got := photoMarkersToPeople(ph)
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		if got[0].Name != "Bob" {
			t.Errorf("got[0].Name = %q", got[0].Name)
		}
	})
	t.Run("name fallback when empty", func(t *testing.T) {
		ph := &Photo{Files: []File{{
			Markers: []Marker{{SubjUID: "s2", Name: ""}},
		}}}
		got := photoMarkersToPeople(ph)
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		if got[0].Name != "s2" {
			t.Errorf("got[0].Name = %q, want s2", got[0].Name)
		}
	})
}
