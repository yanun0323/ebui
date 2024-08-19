package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implement */
var _ SomeView = (*hstackView)(nil)

func HStack(views ...View) *hstackView {
	v := &hstackView{}
	v.viewOption = newViewOption(v, views...)
	return v
}

type hstackView struct {
	viewOption
}

func (v *hstackView) Body() SomeView {
	return v
}

func (v *hstackView) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	current := parent.calculateViewOption(v.viewOption)
	logs.Debugf("HStackView: %+v", current)
	current.Draw(screen, func(screen *ebiten.Image) {
		dx := 0
		current.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, &current)
			if v != nil {
				r := v.draw(screen, current.CreateChild(v, current.X()+dx, current.Y(), current.XX()+dx, current.YY(), current.flexibleW, current.Height()))
				dx += r.Dx()
			}
		})
	})

	return current.DrawnArea()
}

func (v *hstackView) bounds() (int, int) {
	return v.w, v.h
}
