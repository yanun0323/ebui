package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

type View interface {
	Body() SomeView
}

type SomeView interface {
	someViewOption

	draw(screen *ebiten.Image, parent viewOption) image.Rectangle
	bounds() (int, int)
}

//go:generate domaingen -destination=view.option.go -package=ebui -name=viewOption
type someViewOption interface {
	View

	option() *viewOption

	Frame(w, h int) SomeView
	ForegroundColor(clr color.Color) SomeView
	BackgroundColor(clr color.Color) SomeView
	Padding(top, right, bottom, left int) SomeView
	FontSize(size font.Size) SomeView
	CornerRadius(radius ...int) SomeView
}
