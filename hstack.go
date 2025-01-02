package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// HStack 水平排列
func HStack(views ...SomeView) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			children := make([]View, len(views))
			for i, v := range views {
				children[i] = v.Build()
			}
			return &hStackImpl{children: children}
		},
	}
}

type hStackImpl struct {
	children []View
	frame    image.Rectangle
}

// HStack 的佈局邏輯
func (h *hStackImpl) Layout(bounds image.Rectangle) image.Rectangle {
	x := bounds.Min.X
	maxHeight := 0

	// 先讓子視圖計算各自的大小
	for _, child := range h.children {
		childBounds := child.Layout(image.Rect(
			x,
			bounds.Min.Y,
			bounds.Max.X,
			bounds.Max.Y,
		))

		// 更新最大高度
		if childBounds.Dy() > maxHeight {
			maxHeight = childBounds.Dy()
		}

		// 累加寬度
		x += childBounds.Dx()
	}

	return image.Rect(bounds.Min.X, bounds.Min.Y,
		x, bounds.Min.Y+maxHeight)
}

func (h *hStackImpl) Draw(screen *ebiten.Image) {
	for _, child := range h.children {
		child.Draw(screen)
	}
}

func (h *hStackImpl) Build() View {
	return h
}
