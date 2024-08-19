package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
	view modifier
*/

type viewModifier func(screen *ebiten.Image, current *viewOption) SomeView

func newViewModifiers(svs ...View) []viewModifier {
	vms := make([]viewModifier, 0, len(svs))
	for _, sv := range svs {
		body := sv.Body()
		vms = append(vms, func(*ebiten.Image, *viewOption) SomeView {
			return body
		})
	}
	return vms
}

func backgroundColorViewModifier(clr color.Color) viewModifier {
	return func(screen *ebiten.Image, current *viewOption) SomeView {
		if screen == nil || current == nil {
			return nil
		}

		bColor := current.bColor
		if bColor == nil {
			bColor = clr
		}

		w, h := current.Width(), current.Height()
		if w <= 0 || h <= 0 {
			return nil
		}

		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(current.opacity())
		op.GeoM.Translate(float64(current.XX()), float64(current.YY()))
		img := ebiten.NewImage(w, h)
		img.Fill(bColor)
		screen.DrawImage(img, op)

		return nil
	}
}

func paddingViewModifier(top, right, bottom, left int) viewModifier {
	return func(screen *ebiten.Image, current *viewOption) SomeView {
		if screen == nil || current == nil {
			return nil
		}

		current.paddingTop += top
		current.paddingBottom += bottom
		current.paddingLeft += left
		current.paddingRight += right

		current.calculateFlexibleSizeTo()

		return nil
	}
}
