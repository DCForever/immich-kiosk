package photoprism

import (
	"fmt"
	"strings"
)

// thumbSizes are stock PhotoPrism fit sizes. fit_720 is tried first; later
// sizes are used only when an earlier request fails.
var thumbSizes = []string{"fit_720", "fit_1280", "fit_1920"}

func (p *Provider) ImagePreview() ([]byte, string, error) {
	ph := p.currentPhoto()
	if ph == nil {
		return nil, "", fmt.Errorf("photoprism: no current photo")
	}
	return p.previewWithSizes(ph.PrimaryHash(), thumbSizes)
}

func (p *Provider) previewWithSizes(hash string, sizes []string) ([]byte, string, error) {
	token := p.previewToken
	if token == "" {
		token = "public"
	}
	var last error
	for _, size := range sizes {
		path := fmt.Sprintf("%s/t/%s/%s/%s", apiPrefix, hash, token, size)
		body, contentType, err := p.client.getBytes(p.ctx, path)
		if err == nil {
			return body, contentType, nil
		}
		last = err
	}
	if last == nil {
		last = fmt.Errorf("photoprism: no thumbnail")
	}
	return nil, "", last
}

func (p *Provider) Video() ([]byte, string, error) {
	ph := p.currentPhoto()
	if ph == nil || !strings.EqualFold(ph.Type, "video") {
		return nil, "", fmt.Errorf("photoprism: no current video")
	}
	token := p.previewToken
	if token == "" {
		token = "public"
	}
	path := fmt.Sprintf("%s/videos/%s/%s/avc", apiPrefix, ph.PrimaryHash(), token)
	return p.client.getBytes(p.ctx, path)
}
