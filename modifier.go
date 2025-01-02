package ebui

import "image/color"

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
	foregroundColor(color Binding[color.Color]) ViewBuilder
	frame(width, height float64) ViewBuilder
	padding(dp float64) ViewBuilder
}

// 實現基礎 modifier
type foregroundColorModifier struct {
	color Binding[color.Color]
	view  View
}

func (v ViewBuilder) foregroundColor(color Binding[color.Color]) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &foregroundColorModifier{
				color: color,
				view:  v.Build(),
			}
		},
	}
}
