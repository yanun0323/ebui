package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

var _ SomeView = &viewContext{}

const (
	_defaultRoundCorner float64 = 8.0
)

// viewContext 提供所有 View 共用的修飾器方法和狀態
type viewContext struct {
	owner SomeView
	size  image.Rectangle

	padding         *Binding[float64]
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
func NewViewContext(owner SomeView) *viewContext {
	return &viewContext{
		owner:           owner,
		size:            image.Rectangle{},
		padding:         NewBinding(0.0),
		backgroundColor: NewBinding[color.Color](nil),
		foregroundColor: NewBinding[color.Color](color.White),

		fontSize:          NewBinding(font.Body),
		fontWeight:        NewBinding(font.Normal),
		fontLineHeight:    NewBinding(0.0),
		fontLetterSpacing: NewBinding(0.0),
		fontAlignment:     NewBinding(font.AlignLeft),
	}
}

func (ctx *viewContext) Body() SomeView {
	return ctx.owner
}

func (ctx *viewContext) layout(bounds image.Rectangle) image.Rectangle {
	if padding := ctx.padding.Get(); padding > 0 {
		p := int(padding)
		return image.Rect(
			bounds.Min.X+p,
			bounds.Min.Y+p,
			bounds.Max.X-p,
			bounds.Max.Y-p,
		)
	}

	return bounds
}

func (ctx *viewContext) draw(screen *ebiten.Image) {
	if bgColor := ctx.backgroundColor.Get(); bgColor != nil {
		if radius := ctx.roundCorner.Get(); radius > 0 {
			drawRoundedRect(screen, screen.Bounds(), radius, bgColor)
		} else {
			screen.Fill(bgColor)
		}
	}
}

// drawHelper 輔助方法，處理所有修飾器的繪製
func (ctx *viewContext) drawHelper(screen *ebiten.Image, bounds image.Rectangle, drawContent func(screen *ebiten.Image)) {
	// 處理背景色
	if bgColor := ctx.backgroundColor.Get(); bgColor != nil {
		screen.SubImage(bounds).(*ebiten.Image).Fill(bgColor)
	}

	// 處理內容區域
	if padding := ctx.padding.Get(); padding > 0 {
		p := int(padding)
		bounds = image.Rect(
			bounds.Min.X+p,
			bounds.Min.Y+p,
			bounds.Max.X-p,
			bounds.Max.Y-p,
		)
	}

	// 繪製內容
	drawContent(screen)
}

// Padding 修飾器
func (ctx *viewContext) Padding(padding *Binding[float64]) SomeView {
	ctx.padding = padding
	return ctx
}

// BackgroundColor 修飾器
func (ctx *viewContext) BackgroundColor(color *Binding[color.Color]) SomeView {
	ctx.backgroundColor = color
	return ctx
}

// ForegroundColor 修飾器
func (ctx *viewContext) ForegroundColor(color *Binding[color.Color]) SomeView {
	ctx.foregroundColor = color
	return ctx
}

func (ctx *viewContext) FontSize(size *Binding[font.Size]) SomeView {
	ctx.fontSize = size
	return ctx
}

func (ctx *viewContext) FontWeight(weight *Binding[font.Weight]) SomeView {
	ctx.fontWeight = weight
	return ctx
}

func (ctx *viewContext) FontLineHeight(height *Binding[float64]) SomeView {
	ctx.fontLineHeight = height
	return ctx
}

func (ctx *viewContext) FontLetterSpacing(spacing *Binding[float64]) SomeView {
	ctx.fontLetterSpacing = spacing
	return ctx
}

func (ctx *viewContext) FontAlignment(alignment *Binding[font.Alignment]) SomeView {
	ctx.fontAlignment = alignment
	return ctx
}

func (ctx *viewContext) FontItalic(italic ...*Binding[bool]) SomeView {
	if len(italic) != 0 {
		ctx.fontItalic = italic[0]
	} else {
		ctx.fontItalic.Set(true)
	}

	return ctx
}

func (ctx *viewContext) RoundCorner(radius ...*Binding[float64]) SomeView {
	if len(radius) != 0 {
		ctx.roundCorner = radius[0]
	} else {
		ctx.roundCorner = NewBinding(_defaultRoundCorner)
	}
	return ctx
}
