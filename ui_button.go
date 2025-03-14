package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	buttonDefaultLabel = func(key string) func() SomeView {
		return func() SomeView {
			return Text(key).
				Padding(Bind(NewInset(5, 15, 5, 15))).
				BackgroundColor(AccentColor).
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

func (b *buttonImpl) preload(parent *viewCtxEnv, _ ...formulaType) (preloadData, layoutFunc) {
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
	hooks := make([]func(*ebiten.DrawImageOptions), 0, len(hook)+1)
	hooks = append(hooks, hook...)
	hooks = append(hooks, func(opt *ebiten.DrawImageOptions) {
		if b.isPressed {
			opt.ColorScale.ScaleAlpha(0.5)
		}
	})

	b.viewCtx.draw(screen, hooks...)
	b.labelLoaded.draw(screen, hooks...)
}

func (b *buttonImpl) HandleWheelEvent(event wheelEvent) {}

func (b *buttonImpl) HandleTouchEvent(event touchEvent) {
	if b.viewCtxEnv.disabled.Value() {
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
