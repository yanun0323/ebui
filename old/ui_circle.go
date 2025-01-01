package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type circle struct {
	*view
}

func Circle() SomeView {
	v := &circle{}
	v.view = newView(idCircle, v)
	return v
}

func (v *circle) getRenderImage() *ebiten.Image {
	if v.noChange.Swap(true) {
		return v.img.Load()
	}

	size := v.render.size
	clr := v.param.foregroundColor.Get()
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

	v.img.Store(img)

	return img
}
