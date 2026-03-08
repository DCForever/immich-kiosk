package components

import (
	"strings"

	"github.com/damongolding/immich-kiosk/internal/source"
	"github.com/damongolding/immich-kiosk/internal/utils"
)

func CreateDataTag(tags source.Tags) string {
	clean := make([]string, 0, len(tags))

	for _, tag := range tags {
		clean = append(clean, utils.SanitizeClassName(tag.Value))
	}

	return strings.Join(clean, " ")
}
