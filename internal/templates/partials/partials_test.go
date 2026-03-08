package partials

import (
	"testing"

	"github.com/damongolding/immich-kiosk/internal/source"
)

func TestAssetCameraData(t *testing.T) {
	tests := []struct {
		name string
		exif source.ExifInfo
		want string
	}{
		{
			name: "Make already in model (whole word)",
			exif: source.ExifInfo{
				Make:  "Canon",
				Model: "Canon EOS 5D Mark IV",
			},
			want: "Canon EOS 5D Mark IV",
		},
		{
			name: "Make not in model",
			exif: source.ExifInfo{
				Make:  "Canon",
				Model: "EOS 5D Mark IV",
			},
			want: "Canon EOS 5D Mark IV",
		},
		{
			name: "Make is substring but not whole word",
			exif: source.ExifInfo{
				Make:  "Canon",
				Model: "Canonic EOS 5D Mark IV",
			},
			want: "Canon Canonic EOS 5D Mark IV",
		},
		{
			name: "Extra whitespace trimmed",
			exif: source.ExifInfo{
				Make:  " Canon ",
				Model: " EOS 5D Mark IV ",
			},
			want: "Canon EOS 5D Mark IV",
		},
		{
			name: "Empty make",
			exif: source.ExifInfo{
				Make:  "",
				Model: "EOS 5D Mark IV",
			},
			want: "EOS 5D Mark IV",
		},
		{
			name: "Empty model",
			exif: source.ExifInfo{
				Make:  "Canon",
				Model: "",
			},
			want: "Canon",
		},
		{
			name: "Both empty",
			exif: source.ExifInfo{
				Make:  "",
				Model: "",
			},
			want: "",
		},
		{
			name: "Case insensitive whole word match",
			exif: source.ExifInfo{
				Make:  "canon",
				Model: "Canon EOS 5D Mark IV",
			},
			want: "Canon EOS 5D Mark IV",
		},
		{
			name: "My camera",
			exif: source.ExifInfo{
				Make:  "NIKON CORPORATION",
				Model: "NIKON D90",
			},
			want: "NIKON CORPORATION NIKON D90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssetCameraData(tt.exif)
			if got != tt.want {
				t.Errorf("AssetCameraData() = %q, want %q", got, tt.want)
			}
		})
	}
}
