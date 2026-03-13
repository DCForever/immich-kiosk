package collage

import (
	"math"

	"github.com/damongolding/immich-kiosk/internal/source"
)

// computeTreemap arranges photos using treemap-style recursive subdivision.
// Photos get area proportional to aspect-ratio-adjusted size.
func computeTreemap(assets []source.DisplayAsset, width, height float64, algo string) (*LayoutResult, error) {
	ratios := make([]float64, len(assets))
	areas := make([]float64, len(assets))
	totalArea := 0.0
	for i := range assets {
		w := float64(assets[i].ExifInfo.ExifImageWidth)
		h := float64(assets[i].ExifInfo.ExifImageHeight)
		if w <= 0 || h <= 0 {
			w, h = 1, 1
		}
		ratios[i] = w / h
		areas[i] = math.Sqrt(w * h)
		totalArea += areas[i]
	}
	if totalArea <= 0 {
		totalArea = 1
	}

	const gapFrac = 0.008
	availW := width * (1 - gapFrac*2)
	availH := height * (1 - gapFrac*2)

	cells := subdivideTreemap(ratios, areas, totalArea, 0, 0, availW, availH, width, height, 0)
	return &LayoutResult{
		Cells:          cells,
		Algorithm:      algo,
		ContainerWidth:  width,
		ContainerHeight: height,
	}, nil
}

func subdivideTreemap(ratios, areas []float64, totalArea, x, y, w, h, containerW, containerH float64, baseIndex int) []LayoutCell {
	if len(ratios) == 0 {
		return nil
	}
	if len(ratios) == 1 {
		ar := ratios[0]
		cellW := math.Min(w, h*ar)
		cellH := cellW / ar
		if cellH > h {
			cellH = h
			cellW = cellH * ar
		}
		return []LayoutCell{{
			X:          (x + (w-cellW)/2) / containerW,
			Y:          (y + (h-cellH)/2) / containerH,
			Width:      cellW / containerW,
			Height:     cellH / containerH,
			AssetIndex: baseIndex,
		}}
	}

	mid := len(ratios) / 2
	leftArea := 0.0
	for i := 0; i < mid; i++ {
		leftArea += areas[i]
	}
	rightArea := totalArea - leftArea
	if rightArea <= 0 {
		rightArea = 1
	}

	var cells []LayoutCell
	if w >= h {
		leftW := w * (leftArea / totalArea)
		leftCells := subdivideTreemap(ratios[:mid], areas[:mid], leftArea, x, y, leftW, h, containerW, containerH, baseIndex)
		rightCells := subdivideTreemap(ratios[mid:], areas[mid:], rightArea, x+leftW, y, w-leftW, h, containerW, containerH, baseIndex+mid)
		cells = append(leftCells, rightCells...)
	} else {
		leftH := h * (leftArea / totalArea)
		leftCells := subdivideTreemap(ratios[:mid], areas[:mid], leftArea, x, y, w, leftH, containerW, containerH, baseIndex)
		rightCells := subdivideTreemap(ratios[mid:], areas[mid:], rightArea, x, y+leftH, w, h-leftH, containerW, containerH, baseIndex+mid)
		cells = append(leftCells, rightCells...)
	}
	return cells
}
