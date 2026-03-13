package collage

import (
	"math"
	"math/rand"

	"github.com/damongolding/immich-kiosk/internal/source"
)

// computeHero arranges one large photo (hero) plus remaining photos in a cluster.
// Hero is chosen at random; cluster uses justified rows.
func computeHero(assets []source.DisplayAsset, width, height float64, algo string) (*LayoutResult, error) {
	n := len(assets)
	if n == 0 {
		return nil, ErrInvalidAssetCount
	}

	// Pick hero index (0 for deterministic test; in production use rand)
	heroIndex := 0
	if n > 1 {
		heroIndex = rand.Intn(n)
	}

	ratios := make([]float64, n)
	minRatios := make([]float64, n)
	maxRatios := make([]float64, n)
	for i := range assets {
		w := float64(assets[i].ExifInfo.ExifImageWidth)
		h := float64(assets[i].ExifInfo.ExifImageHeight)
		if w <= 0 || h <= 0 {
			w, h = 1, 1
		}
		ratios[i] = w / h
		minRatios[i], maxRatios[i] = AspectRange(w, h)
	}

	const gapFrac = 0.008
	availW := width * (1 - gapFrac*2)
	availH := height * (1 - gapFrac*2)

	// Hero gets ~50% of area; cluster gets the rest
	heroFrac := 0.5
	if n == 1 {
		heroFrac = 1.0
	}

	var cells []LayoutCell

	// Place hero: use ratio within [min,max] to fit available space
	heroMin, heroMax := minRatios[heroIndex], maxRatios[heroIndex]
	heroAr := ratios[heroIndex]
	heroArea := availW * availH * heroFrac
	heroW := math.Sqrt(heroArea * heroAr)
	heroH := heroW / heroAr
	if heroH > availH {
		heroH = availH
		heroW = heroH * heroAr
	}
	if heroW > availW {
		heroW = availW
		heroH = heroW / heroAr
	}
	// Clamp hero aspect to allowed range
	heroCellRatio := heroW / heroH
	if heroCellRatio < heroMin {
		heroCellRatio = heroMin
		heroW = math.Sqrt(heroArea * heroCellRatio)
		heroH = heroW / heroCellRatio
		if heroH > availH {
			heroH = availH
			heroW = heroH * heroCellRatio
		}
		if heroW > availW {
			heroW = availW
			heroH = heroW / heroCellRatio
		}
	} else if heroCellRatio > heroMax {
		heroCellRatio = heroMax
		heroW = math.Sqrt(heroArea * heroCellRatio)
		heroH = heroW / heroCellRatio
		if heroH > availH {
			heroH = availH
			heroW = heroH * heroCellRatio
		}
		if heroW > availW {
			heroW = availW
			heroH = heroW / heroCellRatio
		}
	}
	// When cluster will be to the side, left-align hero to make room; otherwise center hero
	clusterToSide := false
	heroX := (gapFrac + (availW-heroW)/2) / width
	if n > 1 {
		clusterH := availH - heroH - gapFrac*height
		if clusterH < availH*0.2 {
			clusterToSide = true
			heroX = gapFrac / width // Left-align hero to make room for cluster on the right
		}
	}

	cells = append(cells, LayoutCell{
		X:          heroX,
		Y:          gapFrac / height,
		Width:      heroW / width,
		Height:     heroH / height,
		AssetIndex: heroIndex,
	})

	if n > 1 {
		// Cluster: remaining photos in justified rows below or beside hero
		clusterW := availW
		clusterH := availH - heroH - gapFrac*height
		clusterY := heroH + gapFrac*height
		if clusterToSide {
			clusterW = availW - heroW - 2*gapFrac*width // Space to right of hero minus gaps
			clusterH = availH
			clusterY = 0
		}

		// Cluster starts at hero's right edge + gap
		clusterOffsetX := (gapFrac*width + heroW + gapFrac*width) / width
		clusterRatios := make([]float64, 0, n-1)
		clusterMinRatios := make([]float64, 0, n-1)
		clusterMaxRatios := make([]float64, 0, n-1)
		clusterIndices := make([]int, 0, n-1)
		for i := 0; i < n; i++ {
			if i != heroIndex {
				clusterRatios = append(clusterRatios, ratios[i])
				clusterMinRatios = append(clusterMinRatios, minRatios[i])
				clusterMaxRatios = append(clusterMaxRatios, maxRatios[i])
				clusterIndices = append(clusterIndices, i)
			}
		}
		clusterCells := packJustifiedRowsForHero(clusterRatios, clusterMinRatios, clusterMaxRatios, clusterIndices, clusterW, clusterH, width, height, gapFrac, clusterY/height, clusterToSide, clusterOffsetX)
		cells = append(cells, clusterCells...)
	}

	return &LayoutResult{
		Cells:          cells,
		Algorithm:      algo,
		ContainerWidth:  width,
		ContainerHeight: height,
	}, nil
}

func packJustifiedRowsForHero(ratios, minRatios, maxRatios []float64, indices []int, containerW, containerH, totalW, totalH, gapFrac, offsetY float64, clusterToSide bool, clusterOffsetX float64) []LayoutCell {
	if len(ratios) == 0 {
		return nil
	}
	availW := containerW * (1 - gapFrac*2)
	availH := containerH * (1 - gapFrac*2)
	cells := packJustifiedRows(ratios, minRatios, maxRatios, availW, availH, containerW, containerH, gapFrac)
	for i := range cells {
		cells[i].AssetIndex = indices[cells[i].AssetIndex]
		cells[i].X = cells[i].X * (containerW / totalW)
		if clusterToSide {
			cells[i].X += clusterOffsetX
		}
		cells[i].Y = offsetY + cells[i].Y*(containerH/totalH)
		cells[i].Width *= containerW / totalW
		cells[i].Height *= containerH / totalH
	}
	return cells
}
