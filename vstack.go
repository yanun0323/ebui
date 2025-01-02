package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...SomeView) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			children := make([]View, len(views))
			for i, v := range views {
				children[i] = v.Build()
			}

			return &vstackImpl{
				children: children,
			}
		},
	}
}

type vstackImpl struct {
	children []View
	frame    image.Rectangle
}

// VStack 的佈局邏輯
func (v *vstackImpl) Layout(bounds image.Rectangle) image.Rectangle {
	y := bounds.Min.Y
	maxWidth := 0

	// 先讓子視圖計算各自的大小
	for _, child := range v.children {
		childBounds := child.Layout(image.Rect(
			bounds.Min.X,
			y,
			bounds.Max.X,
			bounds.Max.Y,
		))

		// 更新最大寬度
		if childBounds.Dx() > maxWidth {
			maxWidth = childBounds.Dx()
		}

		// 累加高度
		y += childBounds.Dy()
	}

	return image.Rect(bounds.Min.X, bounds.Min.Y,
		bounds.Min.X+maxWidth, y)
}

func (v *vstackImpl) Draw(screen *ebiten.Image) {
	for _, child := range v.children {
		child.Draw(screen)
	}
}

func (v *vstackImpl) Build() View {
	return v
}
