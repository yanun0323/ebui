package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func ZStack(views ...SomeView) SomeView {
	zs := &zstackImpl{
		children: views,
	}
	zs.viewContext = NewViewContext(zs)
	return zs
}

type zstackImpl struct {
	*viewContext

	children []SomeView
	frame    image.Rectangle
}

func (z *zstackImpl) layout(bounds image.Rectangle) image.Rectangle {
	// 先應用修飾器的佈局
	bounds = z.viewContext.layout(bounds)

	maxWidth, maxHeight := 0, 0

	for _, child := range z.children {
		childBounds := child.Body().layout(bounds)
		if childBounds.Dx() > maxWidth {
			maxWidth = childBounds.Dx()
		}
		if childBounds.Dy() > maxHeight {
			maxHeight = childBounds.Dy()
		}
	}

	z.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X+maxWidth,
		bounds.Min.Y+maxHeight,
	)
	return z.frame
}

func (z *zstackImpl) draw(screen *ebiten.Image) {
	z.viewContext.draw(screen)
	z.viewContext.drawHelper(screen, z.frame, func(screen *ebiten.Image) {
		for _, child := range z.children {
			child.Body().draw(screen)
		}
	})
}
