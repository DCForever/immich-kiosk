// Package collage provides a dynamic layout engine for memories collages.
// It computes arrangements from photo aspect ratios and outputs layout data
// for the Templ collage template.
package collage

import (
	"context"
	"errors"
	"math/rand"

	"github.com/damongolding/immich-kiosk/internal/source"
)

// ErrInvalidAssetCount is returned when asset count is outside 3–16.
var ErrInvalidAssetCount = errors.New("collage: asset count must be 3–16")

// LayoutResult holds the computed layout for a collage.
// Passed to the template for rendering.
type LayoutResult struct {
	Cells          []LayoutCell
	Algorithm      string  // "justified", "treemap", "hero"
	ContainerWidth  float64
	ContainerHeight float64
}

// LayoutCell holds position and size for a single photo in the collage.
// Positions and sizes are in normalized units (0–1) for percentage-based CSS.
type LayoutCell struct {
	X          float64 // Left position (0–1)
	Y          float64 // Top position (0–1)
	Width      float64 // Cell width (0–1)
	Height     float64 // Cell height (0–1)
	AssetIndex int     // Index into original assets slice
}

// DefaultContainerWidth is the default width for layout computation (logical units).
const DefaultContainerWidth = 1920

// DefaultContainerHeight is the default height for layout computation (logical units).
const DefaultContainerHeight = 1080

// ComputeLayout computes a dynamic layout for the given assets.
// Uses width and height for container dimensions; selects an algorithm at random.
// Returns error if assets are invalid or layout cannot be computed.
func ComputeLayout(ctx context.Context, assets []source.DisplayAsset, width, height float64) (*LayoutResult, error) {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	if height <= 0 {
		height = DefaultContainerHeight
	}
	if len(assets) < source.MinCollageAssets || len(assets) > source.MaxCollageAssets {
		return nil, ErrInvalidAssetCount
	}

	// For now, use justified only; treemap and hero will be added in US2
	algo := "justified"
	return computeJustified(assets, width, height, algo)
}

// algorithmNames lists available algorithms for random selection (US2).
var algorithmNames = []string{"justified", "treemap", "hero"}

// pickAlgorithm returns a random algorithm name.
func pickAlgorithm() string {
	return algorithmNames[rand.Intn(len(algorithmNames))]
}
