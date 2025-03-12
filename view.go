package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
	layout "github.com/yanun0323/ebui/layout"
)

type View interface {
	Body() SomeView
}

type ViewModifier interface {
	Body(SomeView) SomeView
}

// SomeView represents the base of all View
//
// Frame: internal bounds (exclude padding, border ...)
// Bounds: external bounds (include padding, border ...)
type SomeView interface {
	View

	ctx() *viewCtx
	// preload 設置環境變量，並回傳 View 的佈局資訊
	preload(parent *viewCtxEnv) (preloadData, layoutFunc)

	drawOption(rect CGRect, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions
	// draw 繪製 View
	draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions))

	debugPrint(prefix string, format string, a ...any)

	setEnv(env *viewCtxEnv)

	count() int
	isSpacer() bool

	userSetFrameSize() CGSize
	systemSetFrame() CGRect
	systemSetBounds() CGRect

	DebugPrint(tag string) SomeView
	Frame(*Binding[CGSize]) SomeView
	Padding(inset *Binding[CGInset]) SomeView
	ForegroundColor(color *Binding[CGColor]) SomeView
	BackgroundColor(color *Binding[CGColor]) SomeView
	FontSize(size *Binding[font.Size]) SomeView
	FontWeight(weight *Binding[font.Weight]) SomeView
	FontLineHeight(height *Binding[float64]) SomeView
	FontLetterSpacing(spacing *Binding[float64]) SomeView
	FontAlignment(alignment *Binding[font.TextAlign]) SomeView
	FontItalic(italic ...*Binding[bool]) SomeView
	RoundCorner(radius ...*Binding[float64]) SomeView
	ScaleToFit(enable ...*Binding[bool]) SomeView
	KeepAspectRatio(enable ...*Binding[bool]) SomeView
	Border(border *Binding[CGInset], color ...*Binding[CGColor]) SomeView
	Opacity(opacity *Binding[float64]) SomeView
	Modifier(ViewModifier) SomeView
	Modify(with func(SomeView) SomeView) SomeView
	Align(alignment *Binding[layout.Align]) SomeView
	Center() SomeView
	Debug() SomeView
}

type alignFunc func(offset CGPoint)

// layoutFunc: 用於設置 View 的位置及大小，並回傳實際佔用的空間
//
//	start: 給這個 View 的起始座標
//	flexBoundsSize: 給這個 View 的外部邊界彈性大小
//	bounds: 回傳實際佔用的空間(包含 padding 的最外圍邊界)
type layoutFunc func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect, alignFunc alignFunc)

func newPreloadData(frameSize CGSize, padding CGInset, border CGInset) preloadData {
	return preloadData{
		FrameSize:   frameSize,
		IsInfWidth:  frameSize.IsInfWidth(),
		IsInfHeight: frameSize.IsInfHeight(),
		Padding:     padding,
		Border:      border,
	}
}

type preloadData struct {
	FrameSize   CGSize // 實際的內部邊界尺寸
	IsInfWidth  bool
	IsInfHeight bool
	Padding     CGInset // padding 是 View 用 Padding 設置的 padding
	Border      CGInset // border 是 View 用 Border 設置的 border
}

func (p *preloadData) BoundsSize() CGSize {
	return p.FrameSize.Expand(p.Padding).Expand(p.Border)
}
