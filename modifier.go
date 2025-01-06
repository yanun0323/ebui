package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ViewModifier interface {
	Modify(View) View
}

type ModifiedView struct {
	source   ViewBuilder
	modifier ViewModifier
}

func (v *ModifiedView) Build() View {
	return v.modifier.Modify(v.source.Build())
}

// 為所有 View 添加 modifier 支援
type ViewModifiable interface {
	foregroundColor(color *Binding[color.Color]) ViewBuilder
	frame(width, height float64) ViewBuilder
	padding(dp float64) ViewBuilder
}

// 實現基礎 modifier
type foregroundColorModifier struct {
	color *Binding[color.Color]
	view  View
}

// 實現 View 介面
func (m *foregroundColorModifier) Build() View {
	return m
}

func (m *foregroundColorModifier) Layout(bounds image.Rectangle) image.Rectangle {
	return m.view.Layout(bounds)
}

func (m *foregroundColorModifier) Draw(screen *ebiten.Image) {
	// 保存原始顏色
	m.view.Draw(screen)
}

func (v ViewBuilder) foregroundColor(color *Binding[color.Color]) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &foregroundColorModifier{
				color: color,
				view:  v.Build(),
			}
		},
	}
}
