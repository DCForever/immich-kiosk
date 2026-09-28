package photoprism

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImagePreview_ThumbSizeFallback(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if r.URL.Path != "/api/v1/t/abc/tok/fit_1280" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte("jpeg"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	client.HTTPClient = server.Client()
	p := &Provider{
		client:       client,
		ctx:          context.Background(),
		current:      []Photo{{Hash: "abc", Type: "image"}},
		previewToken: "tok",
	}
	body, contentType, err := p.ImagePreview()
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "jpeg" || contentType != "image/jpeg" {
		t.Fatalf("body %q content-type %q", body, contentType)
	}
	if len(paths) != 2 || paths[0] != "/api/v1/t/abc/tok/fit_720" || paths[1] != "/api/v1/t/abc/tok/fit_1280" {
		t.Fatalf("paths = %v", paths)
	}
}
