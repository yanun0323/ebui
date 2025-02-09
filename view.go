package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

/*
	Old:

type View interface {
	SomeView

	Layout(bounds image.Rectangle) image.Rectangle
	Draw(screen *ebiten.Image)
}

type SomeView interface {
	Build() View
}

*/

type View interface {
	Body() SomeView
}

// SomeView 是所有 View 的基礎介面
type SomeView interface {
	View

	layout(bounds image.Rectangle) image.Rectangle
	draw(screen *ebiten.Image)

	Padding(padding *Binding[float64]) SomeView
	BackgroundColor(color *Binding[color.Color]) SomeView
	ForegroundColor(color *Binding[color.Color]) SomeView
	FontSize(size *Binding[font.Size]) SomeView
	FontWeight(weight *Binding[font.Weight]) SomeView
	FontLineHeight(height *Binding[float64]) SomeView
	FontLetterSpacing(spacing *Binding[float64]) SomeView
	FontAlignment(alignment *Binding[font.Alignment]) SomeView
	FontItalic(italic ...*Binding[bool]) SomeView
	RoundCorner(radius ...*Binding[float64]) SomeView
}
