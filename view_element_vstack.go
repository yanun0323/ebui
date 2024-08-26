package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implementation */
var _ SomeView = (*vstackView)(nil)

func VStack(views ...View) *vstackView {
	v := &vstackView{}
	v.view = newView(typeVStack, v, views...)
	return v
}

type vstackView struct {
	*view
}

func (v *vstackView) draw(screen *ebiten.Image) {
	cache := v.view.Copy()
	cache.Draw(screen, func(screen *ebiten.Image) {
		cache.IterateViewModifiers(func(vm viewModifier) {
			vv := vm(screen, cache)
			if vv != nil {
				vv.draw(screen)
			}
		})
	})
}
