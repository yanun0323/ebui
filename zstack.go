package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type zStackImpl struct {
	children []View
	frame    image.Rectangle
}

func ZStack(views ...SomeView) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			children := make([]View, len(views))
			for i, v := range views {
				children[i] = v.Build()
			}
			return &zStackImpl{children: children}
		},
	}
}

// 實現 View 介面
func (z *zStackImpl) Build() View {
	return z
}

func (z *zStackImpl) Layout(bounds image.Rectangle) image.Rectangle {
	z.frame = bounds
	maxWidth, maxHeight := 0, 0

	for _, child := range z.children {
		childBounds := child.Layout(bounds)
		if childBounds.Dx() > maxWidth {
			maxWidth = childBounds.Dx()
		}
		if childBounds.Dy() > maxHeight {
			maxHeight = childBounds.Dy()
		}
	}

	return image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X+maxWidth,
		bounds.Min.Y+maxHeight,
	)
}

func (z *zStackImpl) Draw(screen *ebiten.Image) {
	for _, child := range z.children {
		child.Draw(screen)
	}
}
