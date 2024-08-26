package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implementation */
var _ SomeView = (*zstackView)(nil)

func ZStack(views ...View) *zstackView {
	v := &zstackView{}
	v.uiView = newUIView(typeZStack, v, views...)
	return v
}

type zstackView struct {
	*uiView
}

func (v *zstackView) draw(screen *ebiten.Image) {
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
