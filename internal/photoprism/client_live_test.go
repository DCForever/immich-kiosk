package photoprism

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"
)

// TestClient_Live_ValidToken runs against a real PhotoPrism instance when
// KIOSK_PHOTOPRISM_URL and KIOSK_PHOTOPRISM_TOKEN are set. It verifies that
// Bearer token authentication works and the API returns a valid photos response.
func TestClient_Live_ValidToken(t *testing.T) {
	baseURL := os.Getenv("KIOSK_PHOTOPRISM_URL")
	token := os.Getenv("KIOSK_PHOTOPRISM_TOKEN")
	if baseURL == "" || token == "" {
		t.Skip("live PhotoPrism test skipped: set KIOSK_PHOTOPRISM_URL and KIOSK_PHOTOPRISM_TOKEN")
	}

	client := NewClient(baseURL, token)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	q := url.Values{}
	q.Set("count", "1")
	q.Set("merged", "true")
	q.Set("primary", "true")

	var list []Photo
	headers, err := client.getJSON(ctx, apiPrefix+"/photos", q, &list)
	if err != nil {
		t.Fatalf("getJSON with valid token: %v", err)
	}

	// Response must be a JSON array (possibly empty)
	if list == nil {
		t.Error("response list is nil, expected non-nil slice")
	}
	// PhotoPrism may return X-Preview-Token and X-Download-Token on photos endpoint
	_ = headers
}

// TestClient_Live_InvalidToken runs against a real PhotoPrism instance when
// KIOSK_PHOTOPRISM_URL is set. It verifies that an invalid token is rejected (401 or error).
func TestClient_Live_InvalidToken(t *testing.T) {
	baseURL := os.Getenv("KIOSK_PHOTOPRISM_URL")
	if baseURL == "" {
		t.Skip("live PhotoPrism test skipped: set KIOSK_PHOTOPRISM_URL")
	}

	client := NewClient(baseURL, "invalid-token")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	q := url.Values{}
	q.Set("count", "1")

	var list []Photo
	_, err := client.getJSON(ctx, apiPrefix+"/photos", q, &list)
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
