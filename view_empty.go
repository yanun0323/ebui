package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implement */
var _ SomeView = (*emptyView)(nil)

func Empty() *emptyView {
	v := &emptyView{}
	v.viewOption = newViewOption(v)
	return v
}

type emptyView struct {
	viewOption
}

func (v *emptyView) Body() SomeView {
	return v
}

func (v *emptyView) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	return image.Rectangle{}
}

func (v *emptyView) bounds() (int, int) {
	return 0, 0
}
