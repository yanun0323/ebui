package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
	view modifier
*/

type viewModifier func(screen *ebiten.Image, current *view) SomeView

func backgroundColorViewModifier(clr color.Color) viewModifier {
	return func(screen *ebiten.Image, current *view) SomeView {
		if current == nil {
			return nil
		}

		if screen == nil {
			return nil
		}

		// bColor := current.bColor
		// if bColor == nil {
		// 	bColor = clr
		// }

		w, h := current.Width(), current.Height()
		if w <= 0 || h <= 0 {
			return nil
		}

		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(current.opacity())
		op.GeoM.Translate(float64(current.paddingLeft), float64(current.paddingTop))
		img := ebiten.NewImage(w, h)
		img.Fill(clr)
		screen.DrawImage(img, op)

		return nil
	}
}

func paddingViewModifier(top, right, bottom, left int) viewModifier {
	return func(screen *ebiten.Image, current *view) SomeView {
		if current == nil {
			return nil
		}

		current.paddingTop += top
		current.paddingBottom += bottom
		current.paddingLeft += left
		current.paddingRight += right

		return nil
	}
}
