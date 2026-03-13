package collage

import (
	"math"

	"github.com/damongolding/immich-kiosk/internal/source"
)

// computeJustified arranges photos in justified rows (similar height per row).
// Uses aspect ratio range (up to 10% crop) to improve packing.
// Returns layout with normalized positions (0–1).
func computeJustified(assets []source.DisplayAsset, width, height float64, algo string) (*LayoutResult, error) {
	ratios := make([]float64, len(assets))
	minRatios := make([]float64, len(assets))
	maxRatios := make([]float64, len(assets))
	for i := range assets {
		w := float64(assets[i].ExifInfo.ExifImageWidth)
		h := float64(assets[i].ExifInfo.ExifImageHeight)
		if w <= 0 || h <= 0 {
			w, h = 1, 1
		}
		ratios[i] = w / h
		minRatios[i], maxRatios[i] = AspectRange(w, h)
	}

	const gapFrac = 0.008 // ~0.8% gap
	availW := width * (1 - gapFrac*2)
	availH := height * (1 - gapFrac*2)
	cells := packJustifiedRows(ratios, minRatios, maxRatios, availW, availH, width, height, gapFrac)

	return &LayoutResult{
		Cells:          cells,
		Algorithm:      algo,
		ContainerWidth:  width,
		ContainerHeight: height,
	}, nil
}

// packJustifiedRows packs photos into justified rows. Uses aspect ratio range to improve packing.
// Returns cells with normalized positions (0–1).
func packJustifiedRows(ratios, minRatios, maxRatios []float64, availW, availH, width, height, gapFrac float64) []LayoutCell {
	n := len(ratios)
	if n == 0 {
		return nil
	}
	numRows := int(math.Ceil(math.Sqrt(float64(n))))
	if numRows < 1 {
		numRows = 1
	}
	if numRows > n {
		numRows = n
	}

	cells := make([]LayoutCell, 0, n)
	photosPerRow := (n + numRows - 1) / numRows
	y := gapFrac * height

	for row := 0; row < numRows; row++ {
		start := row * photosPerRow
		end := start + photosPerRow
		if end > n {
			end = n
		}
		if start >= n {
			break
		}

		numInRow := end - start
		rowGaps := float64(numInRow-1) * gapFrac * width
		rowAvailW := availW - rowGaps
		rowGapTotal := float64(numRows-1) * gapFrac * height
		rowHeight := (availH - rowGapTotal) / float64(numRows)

		// Choose ratios within [min,max] to fill row width (after reserving gaps)
		chosen := chooseRatiosForRow(ratios[start:end], minRatios[start:end], maxRatios[start:end], rowAvailW, rowHeight)

		sumChosen := 0.0
		for _, r := range chosen {
			sumChosen += r
		}
		if sumChosen <= 0 {
			sumChosen = 1
		}

		scaleToWidth := rowAvailW / sumChosen
		if scaleToWidth < rowHeight {
			rowHeight = scaleToWidth
		}

		x := gapFrac * width
		for i := start; i < end; i++ {
			cellW := (chosen[i-start] / sumChosen) * rowAvailW
			cellH := rowHeight

			cells = append(cells, LayoutCell{
				X:          x / width,
				Y:          y / height,
				Width:      cellW / width,
				Height:     cellH / height,
				AssetIndex: i,
			})
			x += cellW + gapFrac*width
		}

		y += rowHeight + gapFrac*height
	}

	return cells
}

// chooseRatiosForRow picks ratios in [min,max] so sum(chosen) ≈ availW/rowHeight.
// Uses linear interpolation when target is within [sum(min), sum(max)].
func chooseRatiosForRow(ratios, minR, maxR []float64, availW, rowHeight float64) []float64 {
	n := len(ratios)
	if n == 0 {
		return nil
	}
	targetSum := availW / rowHeight

	sumMin := 0.0
	sumMax := 0.0
	for i := 0; i < n; i++ {
		sumMin += minR[i]
		sumMax += maxR[i]
	}

	chosen := make([]float64, n)
	if targetSum <= sumMin {
		copy(chosen, minR)
		return chosen
	}
	if targetSum >= sumMax {
		copy(chosen, maxR)
		return chosen
	}
	// Interpolate: chosen[i] = minR[i] + (maxR[i]-minR[i])*t, sum(chosen)=targetSum
	// sumMin + (sumMax-sumMin)*t = targetSum => t = (targetSum-sumMin)/(sumMax-sumMin)
	span := sumMax - sumMin
	if span <= 0 {
		copy(chosen, ratios)
		return chosen
	}
	t := (targetSum - sumMin) / span
	for i := 0; i < n; i++ {
		chosen[i] = minR[i] + (maxR[i]-minR[i])*t
	}
	return chosen
}
