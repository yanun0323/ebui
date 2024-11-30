package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type circle struct {
	view
}

func Circle() SomeView {
	v := &circle{}
	v.view = newView(v)
	return v
}

func (v *circle) draw(screen *ebiten.Image) {
	v.modify()

	size := v.param.frameSize
	clr := v.param.foregroundColor
	diameter := max(size.W, size.H)
	radius := diameter / 2
	img := ebiten.NewImage(diameter, diameter)

	var (
		dx, dy     int
		delta, sqr int
	)

	for x := 0; x < diameter; x++ {
		for y := 0; y < diameter; y++ {
			dx = x - radius
			dy = y - radius
			delta = dx*dx + dy*dy
			sqr = radius * radius
			if delta < sqr {
				img.Set(x, y, clr)
			}
		}
	}

	screen.DrawImage(img, v.drawOption())
}
