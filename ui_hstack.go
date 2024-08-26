package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implementation */
var _ SomeView = (*hstackView)(nil)

func HStack(views ...View) *hstackView {
	v := &hstackView{}
	v.uiView = newUIView(typeHStack, v, views...)
	return v
}

type hstackView struct {
	*uiView
}

func (v *hstackView) draw(screen *ebiten.Image) {
	cache := v.uiView.Copy()
	cache.Draw(screen, func(screen *ebiten.Image) {
		cache.IterateViewModifiers(func(vm viewModifier) {
			vv := vm(screen, cache)
			if vv != nil {
				vv.draw(screen)
			}
		})
	})
}
