package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// rpZero returns result if v is zero value, otherwise returns v.
func rpZero[T comparable](v, result T) T {
	var zero T
	return rpEq(v, zero, result)
}

// rpEq returns result if v is eq, otherwise returns v.
func rpEq[T comparable](v, eq, result T) T {
	if v == eq {
		return result
	}

	return v
}

// rpNeq returns result if v is not n, otherwise returns v.
func rpNeq[T comparable](v, n, result T) T {
	if v != n {
		return result
	}

	return v
}

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
