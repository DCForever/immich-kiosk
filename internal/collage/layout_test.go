package collage

import (
	"context"
	"testing"

	"github.com/damongolding/immich-kiosk/internal/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayoutResult_types(t *testing.T) {
	result := &LayoutResult{
		Cells: []LayoutCell{
			{X: 0, Y: 0, Width: 0.5, Height: 0.5, AssetIndex: 0},
			{X: 0.5, Y: 0, Width: 0.5, Height: 0.5, AssetIndex: 1},
		},
		Algorithm:      "justified",
		ContainerWidth:  1920,
		ContainerHeight: 1080,
	}
	assert.Len(t, result.Cells, 2)
	assert.Equal(t, 0, result.Cells[0].AssetIndex)
	assert.Equal(t, 0.5, result.Cells[0].Width)
}

func TestComputeLayout_invalidCount(t *testing.T) {
	ctx := context.Background()
	assets := make([]source.DisplayAsset, 2) // too few
	_, err := ComputeLayout(ctx, assets, 1920, 1080)
	assert.ErrorIs(t, err, ErrInvalidAssetCount)

	assets = make([]source.DisplayAsset, 20) // too many
	_, err = ComputeLayout(ctx, assets, 1920, 1080)
	assert.ErrorIs(t, err, ErrInvalidAssetCount)
}

func mkAsset(w, h int) source.DisplayAsset {
	return source.DisplayAsset{
		ExifInfo: source.ExifInfo{ExifImageWidth: w, ExifImageHeight: h},
	}
}

func TestComputeLayout_justified(t *testing.T) {
	ctx := context.Background()
	assets := []source.DisplayAsset{
		mkAsset(1920, 1080), // 16:9
		mkAsset(1080, 1920), // 9:16
		mkAsset(1000, 1000), // 1:1
	}

	result, err := ComputeLayout(ctx, assets, 1920, 1080)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "justified", result.Algorithm)
	assert.Len(t, result.Cells, 3)

	for i, c := range result.Cells {
		assert.Equal(t, i, c.AssetIndex)
		assert.Greater(t, c.Width, 0.0)
		assert.Greater(t, c.Height, 0.0)
		assert.GreaterOrEqual(t, c.X, 0.0)
		assert.GreaterOrEqual(t, c.Y, 0.0)
		assert.LessOrEqual(t, c.X+c.Width, 1.01) // allow small float tolerance
		assert.LessOrEqual(t, c.Y+c.Height, 1.01)
	}
}

func TestComputeLayout_zeroDimensions(t *testing.T) {
	ctx := context.Background()
	assets := []source.DisplayAsset{
		mkAsset(0, 0), // fallback to square
		mkAsset(0, 0),
		mkAsset(0, 0),
	}

	result, err := ComputeLayout(ctx, assets, 1920, 1080)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Cells, 3)
}

func TestComputeLayout_allAlgorithms(t *testing.T) {
	ctx := context.Background()
	assets := []source.DisplayAsset{
		mkAsset(1920, 1080),
		mkAsset(1080, 1920),
		mkAsset(1000, 1000),
		mkAsset(800, 600),
		mkAsset(600, 800),
	}

	seen := make(map[string]bool)
	for i := 0; i < 30; i++ {
		result, err := ComputeLayout(ctx, assets, 1920, 1080)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Cells, 5)
		seen[result.Algorithm] = true
	}
	assert.True(t, len(seen) >= 2, "expected at least 2 different algorithms over 30 runs, got %v", seen)
}

func TestComputeLayout_edgeCases(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name   string
		assets []source.DisplayAsset
	}{
		{"3 photos", []source.DisplayAsset{mkAsset(100, 100), mkAsset(200, 200), mkAsset(300, 300)}},
		{"16 photos", make([]source.DisplayAsset, 16)},
		{"panorama", []source.DisplayAsset{mkAsset(4000, 500), mkAsset(1000, 1000), mkAsset(500, 4000)}},
		{"portrait", []source.DisplayAsset{mkAsset(100, 1000), mkAsset(500, 500), mkAsset(1000, 100)}},
	}
	for i := range tests[1].assets {
		tests[1].assets[i] = mkAsset(1000+i*100, 800)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ComputeLayout(ctx, tt.assets, 1920, 1080)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Len(t, result.Cells, len(tt.assets))
			for i, c := range result.Cells {
				assert.Greater(t, c.Width, 0.0, "cell %d width", i)
				assert.Greater(t, c.Height, 0.0, "cell %d height", i)
				assert.GreaterOrEqual(t, c.X, 0.0)
				assert.GreaterOrEqual(t, c.Y, 0.0)
				assert.LessOrEqual(t, c.X+c.Width, 1.01)
				assert.LessOrEqual(t, c.Y+c.Height, 1.01)
			}
		})
	}
}
