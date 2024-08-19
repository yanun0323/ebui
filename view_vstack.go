package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implement */
var _ SomeView = (*vstackView)(nil)

func VStack(views ...View) *vstackView {
	v := &vstackView{}
	v.viewOption = newViewOption(v, views...)
	return v
}

type vstackView struct {
	viewOption
}

func (v *vstackView) Body() SomeView {
	return v
}

func (v *vstackView) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	current := parent.calculateViewOption(v.viewOption)
	logs.Debugf("VStackView: %+v", current)
	current.Draw(screen, func(screen *ebiten.Image) {
		dy := 0
		current.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, &current)
			if v != nil {
				r := v.draw(screen, current.CreateChild(v, current.X(), current.Y()+dy, current.XX(), current.YY()+dy, current.Width(), current.flexibleH))

				dy += r.Dy()
			}
		})
	})

	return current.DrawnArea()
}

func (v *vstackView) bounds() (int, int) {
	return v.w, v.h
}
