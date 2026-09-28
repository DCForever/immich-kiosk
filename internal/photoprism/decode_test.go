package photoprism

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPhoto_DefensiveJSON(t *testing.T) {
	raw := []byte(`{
		"id": 42,
		"uid": "photo-1",
		"type": "image",
		"favorite": 1,
		"takenAt": "",
		"takenAtLocal": "2024-06-15 14:30:00",
		"width": "1920",
		"height": 1080,
		"portrait": 0,
		"fNumber": "2.8",
		"files": [{
			"hash": "abc",
			"primary": "true",
			"mime": "image/jpeg",
			"markers": [{"name": "Ada", "subjectUid": "subj-1"}]
		}]
	}`)
	var photo Photo
	if err := json.Unmarshal(raw, &photo); err != nil {
		t.Fatal(err)
	}
	if photo.ID != "42" || photo.UID != "photo-1" || !photo.Favorite {
		t.Fatalf("identity = %+v", photo)
	}
	if !photo.TakenAt.IsZero() {
		t.Fatalf("TakenAt = %v, want zero", photo.TakenAt)
	}
	wantLocal := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	if !photo.TakenAtLocal.Equal(wantLocal) {
		t.Fatalf("TakenAtLocal = %v", photo.TakenAtLocal)
	}
	if photo.Width != 1920 || photo.Height != 1080 || photo.Portrait || photo.FNumber != 2.8 {
		t.Fatalf("capture = %+v", photo)
	}
	if photo.PrimaryHash() != "abc" {
		t.Fatalf("hash = %s", photo.PrimaryHash())
	}
	if len(photo.Files[0].Markers) != 1 || photo.Files[0].Markers[0].SubjUID != "subj-1" {
		t.Fatalf("markers = %+v", photo.Files[0].Markers)
	}
}

func TestPhoto_RoundTripStandardJSON(t *testing.T) {
	taken := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	in := Photo{
		UID: "u", Favorite: true, TakenAt: taken, Width: 10,
		Files: []File{{Hash: "h", Primary: true, Markers: []Marker{{Name: "A", SubjUID: "s"}}}},
	}
	body, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out Photo
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatal(err)
	}
	if out.UID != "u" || !out.Favorite || out.Width != 10 || out.PrimaryHash() != "h" {
		t.Fatalf("out = %+v", out)
	}
	if !out.TakenAt.Equal(taken) || len(out.Files[0].Markers) != 1 || out.Files[0].Markers[0].Name != "A" {
		t.Fatalf("out = %+v", out)
	}
}

func TestDecodeBody_WrappedListAndEmpty(t *testing.T) {
	var wrapped []Photo
	body := []byte(`{"photos":[{"UID":"p1","Title":"T"}],"previewToken":"pt"}`)
	if err := decodeBody(body, &wrapped); err != nil {
		t.Fatal(err)
	}
	if len(wrapped) != 1 || wrapped[0].UID != "p1" || wrapped[0].Title != "T" {
		t.Fatalf("wrapped = %+v", wrapped)
	}

	var albums []Album
	if err := decodeBody([]byte(`[{"uid":"a1","title":"Trip","photoCount":"12"}]`), &albums); err != nil {
		t.Fatal(err)
	}
	if len(albums) != 1 || albums[0].UID != "a1" || albums[0].PhotoCount != 12 {
		t.Fatalf("albums = %+v", albums)
	}

	var empty []Photo
	if err := decodeBody([]byte("null"), &empty); err != nil {
		t.Fatal(err)
	}
	if empty == nil || len(empty) != 0 {
		t.Fatalf("empty = %#v", empty)
	}
}

func TestClient_WrappedPhotosAndNoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/empty" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("X-Photoprism-Preview-Token", "from-header")
		_, _ = w.Write([]byte(`{"photos":[{"UID":"p9"}],"downloadToken":"dl"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	client.HTTPClient = server.Client()
	ctx := context.Background()

	var list []Photo
	headers, err := client.getJSON(ctx, "/api/v1/photos", nil, &list)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].UID != "p9" {
		t.Fatalf("list = %+v", list)
	}
	if headers["x-preview-token"] != "from-header" || headers["x-download-token"] != "dl" {
		t.Fatalf("headers = %v", headers)
	}

	var none []Photo
	if _, err := client.getJSON(ctx, "/api/v1/empty", nil, &none); err != nil {
		t.Fatal(err)
	}
	if len(none) != 0 {
		t.Fatalf("none = %+v", none)
	}
}

func TestCompactPhotos_DropsBlank(t *testing.T) {
	list := compactPhotos([]Photo{{}, {UID: "keep"}, {Hash: "h"}})
	if len(list) != 2 {
		t.Fatalf("len = %d", len(list))
	}
}
