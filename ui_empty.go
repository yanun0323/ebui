package ebui

import "github.com/hajimehoshi/ebiten/v2"

/* Check Interface Implementation */
var _ SomeView = (*emptyView)(nil)

func EmptyView() *emptyView {
	v := &emptyView{}
	v.uiView = newView(typesEmpty, v)
	return v
}

type emptyView struct {
	*uiView
}

func (p *emptyView) getSize() size {
	return size{}
}

func (p *emptyView) Frame(w, h int) SomeView {
	return p
}

func (v *emptyView) draw(screen *ebiten.Image) {
	v.drawModifiers(screen)
}
