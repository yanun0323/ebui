package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type View interface {
	SomeView

	Layout(bounds image.Rectangle) image.Rectangle
	Draw(screen *ebiten.Image)
}

// SomeView 是所有 View 的基礎介面
type SomeView interface {
	Build() View
}

// ViewBuilder 用來構建具體的 View
type ViewBuilder struct {
	build func() View
}

func (v ViewBuilder) Build() View {
	return v.build()
}

func (v ViewBuilder) WithPadding(padding float64) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			inner := v.Build()
			return &paddingModifier{
				view:    inner,
				padding: padding,
			}
		},
	}
}

type paddingModifier struct {
	view    View
	padding float64
}

func (p *paddingModifier) Build() View {
	return p
}

func (p *paddingModifier) Layout(bounds image.Rectangle) image.Rectangle {
	padding := int(p.padding)
	innerBounds := image.Rect(
		bounds.Min.X+padding,
		bounds.Min.Y+padding,
		bounds.Max.X-padding,
		bounds.Max.Y-padding,
	)
	return p.view.Layout(innerBounds)
}

func (p *paddingModifier) Draw(screen *ebiten.Image) {
	p.view.Draw(screen)
}

// 背景色修飾器
type backgroundColorModifier struct {
	next  View
	color color.Color
}

func (m *backgroundColorModifier) Layout(bounds image.Rectangle) image.Rectangle {
	return m.next.Layout(bounds)
}

func (m *backgroundColorModifier) Draw(screen *ebiten.Image) {
	// 先繪製背景色
	screen.SubImage(m.next.Layout(screen.Bounds())).(*ebiten.Image).Fill(m.color)
	// 再繪製內容
	m.next.Draw(screen)
}

func (m *backgroundColorModifier) Build() View {
	return &backgroundColorModifier{
		next:  m.next.Build(),
		color: m.color,
	}
}

// ViewBuilder 的擴展方法
func (v ViewBuilder) BackgroundColor(color color.Color) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &backgroundColorModifier{
				next:  v.Build(),
				color: color,
			}
		},
	}
}
