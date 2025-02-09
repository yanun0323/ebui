package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...SomeView) SomeView {
	vs := &vstackImpl{
		children: views,
	}
	vs.viewContext = NewViewContext(vs)
	return vs
}

type vstackImpl struct {
	*viewContext

	children []SomeView
	frame    image.Rectangle
}

// VStack 的佈局邏輯
func (v *vstackImpl) layout(bounds image.Rectangle) image.Rectangle {
	// 先應用修飾器的佈局
	bounds = v.viewContext.layout(bounds)

	y := bounds.Min.Y
	maxWidth := 0

	// 先讓子視圖計算各自的大小
	for _, child := range v.children {
		childBounds := child.layout(image.Rect(
			bounds.Min.X,
			y,
			bounds.Max.X,
			bounds.Max.Y,
		))

		// 更新最大寬度
		if childBounds.Dx() > maxWidth {
			maxWidth = childBounds.Dx()
		}

		// 累加高度，添加間距
		y = childBounds.Max.Y + 8 // 添加固定間距
	}

	// 水平置中對齊所有子視圖
	y = bounds.Min.Y // 重置 y 座標
	for _, child := range v.children {
		childBounds := child.layout(image.Rect(
			bounds.Min.X+(bounds.Dx()-maxWidth)/2, // 水平置中
			y,
			bounds.Min.X+(bounds.Dx()+maxWidth)/2,
			bounds.Max.Y,
		))
		y = childBounds.Max.Y + 8 // 更新 y 座標，添加間距
	}

	v.frame = image.Rect(bounds.Min.X, bounds.Min.Y,
		bounds.Min.X+maxWidth, y-8) // 減去最後一個間距
	return v.frame
}

func (v *vstackImpl) draw(screen *ebiten.Image) {
	v.viewContext.draw(screen)
	v.viewContext.drawHelper(screen, v.frame, func(screen *ebiten.Image) {
		for _, child := range v.children {
			child.Body().draw(screen)
		}
	})
}
