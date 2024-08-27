package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
	view modifier
*/

type viewModifier func(screen *ebiten.Image, current *uiView)

func backgroundColorViewModifier(clr color.Color) viewModifier {
	return func(screen *ebiten.Image, current *uiView) {
		if current == nil || screen == nil {
			return
		}

		w, h := current.Width(), current.Height()
		if w <= 0 || h <= 0 {
			return
		}

		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(current.opacity())
		op.GeoM.Translate(float64(current.padding.left), float64(current.padding.top))
		img := ebiten.NewImage(w, h)
		img.Fill(clr)
		screen.DrawImage(img, op)
	}
}

func frameViewModifier(w, h int) viewModifier {
	return func(_ *ebiten.Image, current *uiView) {
		if current == nil {
			return
		}

		if w < 0 {
			w = -1
		}

		if h < 0 {
			h = -1
		}

		current.initSize = frame{rpEq(w, -1, current.initSize.w), rpEq(h, -1, current.initSize.h)}
		current.size = frame{rpEq(w, -1, current.size.w), rpEq(h, -1, current.size.h)}
	}
}

func paddingViewModifier(top, right, bottom, left int) viewModifier {
	return func(_ *ebiten.Image, current *uiView) {
		if current == nil {
			return
		}

		if top < 0 {
			top = 0
		}

		if right < 0 {
			right = 0
		}

		if bottom < 0 {
			bottom = 0
		}

		if left < 0 {
			left = 0
		}

		current.padding.left += left
		current.padding.right += right
		current.padding.top += top
		current.padding.bottom += bottom
	}
}

func cornerRadiusViewModifier(radius ...int) viewModifier {
	r := 7
	if len(radius) != 0 && radius[0] >= 0 {
		r = radius[0]
	}

	return func(screen *ebiten.Image, current *uiView) {
		if screen == nil || current == nil {
			return
		}

		if r == 0 {
			return
		}

		makeImageRounded(screen, current, r)
	}
}
