package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type View interface {
	Body() SomeView
}

type SomeView interface {
	internalView

	Frame(Binding[Size]) SomeView
	Offset(Binding[Point]) SomeView
	ForegroundColor(Binding[color.Color]) SomeView
	BackgroundColor(Binding[color.Color]) SomeView
	// Opacity(Binding[float64]) SomeView
	// Disable(...Binding[bool]) SomeView

	// Padding makes view have padding.
	//
	// # []int slice length:
	// 	 - [0] default (all)
	// 	 - [1] all
	// 	 - [2] vertical, horizontal
	// 	 - [4] top, right, bottom, left
	// Padding(...Binding[int]) SomeView

	// FontSize(Binding[font.Size]) SomeView
	// FontWeight(Binding[font.Weight]) SomeView
	// LineSpacing(Binding[float64]) SomeView
	// Kerning(Binding[int]) SomeView
	// Italic(...Binding[bool]) SomeView

	// CornerRadius(radius ...int) SomeView

	// Resizable() SomeView
	// AspectRatio(ratio ...float64) SomeView

	// Border(clr color.Color, width ...int) SomeView
}

type internalView interface {
	id() identity
	bounds() (min, current, max Size)
	update(container Size)
	getRenderImage() *ebiten.Image
	drawable() bool
	draw(screen *ebiten.Image)
}
