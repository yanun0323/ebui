package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
	view modifier
*/

type viewModifier func(screen *ebiten.Image, current *uiView) SomeView

func backgroundColorViewModifier(clr color.Color) viewModifier {
	return func(screen *ebiten.Image, current *uiView) SomeView {
		if current == nil {
			return nil
		}

		if screen == nil {
			return nil
		}

		w, h := current.Width(), current.Height()
		if w <= 0 || h <= 0 {
			return nil
		}

		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(current.opacity())
		op.GeoM.Translate(float64(current.padding.left), float64(current.padding.top))
		img := ebiten.NewImage(w, h)
		img.Fill(clr)
		screen.DrawImage(img, op)

		return nil
	}
}

func paddingViewModifier(top, right, bottom, left int) viewModifier {
	return func(screen *ebiten.Image, current *uiView) SomeView {
		if current == nil {
			return nil
		}

		current.padding.top += top
		current.padding.bottom += bottom
		current.padding.left += left
		current.padding.right += right

		return nil
	}
}
