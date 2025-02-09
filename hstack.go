package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// HStack 水平排列
func HStack(views ...SomeView) SomeView {
	hs := &hstackImpl{
		children: views,
	}
	hs.viewContext = NewViewContext(hs)
	return hs
}

type hstackImpl struct {
	*viewContext

	children []SomeView
	frame    image.Rectangle
}

// HStack 的佈局邏輯
func (h *hstackImpl) layout(bounds image.Rectangle) image.Rectangle {
	// 先應用修飾器的佈局
	bounds = h.viewContext.layout(bounds)

	x := bounds.Min.X
	maxHeight := 0

	// 先讓子視圖計算各自的大小
	for _, child := range h.children {
		childBounds := child.layout(image.Rect(
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

	h.frame = image.Rect(bounds.Min.X, bounds.Min.Y,
		x, bounds.Min.Y+maxHeight)
	return h.frame
}

func (h *hstackImpl) draw(screen *ebiten.Image) {
	h.viewContext.draw(screen)
	h.viewContext.drawHelper(screen, h.frame, func(screen *ebiten.Image) {
		for _, child := range h.children {
			child.draw(screen)
		}
	})
}
