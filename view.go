package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

type View interface {
	Body() SomeView
}

// layoutFunc: 用於設置 View 的位置及大小，並回傳實際佔用的空間
//
//	start: 給這個 View 的起始座標
//	flexBoundsSize: 給這個 View 的外部邊界彈性大小
//	bounds: 回傳實際佔用的空間(包含 padding 的最外圍邊界)
type layoutFunc func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect)

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

// SomeView 是所有 View 的基礎介面
type SomeView interface {
	View

	// preload 設置環境變量，並回傳 View 的佈局資訊
	preload(parent *viewCtxEnv) (preloadData, layoutFunc)

	// draw 繪製 View
	draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions))
	debugPrint(stage string, frame CGRect, flexFrameSize CGSize, sData preloadData)

	setEnv(env *viewCtxEnv)

	count() int
	isSpacer() bool

	userSetFrameSize() CGSize
	systemSetFrame() CGRect
	systemSetBounds() CGRect

	Debug(tag string) SomeView
	Frame(*Binding[CGSize]) SomeView
	Padding(inset *Binding[CGInset]) SomeView
	ForegroundColor(color *Binding[AnyColor]) SomeView
	BackgroundColor(color *Binding[AnyColor]) SomeView
	FontSize(size *Binding[font.Size]) SomeView
	FontWeight(weight *Binding[font.Weight]) SomeView
	FontLineHeight(height *Binding[float64]) SomeView
	FontLetterSpacing(spacing *Binding[float64]) SomeView
	FontAlignment(alignment *Binding[font.Alignment]) SomeView
	FontItalic(italic ...*Binding[bool]) SomeView
	RoundCorner(radius ...*Binding[float64]) SomeView
	ScaleToFit(enable ...*Binding[bool]) SomeView
	KeepAspectRatio(enable ...*Binding[bool]) SomeView
	Border(border *Binding[CGInset], color ...*Binding[AnyColor]) SomeView
	Opacity(opacity *Binding[float64]) SomeView
}

// Frame: 不包含 Padding 的內部邊界
// Bounds: 包含 Padding 的外部邊界
