package ebui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

var _ SomeView = &ctx{}

var Inf = math.Inf(1)
var isInf = func(f float64) bool {
	return math.IsInf(f, 1)
}

const (
	_defaultRoundCorner float64 = 8.0
)

// ctx 提供所有 View 共用的修飾器方法和狀態
type ctx struct {
	_owner          SomeView
	_tag            string
	_systemSetFrame CGRect // 不包含 Padding 的內部邊界

	frameWidth      *Binding[float64]
	frameHeight     *Binding[float64]
	padding         *Binding[Inset]
	backgroundColor *Binding[color.Color]
	foregroundColor *Binding[color.Color]

	fontSize          *Binding[font.Size]
	fontWeight        *Binding[font.Weight]
	fontLineHeight    *Binding[float64]
	fontLetterSpacing *Binding[float64]
	fontAlignment     *Binding[font.Alignment]
	fontItalic        *Binding[bool]
	roundCorner       *Binding[float64]

	scaleToFit      *Binding[bool]
	keepAspectRatio *Binding[bool]
}

// 初始化 ViewContext
func newViewContext(owner SomeView) *ctx {
	return &ctx{
		_owner: owner,

		frameWidth:      Bind(Inf),
		frameHeight:     Bind(Inf),
		_systemSetFrame: CGRect{},
		padding:         Bind(Inset{}),
		backgroundColor: Bind[color.Color](nil),
		foregroundColor: Bind[color.Color](color.White),

		fontSize:          Bind(font.Body),
		fontWeight:        Bind(font.Normal),
		fontLineHeight:    Bind(0.0),
		fontLetterSpacing: Bind(0.0),
		fontAlignment:     Bind(font.AlignLeft),
	}
}

func (ctx *ctx) Body() SomeView {
	return ctx._owner
}

func (ctx *ctx) userSetFrameSize() flexibleCGSize {
	return newFlexibleCGSize(ctx.frameWidth.Get(), ctx.frameHeight.Get())
}

// systemSetFrame 回傳的是內部邊界
func (ctx *ctx) systemSetFrame() CGRect {
	return ctx._systemSetFrame
}

func (ctx *ctx) debugPrint(frame CGRect) {
	if len(ctx._tag) != 0 {
		logf("\x1b[35m[%s]\x1b[0m\tStart(%4.f, %4.f)\tEnd(%4.f, %4.f)\tSize(%4.f, %4.f)",
			ctx._tag,
			frame.Start.X, frame.Start.Y,
			frame.End.X, frame.End.Y,
			frame.Dx(), frame.Dy(),
		)
	}
}

func (ctx *ctx) preload() (flexibleCGSize, Inset, layoutFunc) {
	padding := ctx.padding.Get()
	userSetFrameSize := ctx._owner.userSetFrameSize()
	return userSetFrameSize, padding, func(start CGPoint, flexFrameSize CGSize) CGRect {
		finalFrame := CGRect{start, start.Add(flexFrameSize.ToCGPoint())}
		finalFrameSize := userSetFrameSize
		if !finalFrameSize.IsInfX {
			finalFrame.End.X = start.X + finalFrameSize.Frame.Width
		}

		if !finalFrameSize.IsInfY {
			finalFrame.End.Y = start.Y + finalFrameSize.Frame.Height
		}

		ctx._systemSetFrame = rect(
			finalFrame.Start.X+padding.Left,
			finalFrame.Start.Y+padding.Top,
			finalFrame.End.X+padding.Left,
			finalFrame.End.Y+padding.Top,
		)

		result := finalFrame.Expand(padding)
		ctx.debugPrint(result)
		return result
	}
}

func (ctx *ctx) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {

	drawFrame := ctx._owner.systemSetFrame()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(drawFrame.Start.X, drawFrame.Start.Y)
	for _, h := range hook {
		h(opt)
	}

	bgColor := ctx.backgroundColor.Get()
	if bgColor == nil {
		return opt
	}

	if !drawFrame.drawable() {
		return opt
	}

	if radius := ctx.roundCorner.Get(); radius > 0 {
		drawRoundedRect(screen, drawFrame, radius, bgColor, opt)
	} else {
		img := ebiten.NewImage(int(drawFrame.Dx()), int(drawFrame.Dy()))
		img.Fill(bgColor)
		screen.DrawImage(img, opt)
	}

	return opt
}

// Frame 修飾器
func (ctx *ctx) Frame(width *Binding[float64], height *Binding[float64]) SomeView {
	if width == nil {
		width = Bind(Inf)
	}

	if height == nil {
		height = Bind(Inf)
	}

	ctx.frameWidth = width
	ctx.frameHeight = height

	return ctx._owner
}

// Padding 修飾器
func (ctx *ctx) Padding(padding *Binding[Inset]) SomeView {
	ctx.padding = padding
	return ctx._owner
}

// BackgroundColor 修飾器
func (ctx *ctx) BackgroundColor(color *Binding[color.Color]) SomeView {
	ctx.backgroundColor = color
	return ctx._owner
}

// ForegroundColor 修飾器
func (ctx *ctx) ForegroundColor(color *Binding[color.Color]) SomeView {
	ctx.foregroundColor = color
	return ctx._owner
}

func (ctx *ctx) FontSize(size *Binding[font.Size]) SomeView {
	ctx.fontSize = size
	return ctx._owner
}

func (ctx *ctx) FontWeight(weight *Binding[font.Weight]) SomeView {
	ctx.fontWeight = weight
	return ctx._owner
}

func (ctx *ctx) FontLineHeight(height *Binding[float64]) SomeView {
	ctx.fontLineHeight = height
	return ctx._owner
}

func (ctx *ctx) FontLetterSpacing(spacing *Binding[float64]) SomeView {
	ctx.fontLetterSpacing = spacing
	return ctx._owner
}

func (ctx *ctx) FontAlignment(alignment *Binding[font.Alignment]) SomeView {
	ctx.fontAlignment = alignment
	return ctx._owner
}

func (ctx *ctx) FontItalic(italic ...*Binding[bool]) SomeView {
	if len(italic) != 0 {
		ctx.fontItalic = italic[0]
	} else {
		ctx.fontItalic.Set(true)
	}

	return ctx._owner
}

func (ctx *ctx) RoundCorner(radius ...*Binding[float64]) SomeView {
	if len(radius) != 0 {
		ctx.roundCorner = radius[0]
	} else {
		ctx.roundCorner = Bind(_defaultRoundCorner)
	}
	return ctx._owner
}

func (ctx *ctx) Debug(tag string) SomeView {
	ctx._tag = tag
	return ctx._owner
}

func (ctx *ctx) ScaleToFit(enable ...*Binding[bool]) SomeView {
	if len(enable) != 0 {
		ctx.scaleToFit = enable[0]
	} else {
		ctx.scaleToFit = Bind(true)
	}

	return ctx._owner
}

func (ctx *ctx) KeepAspectRatio(enable ...*Binding[bool]) SomeView {
	if len(enable) != 0 {
		ctx.keepAspectRatio = enable[0]
	} else {
		ctx.keepAspectRatio = Bind(true)
	}

	return ctx._owner
}
