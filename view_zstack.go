package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implement */
var _ SomeView = (*zstackView)(nil)

func ZStack(views ...View) *zstackView {
	v := &zstackView{}
	v.viewOption = newViewOption(v, views...)
	return v
}

type zstackView struct {
	viewOption
}

func (v *zstackView) Body() SomeView {
	return v
}

func (v *zstackView) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	current := parent.calculateViewOption(v.viewOption)
	logs.Debugf("ZStackView: %+v", current)
	current.Draw(screen, func(screen *ebiten.Image) {
		current.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, &current)
			if v != nil {
				_ = v.draw(screen, current.CreateChild(v, current.X(), current.Y(), current.XX(), current.YY(), current.Width(), current.Height()))
			}
		})
	})
	return current.DrawnArea()
}

func (v *zstackView) bounds() (int, int) {
	return v.w, v.h
}
