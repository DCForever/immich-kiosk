package photoprism

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/photos" {
			t.Errorf("path = %s, want /api/v1/photos", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer mytoken" {
			t.Errorf("Authorization = %q, want Bearer mytoken", auth)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		w.Header().Set("X-Preview-Token", "preview123")
		w.Header().Set("X-Download-Token", "download456")
		w.Header().Set("Content-Type", "application/json")
		body := []Photo{{UID: "p1", Type: "image", Title: "Test"}}
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer server.Close()

	client := NewClient(server.URL, "mytoken")
	client.HTTPClient = server.Client()

	ctx := context.Background()
	var list []Photo
	headers, err := client.getJSON(ctx, "/api/v1/photos", nil, &list)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].UID != "p1" || list[0].Title != "Test" {
		t.Errorf("list = %v, want one photo UID p1", list)
	}
	if headers["x-preview-token"] != "preview123" || headers["x-download-token"] != "download456" {
		t.Errorf("headers = %v", headers)
	}
}

func TestClient_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad")
	client.HTTPClient = server.Client()

	ctx := context.Background()
	var list []Photo
	_, err := client.getJSON(ctx, "/api/v1/photos", nil, &list)
	if err == nil {
		t.Fatal("expected error for 401")
	}
}
