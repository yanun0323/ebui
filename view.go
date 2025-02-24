package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

type View interface {
	Body() SomeView
}

// layoutFunc: 用於設置 View 的位置及大小，並回傳實際佔用的空間
//
//	start: 給這個 View 的起始座標
//	flexFrameSize: 給這個 View 的內部邊界彈性大小
//	bounds: 回傳實際佔用的空間(包含 padding 的最外圍邊界)
type layoutFunc func(start CGPoint, flexFrameSize CGSize) (bounds CGRect)

// SomeView 是所有 View 的基礎介面
type SomeView interface {
	View

	// preload 回傳的 frameSize 是 View 用 Frame 設置的大小
	// preload 回傳的 padding 是 View 用 Padding 設置的 padding
	// layoutFn: 用於設置 View 的位置及大小，並回傳實際佔用的空間
	// 		start: 給這個 View 的起始座標
	// 		flexFrameSize: 給這個 View 的內部邊界彈性大小
	// 		bounds: 回傳實際佔用的空間(包含 padding 的最外圍邊界)
	preload() (frameSize flexibleCGSize, padding Inset, layoutFn layoutFunc)

	// draw 繪製 View
	draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions
	debugPrint(frame CGRect)

	userSetFrameSize() flexibleCGSize
	systemSetFrame() CGRect

	Debug(tag string) SomeView
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
	ScaleToFit(enable ...*Binding[bool]) SomeView
	KeepAspectRatio(enable ...*Binding[bool]) SomeView
}

// Frame: 不包含 Padding 的內部邊界
// Bounds: 包含 Padding 的外部邊界
