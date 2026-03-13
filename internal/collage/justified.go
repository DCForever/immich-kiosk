package collage

import (
	"math"

	"github.com/damongolding/immich-kiosk/internal/source"
)

// computeJustified arranges photos in justified rows (similar height per row).
// Preserves aspect ratio; scales photos to fit row width.
// Returns layout with normalized positions (0–1).
func computeJustified(assets []source.DisplayAsset, width, height float64, algo string) (*LayoutResult, error) {
	ratios := make([]float64, len(assets))
	for i := range assets {
		w := float64(assets[i].ExifInfo.ExifImageWidth)
		h := float64(assets[i].ExifInfo.ExifImageHeight)
		if w <= 0 || h <= 0 {
			w, h = 1, 1
		}
		ratios[i] = w / h
	}

	const gapFrac = 0.008 // ~0.8% gap
	availW := width * (1 - gapFrac*2)
	availH := height * (1 - gapFrac*2)

	n := len(ratios)
	numRows := int(math.Ceil(math.Sqrt(float64(n))))
	if numRows < 1 {
		numRows = 1
	}
	if numRows > n {
		numRows = n
	}

	cells := make([]LayoutCell, 0, n)
	photosPerRow := (n + numRows - 1) / numRows
	y := gapFrac

	for row := 0; row < numRows; row++ {
		start := row * photosPerRow
		end := start + photosPerRow
		if end > n {
			end = n
		}
		if start >= n {
			break
		}

		sumRatios := 0.0
		for i := start; i < end; i++ {
			sumRatios += ratios[i]
		}
		if sumRatios <= 0 {
			sumRatios = 1
		}

		rowHeight := availH / float64(numRows)
		scaleToWidth := availW / sumRatios
		if scaleToWidth < rowHeight {
			rowHeight = scaleToWidth
		}

		x := gapFrac
		for i := start; i < end; i++ {
			cellW := (ratios[i] / sumRatios) * availW
			cellH := rowHeight

			cells = append(cells, LayoutCell{
				X:          x / width,
				Y:          y / height,
				Width:      cellW / width,
				Height:     cellH / height,
				AssetIndex: i,
			})
			x += cellW/width + gapFrac
		}

		y += rowHeight/height + gapFrac
	}

	return &LayoutResult{
		Cells:          cells,
		Algorithm:      algo,
		ContainerWidth:  width,
		ContainerHeight: height,
	}, nil
}
