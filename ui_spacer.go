package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implementation */
var _ SomeView = (*spacerView)(nil)

func Spacer() *spacerView {
	v := &spacerView{}
	v.uiView = newUIView(typeSpacer, v)
	return v
}

type spacerView struct {
	*uiView
}

func (v *spacerView) draw(screen *ebiten.Image) {}
