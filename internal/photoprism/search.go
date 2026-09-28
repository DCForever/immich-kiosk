package photoprism

import (
	"net/url"
	"strings"

	"github.com/damongolding/immich-kiosk/internal/config"
)

// photoPrismUIDCharset is the stock PhotoPrism base32 alphabet.
const photoPrismUIDCharset = "abcdefghijklmnopqrstuvwxyz234567"

// buildPersonSearchQuery returns a stock PhotoPrism search filter.
// RequireAllPeople with several configured names uses people:"A & B".
// A subject UID uses face:<uid>. Any other value uses person:"Name".
func buildPersonSearchQuery(cfg config.Config, singlePersonID string) string {
	if cfg.RequireAllPeople && len(cfg.People) > 1 {
		return peopleSearchQuery(cfg.People)
	}
	return personOrFaceQuery(singlePersonID)
}

func peopleSearchQuery(names []string) string {
	escaped := make([]string, 0, len(names))
	for _, name := range names {
		escaped = append(escaped, strings.ReplaceAll(name, `"`, `\"`))
	}
	return `people:"` + strings.Join(escaped, " & ") + `"`
}

func personOrFaceQuery(id string) string {
	id = strings.TrimSpace(id)
	if isPhotoPrismUID(id) {
		return "face:" + id
	}
	return `person:"` + strings.ReplaceAll(id, `"`, `\"`) + `"`
}

func isPhotoPrismUID(id string) bool {
	if len(id) < 16 {
		return false
	}
	for _, r := range id {
		if !strings.ContainsRune(photoPrismUIDCharset, r) {
			return false
		}
	}
	return true
}

func stockAlbumsQuery() url.Values {
	q := url.Values{}
	q.Set("count", "500")
	q.Set("type", "album")
	return q
}

func regularAlbums(list []Album) []Album {
	out := make([]Album, 0, len(list))
	for _, album := range list {
		if isRegularAlbum(album.Type) {
			out = append(out, album)
		}
	}
	return out
}

func isRegularAlbum(kind string) bool {
	kind = strings.ToLower(strings.TrimSpace(kind))
	return kind == "" || kind == "album"
}
