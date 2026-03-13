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
