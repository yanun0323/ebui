package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonStyle struct {
	backgroundColor color.Color
	pressedColor    color.Color
	cornerRadius    float64
	borderWidth     float64
	borderColor     color.Color
	padding         float64
}

type buttonImpl struct {
	action   func()
	label    View
	frame    image.Rectangle
	style    ButtonStyle
	isPressed bool
	cache    *ViewCache
}

func Button(action func(), label SomeView) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &buttonImpl{
				action: action,
				label:  label.Build(),
				style: ButtonStyle{
					backgroundColor: color.RGBA{0x4C, 0xAF, 0x50, 0xFF}, // Material Design Green
					pressedColor:    color.RGBA{0x38, 0x8E, 0x3C, 0xFF},
					cornerRadius:    8,
					padding:         16,
				},
				cache: NewViewCache(),
			}
		},
	}
}

func (b *buttonImpl) Layout(bounds image.Rectangle) image.Rectangle {
	labelBounds := b.label.Layout(bounds)
	padding := int(b.style.padding)
	
	b.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X + labelBounds.Dx() + padding*2,
		bounds.Min.Y + labelBounds.Dy() + padding*2,
	)
	
	// 調整標籤位置到按鈕中心
	labelX := b.frame.Min.X + (b.frame.Dx()-labelBounds.Dx())/2
	labelY := b.frame.Min.Y + (b.frame.Dy()-labelBounds.Dy())/2
	b.label.Layout(image.Rect(labelX, labelY, labelX+labelBounds.Dx(), labelY+labelBounds.Dy()))
	
	return b.frame
}

func (b *buttonImpl) Draw(screen *ebiten.Image) {
	// 繪製背景
	bgColor := b.style.backgroundColor
	if b.isPressed {
		bgColor = b.style.pressedColor
	}
	
	// 使用圓角矩形
	drawRoundedRect(screen, b.frame, b.style.cornerRadius, bgColor)
	
	// 繪製邊框
	if b.style.borderWidth > 0 {
		drawRoundedRectBorder(screen, b.frame, b.style.cornerRadius, b.style.borderWidth, b.style.borderColor)
	}
	
	// 繪製標籤
	b.label.Draw(screen)
}

func (b *buttonImpl) HandleInput(x, y int, pressed bool) bool {
	if pressed && b.frame.Contains(image.Pt(x, y)) {
		b.action()
		return true
	}
	return false
}

// Button 的事件處理
func (b *buttonImpl) HandleTouchEvent(event TouchEvent) bool {
	if event.Phase == TouchPhasePressed {
		if b.frame.Contains(event.Position) {
			b.action()  // 執行點擊動作
			return true // 事件已處理
		}
	}
	return false // 事件未處理，繼續傳遞
}
