package immich

import (
	"net/http"
	"net/url"
	"path"
)

// Video retrieves the video asset from Immich server.
// Returns the video data as a byte slice, the contentType, and any error encountered.
// The video is returned in octet-stream format.
func (a *Asset) Video() ([]byte, string, error) {
	var responseBody []byte
	var contentType string

	u, err := url.Parse(a.requestConfig.ImmichURL)
	if err != nil {
		return responseBody, "", err
	}

	apiURL := url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   path.Join("api", "assets", a.ID, "video", "playback"),
	}

	octetStreamHeader := map[string]string{"Accept": "application/octet-stream"}

	responseBody, contentType, _, err = a.immichAPICall(a.ctx, http.MethodGet, apiURL.String(), nil, octetStreamHeader)
	if err != nil {
		return responseBody, contentType, err
	}

	return responseBody, contentType, nil
}

// durationCheck verifies that the video duration string in the Asset is valid and represents
// a duration of at least one second.
// Has been substantially simplified from the original implementation but to Immich V3.
//
// Returns true if the duration is valid and at least one second, false otherwise.
func (a *Asset) durationCheck() bool {
	// Ignore videos with duration less than one second
	if a.Duration < 1000 {
		return false
	}

	// Check maximum duration if configured
	if a.requestConfig.ExcludeVideosOver > 0 && a.Duration > int64(a.requestConfig.ExcludeVideosOver)*1000 {
		return false
	}

	return true
}
