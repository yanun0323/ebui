package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func makeImageRounded(img *ebiten.Image, round int) {
	if img == nil || round <= 0 {
		return
	}

	h, w := img.Bounds().Dy(), img.Bounds().Dx()

	for x := 0; x < round; x++ {
		// left-top
		ltX, ltY := round, round
		for y := 0; y < round; y++ {
			dx := x - ltX
			dy := y - ltY
			if dx*dx+dy*dy > round*round {
				img.Set(x, y, color.Transparent)
			}
		}

		// left-bottom
		lbX, lbY := round, h-round
		for y := h - round; y < h; y++ {
			dx := x - lbX
			dy := y - lbY
			if dx*dx+dy*dy > round*round {
				img.Set(x, y, color.Transparent)
			}
		}
	}

	for x := w - round; x < w; x++ {
		// right-top
		rtX, rtY := w-round, round
		for y := 0; y < round; y++ {
			dx := x - rtX
			dy := y - rtY
			if dx*dx+dy*dy > round*round {
				img.Set(x, y, color.Transparent)
			}
		}

		// right-bottom
		rbX, rbY := w-round, h-round
		for y := h - round; y < h; y++ {
			dx := x - rbX
			dy := y - rbY
			if dx*dx+dy*dy > round*round {
				img.Set(x, y, color.Transparent)
			}
		}
	}
}
