package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implementation */
var _ SomeView = (*vstackView)(nil)

func VStack(views ...View) *vstackView {
	v := &vstackView{}
	v.uiView = newUIView(typesVStack, v, views...)
	return v
}

type vstackView struct {
	*uiView
}

func (v *vstackView) draw(screen *ebiten.Image) {
	cache := v.uiView.Copy()
	cache.Draw(screen, func(screen *ebiten.Image) {
		cache.ApplyViewModifiers(screen)
	})
}
