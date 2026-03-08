package photoprism

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// livePhotoPrismCreds returns baseURL and token for live tests: from env first,
// then from config.yaml (photoprism_url, photoprism_token) if present.
func livePhotoPrismCreds() (baseURL, token string) {
	baseURL = strings.TrimSpace(os.Getenv("KIOSK_PHOTOPRISM_URL"))
	token = strings.TrimSpace(os.Getenv("KIOSK_PHOTOPRISM_TOKEN"))
	if baseURL != "" && token != "" {
		return baseURL, token
	}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("../..")
	if err := v.ReadInConfig(); err != nil {
		return baseURL, token
	}
	if baseURL == "" {
		baseURL = strings.TrimSpace(v.GetString("photoprism_url"))
	}
	if token == "" {
		token = strings.TrimSpace(v.GetString("photoprism_token"))
	}
	return baseURL, token
}

// TestClient_Live_ValidToken runs against a real PhotoPrism instance when
// KIOSK_PHOTOPRISM_URL and KIOSK_PHOTOPRISM_TOKEN are set (or photoprism_url
// and photoprism_token in config.yaml). It verifies that Bearer token
// authentication works and the API returns a valid photos response.
func TestClient_Live_ValidToken(t *testing.T) {
	baseURL, token := livePhotoPrismCreds()
	if baseURL == "" || token == "" {
		t.Skip("live PhotoPrism test skipped: set KIOSK_PHOTOPRISM_URL and KIOSK_PHOTOPRISM_TOKEN, or photoprism_url and photoprism_token in config.yaml")
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
// KIOSK_PHOTOPRISM_URL is set (or photoprism_url in config.yaml). It verifies
// that an invalid token is rejected (401 or error).
func TestClient_Live_InvalidToken(t *testing.T) {
	baseURL, _ := livePhotoPrismCreds()
	if baseURL == "" {
		t.Skip("live PhotoPrism test skipped: set KIOSK_PHOTOPRISM_URL or photoprism_url in config.yaml")
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
