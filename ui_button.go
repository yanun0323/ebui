package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type buttonImpl struct {
	*ctx

	action    func()
	label     func() SomeView
	frame     CGRect
	isPressed bool
	cache     *ViewCache
}

func Button(action func(), label func() SomeView) SomeView {
	btn := &buttonImpl{
		action: action,
		label:  label,
		cache:  NewViewCache(),
	}
	btn.ctx = newViewContext(tagButton, btn)
	RegisterEventHandler(btn)
	return btn
}

func (b *buttonImpl) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	btnFrameSize, btnInset, btnLayoutFn := b.ctx.preload()
	labelFrameSize, labelInset, labelLayoutFn := b.label().preload()

	frameSize := btnFrameSize.MaxWH(labelFrameSize)
	frameInset := btnInset.MaxBounds(labelInset)
	frameSize = frameSize.Expand(frameInset)

	return frameSize, frameInset, func(start CGPoint, flexibleSize CGSize) CGRect {
		btnResult := btnLayoutFn(start, flexibleSize)
		labelResult := labelLayoutFn(btnResult.Start, btnResult.Size())
		return btnResult.MaxStartEnd(labelResult)
	}
}

func (b *buttonImpl) HandleInput(x, y float64, pressed bool) bool {
	if pressed && pt(x, y).In(b.frame) {
		b.action()
		return true
	}
	return false
}

func (b *buttonImpl) draw(screen *ebiten.Image, bounds ...CGRect) {
	b.ctx.draw(screen, bounds...)
	b.label().draw(screen, bounds...)
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
