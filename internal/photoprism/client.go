// Package photoprism provides a minimal HTTP client for the PhotoPrism REST API.
// It uses only the Go standard library (net/http, encoding/json). Authentication
// is via Bearer token (app password or session token); set Authorization header
// on each request. See https://docs.photoprism.app/developer-guide/api/auth/.
package photoprism

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiPrefix      = "/api/v1"
	defaultTimeout = 30 * time.Second
)

// Client performs authenticated requests to a PhotoPrism instance.
type Client struct {
	BaseURL    string       // e.g. https://photoprism.example.com (no trailing slash)
	Token      string       // Bearer token: app password (Settings → Account) or session/access token per Client Authentication docs
	HTTPClient *http.Client // nil uses default with timeout
}

// NewClient returns a client for the given base URL and token.
// baseURL must not have a trailing slash. Token is sent as Authorization: Bearer <token>.
func NewClient(baseURL, token string) *Client {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) client() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: defaultTimeout}
}

// do sends an authenticated GET request and returns the body and response for reading headers.
// Caller must close body. Path is the path without base URL (e.g. /api/v1/photos).
// query may be nil.
func (c *Client) do(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	rawURL := c.BaseURL + path
	if len(query) > 0 {
		rawURL += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("photoprism api: %s %s", resp.Status, rawURL)
	}
	return resp, nil
}

// getJSON performs a GET request and decodes the response body into v.
// It reads the response headers (e.g. X-Preview-Token) into the returned map.
func (c *Client) getJSON(ctx context.Context, path string, query url.Values, v any) (headers map[string]string, err error) {
	resp, err := c.do(ctx, path, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	headers = make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[strings.ToLower(k)] = v[0]
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return nil, fmt.Errorf("photoprism api json: %w", err)
		}
	}
	return headers, nil
}

// getBytes performs a GET request and returns the raw body (e.g. for thumbnail/video binary).
func (c *Client) getBytes(ctx context.Context, path string) ([]byte, string, error) {
	rawURL := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := c.client().Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("photoprism api: %s %s", resp.Status, rawURL)
	}
	ct := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return body, ct, nil
}
