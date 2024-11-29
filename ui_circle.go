package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implement */
var _ SomeView = (*circle)(nil)

func Circle() SomeView {
	v := &circle{}
	v.uiView = newView(typesCircle, v)

	return v
}

type circle struct {
	*uiView
}

func (v *circle) draw(screen *ebiten.Image) {
	drawSize := v.getDrawSize(v.cachedSize)
	diameter := min(drawSize.w, drawSize.h)
	radius := diameter / 2

	img := ebiten.NewImage(diameter, diameter)
	clr := defaultForegroundColor
	if v.fColor != nil {
		clr = v.fColor
	}

	cR, cG, cB, cA := clr.RGBA()
	clr2 := color.Color(color.RGBA64{uint16(cR), uint16(cG), uint16(cB), uint16(cA / 2)})

	for x := 0; x < drawSize.w; x++ {
		for y := 0; y < drawSize.h; y++ {
			a := (x-radius)*(x-radius) + (y-radius)*(y-radius)
			b := radius * radius
			if a < b {
				img.Set(x, y, clr)
			} else if a == b {
				img.Set(x, y, clr2)
			}
		}
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(v.start.x), float64(v.start.y))
	screen.DrawImage(img, opt)
	v.drawModifiers(screen)
}
