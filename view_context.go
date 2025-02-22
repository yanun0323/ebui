package ebui

import (
	"fmt"
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
	_id             int64
	_tag            tag
	_systemSetFrame CGRect // 視圖的邊界，用 padding 修飾過的邊界
	_padding        *Binding[float64]

	frameWidth      *Binding[float64]
	frameHeight     *Binding[float64]
	backgroundColor *Binding[color.Color]
	foregroundColor *Binding[color.Color]

	fontSize          *Binding[font.Size]
	fontWeight        *Binding[font.Weight]
	fontLineHeight    *Binding[float64]
	fontLetterSpacing *Binding[float64]
	fontAlignment     *Binding[font.Alignment]
	fontItalic        *Binding[bool]
	roundCorner       *Binding[float64]
}

// 初始化 ViewContext
func newViewContext(viewType tag, owner SomeView) *ctx {
	return &ctx{
		_owner:          owner,
		_id:             getID(),
		_tag:            viewType,
		_systemSetFrame: rect(0, 0, 0, 0),
		_padding:        Bind(0.0),

		frameWidth:      Bind(Inf),
		frameHeight:     Bind(Inf),
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

func (ctx *ctx) id() int64 {
	if ctx._id == 0 {
		ctx._id = getID()
	}
	return ctx._id
}

func (ctx *ctx) types() string {
	return ctx._tag.String()
}

func (ctx *ctx) debug() string {
	return fmt.Sprintf("\t\x1b[%dm[%d]\x1b[0m \x1b[%dm[%s]\x1b[0m\t", 35, ctx.id(), 32, ctx.types())
}

func (ctx *ctx) userSetFrame() CGRect {
	return CGRect{ptZero, CGPoint{ctx.frameWidth.Get(), ctx.frameHeight.Get()}}
}

// systemSetFrame 回傳的是去除 padding 後的內部邊界
func (ctx *ctx) systemSetFrame() CGRect {
	return ctx._systemSetFrame
}

func (ctx *ctx) padding() Inset {
	p := ctx._padding.Get()
	return ins(p, p, p, p)
}

func (ctx *ctx) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	padding := ctx.padding()
	return ctx.userSetFrame().Size(), padding, func(start CGPoint, flexSize CGSize) CGRect {
		frame := ctx.userSetFrame()
		finalSize := frame.Size()
		if isInf(finalSize.Width) {
			finalSize.Width = flexSize.Width
		}

		if isInf(finalSize.Height) {
			finalSize.Height = flexSize.Height
		}

		frame = rect(start.X, start.Y, start.X+finalSize.Width, start.Y+finalSize.Height)
		ctx._systemSetFrame = frame.Shrink(padding)

		return frame
	}
}

func (ctx *ctx) draw(screen *ebiten.Image, bounds ...CGRect) {
	bgColor := ctx.backgroundColor.Get()
	if bgColor == nil {
		return
	}

	drawFrame := ctx.systemSetFrame()
	if len(bounds) != 0 {
		drawFrame = bounds[0]
	}

	if !drawFrame.drawable() {
		return
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(drawFrame.Start.X, drawFrame.Start.Y)

	if radius := ctx.roundCorner.Get(); radius > 0 {
		drawRoundedRect(screen, drawFrame, radius, bgColor, opt)
	} else {
		img := ebiten.NewImage(int(drawFrame.Dx()), int(drawFrame.Dy()))
		img.Fill(bgColor)
		screen.DrawImage(img, opt)

		// logf("[DRAW] %s frame: %+v\n", ctx._owner.debug(), drawFrame)
	}
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
func (ctx *ctx) Padding(padding *Binding[float64]) SomeView {
	ctx._padding = padding
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
