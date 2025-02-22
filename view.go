package ebui

import (
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

	id() int64

	// preload 回傳的 size 是 View 用 Frame 設置的大小
	// preload 回傳的 padding 是 View 的 padding
	// layout: 用於設置 View 的位置及大小，並回傳實際佔用的空間
	// 		start: 給這個 View 的起始座標
	// 		flexSize: 給這個 View 的彈性大小
	// 		return: 回傳實際佔用的空間(包含 padding 後的最外圍邊界)
	preload() (frameSize CGSize, padding Inset, layoutFn func(start CGPoint, flexSize CGSize) CGRect)

	// draw 繪製 View
	draw(screen *ebiten.Image, bounds ...CGRect)

	Frame(width *Binding[float64], height *Binding[float64]) SomeView
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
