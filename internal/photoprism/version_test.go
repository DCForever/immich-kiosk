package photoprism

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerVersion_FromConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/config" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Errorf("auth = %q", r.Header.Get("Authorization"))
		}
		_, _ = w.Write([]byte(`{"Version":"251010-test"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "tok")
	client.HTTPClient = server.Client()
	if got := client.ServerVersion(context.Background()); got != "251010-test" {
		t.Fatalf("version = %q", got)
	}
	if got := client.ServerVersion(context.Background()); got != "251010-test" {
		t.Fatalf("cached version = %q", got)
	}
}

func TestServerVersion_HeaderAndFailOpen(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("X-Photoprism-Version", "9.9.9")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "tok")
	client.HTTPClient = server.Client()
	if got := client.ServerVersion(context.Background()); got != "9.9.9" {
		t.Fatalf("version = %q", got)
	}
	_ = client.ServerVersion(context.Background())
	if hits != 1 {
		t.Fatalf("hits = %d, want 1", hits)
	}
}

func TestServerVersion_FailOpenCachesMiss(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	first := NewClient(server.URL, "tok")
	first.HTTPClient = server.Client()
	if got := first.ServerVersion(context.Background()); got != "" {
		t.Fatalf("version = %q", got)
	}
	second := NewClient(server.URL, "tok")
	second.HTTPClient = server.Client()
	if got := second.ServerVersion(context.Background()); got != "" {
		t.Fatalf("version = %q", got)
	}
	if hits != 1 {
		t.Fatalf("hits = %d, want 1", hits)
	}
}
