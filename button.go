package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type buttonImpl struct {
	*viewContext

	action    func()
	label     func() SomeView
	frame     image.Rectangle
	isPressed bool
	cache     *ViewCache
}

func Button(action func(), label func() SomeView) SomeView {
	btn := &buttonImpl{
		action: action,
		label:  label,
		cache:  NewViewCache(),
	}
	btn.viewContext = NewViewContext(btn)
	RegisterEventHandler(btn)
	return btn
}

func (b *buttonImpl) layout(bounds image.Rectangle) image.Rectangle {
	return b.label().layout(bounds)
}

func (b *buttonImpl) HandleInput(x, y int, pressed bool) bool {
	if pressed && image.Pt(x, y).In(b.frame) {
		b.action()
		return true
	}
	return false
}

func (b *buttonImpl) draw(screen *ebiten.Image) {
	b.viewContext.draw(screen)
	b.viewContext.drawHelper(screen, b.frame, func(screen *ebiten.Image) {
		b.label().draw(screen)
	})
}

// Button 的事件處理
func (b *buttonImpl) HandleTouchEvent(event TouchEvent) bool {
	switch event.Phase {
	case TouchPhaseBegan:
		if event.Position.In(b.frame) {
			b.isPressed = true
			return true
		}
	case TouchPhaseMoved:
		if b.isPressed {
			b.isPressed = event.Position.In(b.frame)
			return true
		}
	case TouchPhaseEnded, TouchPhaseCancelled:
		if b.isPressed {
			b.isPressed = false
			if event.Position.In(b.frame) {
				b.action()
			}
			return true
		}
	}
	return false
}

func (b *buttonImpl) HandleKeyEvent(event KeyEvent) bool {
	return false
}
