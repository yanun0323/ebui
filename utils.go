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

// makeImageRounded makes image rounded.
//
// FIXME: should be optimized for chaining calls.
func makeImageRounded(img *ebiten.Image, current *uiView, round int) {
	if img == nil || round <= 0 {
		return
	}

	sW, sH := img.Bounds().Dx(), img.Bounds().Dy()
	if sW <= 0 || sH <= 0 {
		return
	}

	set := func(x, y int) {
		// x += current.padding.left
		// y += current.padding.top
		if x >= 0 && y >= 0 && x < sW && y < sH {
			img.Set(x, y, color.Transparent)
		}
	}

	w, h := current.size.w, current.size.h

	for x := 0; x < round; x++ {
		// left-top
		ltX, ltY := round, round
		for y := 0; y < round; y++ {
			dx := x - ltX
			dy := y - ltY
			if dx*dx+dy*dy > round*round {
				set(x, y)
			}
		}

		// left-bottom
		lbX, lbY := round, h-round
		for y := h - round; y < h; y++ {
			dx := x - lbX
			dy := y - lbY
			if dx*dx+dy*dy > round*round {
				set(x, y)
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
				set(x, y)
			}
		}

		// right-bottom
		rbX, rbY := w-round, h-round
		for y := h - round; y < h; y++ {
			dx := x - rbX
			dy := y - rbY
			if dx*dx+dy*dy > round*round {
				set(x, y)
			}
		}
	}
}
