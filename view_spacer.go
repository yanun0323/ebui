package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

/* Check Interface Implement */
var _ SomeView = (*spacerView)(nil)

func Spacer() *spacerView {
	v := &spacerView{}
	v.viewOption = newViewOption(v)
	return v
}

type spacerView struct {
	viewOption
}

func (v *spacerView) Body() SomeView {
	return v
}

func (v *spacerView) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	current := parent.calculateViewOption(v.viewOption)
	return current.DrawnArea()
}

func (v *spacerView) bounds() (int, int) {
	return v.w, v.h
}
