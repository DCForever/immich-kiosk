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

// AllowedCropFrac is the max fraction to crop from each end of the longer dimension (10%).
const AllowedCropFrac = 0.10

// AspectRange returns [minRatio, maxRatio] for an image; cell aspect must stay in this range.
// Allows up to 10% cropping from each end of the longer dimension.
// - Landscape (r >= 1): min = 0.8*r, max = r (crop left/right)
// - Portrait (r < 1): min = r, max = r/0.8 (crop top/bottom)
// - Square (r = 1): min = 0.8, max = 1.25
func AspectRange(w, h float64) (minRatio, maxRatio float64) {
	if w <= 0 || h <= 0 {
		return 1 - 2*AllowedCropFrac, 1 / (1 - 2*AllowedCropFrac)
	}
	r := w / h
	if r > 1 {
		return (1 - 2*AllowedCropFrac) * r, r
	}
	if r < 1 {
		return r, r / (1 - 2*AllowedCropFrac)
	}
	return 1 - 2*AllowedCropFrac, 1 / (1 - 2*AllowedCropFrac)
}

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

	algo := pickAlgorithm(len(assets))
	switch algo {
	case "treemap":
		return computeTreemap(assets, width, height, algo)
	case "hero":
		return computeHero(assets, width, height, algo)
	default:
		return computeJustified(assets, width, height, algo)
	}
}

// algorithmNames lists available algorithms for random selection (US2).
var algorithmNames = []string{"justified", "treemap", "hero"}

// pickAlgorithm returns a random algorithm name. For 8+ photos, excludes treemap
// because it tends to produce narrow vertical strips in wide containers.
func pickAlgorithm(n int) string {
	candidates := algorithmNames
	if n >= 8 {
		candidates = []string{"justified", "hero"}
	}
	return candidates[rand.Intn(len(candidates))]
}
