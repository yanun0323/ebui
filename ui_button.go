package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	buttonDefaultLabel = func(key string) func() SomeView {
		return func() SomeView {
			return Text(key).
				Padding(Bind(NewInset(5, 15, 5, 15))).
				BackgroundColor(Bind(AccentColor)).
				RoundCorner(Bind(15.0))
		}
	}
)

type buttonImpl struct {
	*viewCtx

	action      func()
	label       func() SomeView
	labelLoaded SomeView
	isPressed   bool
}

func Button(key string, action func(), label ...func() SomeView) SomeView {
	lb := buttonDefaultLabel(key)
	if len(label) != 0 && label[0] != nil {
		lb = label[0]
	}

	btn := &buttonImpl{
		action: action,
		label:  lb,
	}
	btn.viewCtx = newViewContext(btn)
	globalEventManager.RegisterHandler(btn)
	return btn
}

func (b *buttonImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	b.labelLoaded = b.label()
	if b.labelLoaded == nil {
		panic("empty view from button label")
	}

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
func (b *buttonImpl) HandleTouchEvent(event touchEvent) {
	if b.viewCtxEnv.disabled.Get() {
		b.isPressed = false
		return
	}

	switch event.Phase {
	case touchPhaseBegan:
		if event.Position.In(b.labelLoaded.systemSetBounds()) {
			b.isPressed = true
		}
	case touchPhaseMoved:
	case touchPhaseEnded, touchPhaseCancelled:
		if b.isPressed {
			b.isPressed = false
			if event.Position.In(b.labelLoaded.systemSetBounds()) {
				if b.action != nil {
					b.action()
				}
			}
		}
	}
}

func (b *buttonImpl) HandleKeyEvent(event keyEvent) {}

func (b *buttonImpl) HandleInputEvent(event inputEvent) {}
