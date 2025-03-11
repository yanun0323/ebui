package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type buttonImpl struct {
	*viewCtx

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
	btn.viewCtx = newViewContext(btn)
	globalEventManager.RegisterHandler(btn)
	return btn
}

func (b *buttonImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	b.labelLoaded = b.label()
	formulaStack := &formulaStack{
		types:    formulaZStack,
		stackCtx: b.viewCtx,
		children: []SomeView{b.labelLoaded},
	}

	return formulaStack.preload(parent)
}

func (b *buttonImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	hook = append(hook, func(opt *ebiten.DrawImageOptions) {
		if b.isPressed {
			opt.ColorScale.ScaleAlpha(0.5)
		}
	})

	b.viewCtx.draw(screen, hook...)
	b.labelLoaded.draw(screen, hook...)
}

// Button 的事件處理
func (b *buttonImpl) HandleTouchEvent(event touchEvent) bool {
	switch event.Phase {
	case touchPhaseBegan:
		if event.Position.In(b.labelLoaded.systemSetBounds()) {
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
			if event.Position.In(b.labelLoaded.systemSetBounds()) {
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
