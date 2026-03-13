package collage

import (
	"math"

	"github.com/damongolding/immich-kiosk/internal/source"
)

// computeTreemap arranges photos using treemap-style recursive subdivision.
// Photos get area proportional to aspect-ratio-adjusted size.
// Uses aspect ratio range (up to 10% crop) to improve packing.
func computeTreemap(assets []source.DisplayAsset, width, height float64, algo string) (*LayoutResult, error) {
	ratios := make([]float64, len(assets))
	minRatios := make([]float64, len(assets))
	maxRatios := make([]float64, len(assets))
	areas := make([]float64, len(assets))
	totalArea := 0.0
	for i := range assets {
		w := float64(assets[i].ExifInfo.ExifImageWidth)
		h := float64(assets[i].ExifInfo.ExifImageHeight)
		if w <= 0 || h <= 0 {
			w, h = 1, 1
		}
		ratios[i] = w / h
		minRatios[i], maxRatios[i] = AspectRange(w, h)
		areas[i] = math.Sqrt(w * h)
		totalArea += areas[i]
	}
	if totalArea <= 0 {
		totalArea = 1
	}

	const gapFrac = 0.008
	availW := width * (1 - gapFrac*2)
	availH := height * (1 - gapFrac*2)
	offsetX := gapFrac * width
	offsetY := gapFrac * height

	cells := subdivideTreemap(ratios, minRatios, maxRatios, areas, totalArea, offsetX, offsetY, availW, availH, width, height, 0)
	return &LayoutResult{
		Cells:          cells,
		Algorithm:      algo,
		ContainerWidth:  width,
		ContainerHeight: height,
	}, nil
}

func subdivideTreemap(ratios, minRatios, maxRatios, areas []float64, totalArea, x, y, w, h, containerW, containerH float64, baseIndex int) []LayoutCell {
	if len(ratios) == 0 {
		return nil
	}
	if len(ratios) == 1 {
		minR, maxR := minRatios[0], maxRatios[0]
		ar := ratios[0]
		// Choose ratio in [minR, maxR] that maximizes cell area within rect (w,h)
		// Cell fits if cellW<=w and cellH<=h. For ratio r: cellW=cellH*r, so cellH*r<=w, cellH<=h => cellH=min(h, w/r)
		// We want max area = cellW*cellH = r*cellH^2. Try r at min and max; pick the one that fits and has larger area.
		cellW, cellH := fitCellInRect(ar, minR, maxR, w, h)
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
		leftCells := subdivideTreemap(ratios[:mid], minRatios[:mid], maxRatios[:mid], areas[:mid], leftArea, x, y, leftW, h, containerW, containerH, baseIndex)
		rightCells := subdivideTreemap(ratios[mid:], minRatios[mid:], maxRatios[mid:], areas[mid:], rightArea, x+leftW, y, w-leftW, h, containerW, containerH, baseIndex+mid)
		cells = append(leftCells, rightCells...)
	} else {
		leftH := h * (leftArea / totalArea)
		leftCells := subdivideTreemap(ratios[:mid], minRatios[:mid], maxRatios[:mid], areas[:mid], leftArea, x, y, w, leftH, containerW, containerH, baseIndex)
		rightCells := subdivideTreemap(ratios[mid:], minRatios[mid:], maxRatios[mid:], areas[mid:], rightArea, x, y+leftH, w, h-leftH, containerW, containerH, baseIndex+mid)
		cells = append(leftCells, rightCells...)
	}
	return cells
}

// fitCellInRect returns cell dimensions (cellW, cellH) for an image with aspect range [minR, maxR]
// that fit in rect (w,h) and maximize area. Uses original ratio ar as hint; result is clamped to [minR, maxR].
func fitCellInRect(ar, minR, maxR, w, h float64) (cellW, cellH float64) {
	// Try to fit with ratio ar (clamped to [minR, maxR])
	r := ar
	if r < minR {
		r = minR
	}
	if r > maxR {
		r = maxR
	}
	cellW = math.Min(w, h*r)
	cellH = cellW / r
	if cellH > h {
		cellH = h
		cellW = cellH * r
	}
	// Ensure resulting ratio is in [minR, maxR] (may have been constrained by rect)
	actualRatio := cellW / cellH
	if actualRatio < minR {
		cellH = cellW / minR
		if cellH > h {
			cellH = h
			cellW = cellH * minR
		}
	} else if actualRatio > maxR {
		cellW = cellH * maxR
		if cellW > w {
			cellW = w
			cellH = cellW / maxR
		}
	}
	return cellW, cellH
}
