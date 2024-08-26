package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implementation */
var _ SomeView = (*spacerView)(nil)

func Spacer() *spacerView {
	v := &spacerView{}
	v.view = newView(typeSpacer, v)
	return v
}

type spacerView struct {
	*view
}

func (v *spacerView) draw(screen *ebiten.Image) {}
