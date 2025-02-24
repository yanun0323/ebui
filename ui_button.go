package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type buttonImpl struct {
	*ctx

	action      func()
	label       func() SomeView
	labelLoaded SomeView
	isPressed   bool
}

func Button(action func(), label func() SomeView) SomeView {
	btn := &buttonImpl{
		action: action,
		label:  label,
	}
	btn.ctx = newViewContext(btn)
	globalEventManager.RegisterHandler(btn)
	return btn
}

func (b *buttonImpl) preload() (flexibleCGSize, CGInset, layoutFunc) {
	b.labelLoaded = b.label()
	formulaStack := &formulaStack{
		types:    formulaZStack,
		stackCtx: b.ctx,
		children: []SomeView{b.labelLoaded},
	}

	return formulaStack.preload()
}

func (b *buttonImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	hook = append(hook, func(opt *ebiten.DrawImageOptions) {
		if b.isPressed {
			opt.ColorScale.ScaleAlpha(0.5)
		}
	})

	op := b.ctx.draw(screen, hook...)
	_ = b.labelLoaded.draw(screen, hook...)

	return op
}

// Button 的事件處理
func (b *buttonImpl) HandleTouchEvent(event touchEvent) bool {
	switch event.Phase {
	case touchPhaseBegan:
		if event.Position.In(b.labelLoaded.systemSetFrame()) {
			b.isPressed = true
			return true
		}
	case touchPhaseMoved:
		if b.isPressed {
			return true
		}
	case touchPhaseEnded, touchPhaseCancelled:
		if b.isPressed {
			b.isPressed = false
			if event.Position.In(b.labelLoaded.systemSetFrame()) {
				b.action()
			}
			return true
		}
	}
	return false
}

func (b *buttonImpl) HandleKeyEvent(event keyEvent) bool {
	return false
}
