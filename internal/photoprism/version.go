package photoprism

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

// versionCache remembers one probe per base URL so each slideshow request
// does not wait on GET /api/v1/config. An empty version is a completed probe.
var (
	versionCacheMu sync.Mutex
	versionCache   = map[string]string{}
	versionProbed  = map[string]bool{}
)

const versionProbeTimeout = 5 * time.Second

// ServerVersion returns the PhotoPrism version from GET /api/v1/config.
// The result is cached on the client. A failed probe returns an empty string
// so callers can continue; fork-specific behavior is not gated on this value.
func (c *Client) ServerVersion(ctx context.Context) string {
	version, _ := c.lookupVersion(ctx)
	return version
}

// lookupVersion returns the cached or freshly probed version.
// The boolean is true only for the call that performed the probe.
func (c *Client) lookupVersion(ctx context.Context) (string, bool) {
	if ctx == nil {
		ctx = context.Background()
	}
	if version, ok := loadCachedVersion(c.BaseURL); ok {
		c.noteVersion(version)
		return version, false
	}
	probedNow := false
	c.versionOnce.Do(func() {
		found := c.probeVersion(ctx)
		storeCachedVersion(c.BaseURL, found)
		c.noteVersion(found)
		probedNow = true
	})
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.version, probedNow
}

func (c *Client) noteVersion(version string) {
	if version == "" {
		return
	}
	c.mu.Lock()
	if c.version == "" {
		c.version = version
	}
	c.mu.Unlock()
	upgradeCachedVersion(c.BaseURL, version)
}

func loadCachedVersion(baseURL string) (string, bool) {
	versionCacheMu.Lock()
	defer versionCacheMu.Unlock()
	version, ok := versionCache[baseURL]
	if !ok && !versionProbed[baseURL] {
		return "", false
	}
	return version, versionProbed[baseURL]
}

func storeCachedVersion(baseURL, version string) {
	versionCacheMu.Lock()
	defer versionCacheMu.Unlock()
	if versionProbed[baseURL] && versionCache[baseURL] != "" {
		return
	}
	versionProbed[baseURL] = true
	if versionCache[baseURL] == "" {
		versionCache[baseURL] = version
	}
}

func upgradeCachedVersion(baseURL, version string) {
	if version == "" || baseURL == "" {
		return
	}
	versionCacheMu.Lock()
	defer versionCacheMu.Unlock()
	if versionCache[baseURL] != "" {
		return
	}
	versionCache[baseURL] = version
	versionProbed[baseURL] = true
}

func (c *Client) probeVersion(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, versionProbeTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+apiPrefix+"/config", nil)
	if err != nil {
		return ""
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.client().Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	headers := headerMap(resp.Header)
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if version := versionFromBody(body); version != "" {
		return version
	}
	return headerVersion(headers)
}

func versionFromBody(body []byte) string {
	fields, err := objectFields(body)
	if err != nil || fields == nil {
		return ""
	}
	return decodeString(fields, "version")
}

func headerVersion(headers map[string]string) string {
	return firstHeader(headers, "x-photoprism-version", "x-version")
}

func firstHeader(headers map[string]string, keys ...string) string {
	for _, key := range keys {
		if headers[key] != "" {
			return headers[key]
		}
	}
	return ""
}
