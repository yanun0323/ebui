package ebui

import "github.com/hajimehoshi/ebiten/v2"

const (
	_defaultToggleSize    = 30
	_defaultTogglePadding = 2
)

var (
	defaultToggleOnColor            = NewColor(239, 239, 239, 255)
	defaultToggleOffColor           = NewColor(239, 239, 239, 255)
	defaultToggleOnBackgroundColor  = NewColor(64, 191, 64, 255)
	defaultToggleOffBackgroundColor = NewColor(64, 64, 64, 128)
	toggleDefaultLabel              = func() func(bool) SomeView {
		return func(enabled bool) SomeView {
			return HStack(
				If(enabled, Spacer(), EmptyView()),
				Circle().
					Frame(Const(NewSize(_defaultToggleSize, _defaultToggleSize))).
					BackgroundColor(BindFunc(func() AnyColor {
						if enabled {
							return defaultToggleOnColor
						}
						return defaultToggleOffColor
					}, func(AnyColor) {})).
					Padding(Const(NewInset(_defaultTogglePadding, _defaultTogglePadding, _defaultTogglePadding, _defaultTogglePadding))),
				If(!enabled, Spacer(), EmptyView()),
			).
				Frame(Const(NewSize(60, _defaultToggleSize+_defaultTogglePadding*2))).
				BackgroundColor(BindFunc(func() AnyColor {
					if enabled {
						return defaultToggleOnBackgroundColor
					}
					return defaultToggleOffBackgroundColor
				}, func(AnyColor) {})).
				RoundCorner(Const(float64(_defaultToggleSize / 2))).
				Padding(Const(NewInset(5, 5, 5, 5)))
		}
	}
)

type toggleImpl struct {
	*viewCtx

	label       func(bool) SomeView
	labelLoaded SomeView
	enabled     *Binding[bool]
	isPressed   bool
}

func Toggle(enabled *Binding[bool], label ...func(bool) SomeView) SomeView {
	lb := toggleDefaultLabel()
	if len(label) != 0 && label[0] != nil {
		lb = label[0]
	}

	t := &toggleImpl{
		enabled: enabled,
		label:   lb,
	}
	t.viewCtx = newViewContext(t)
	globalEventManager.RegisterHandler(t)
	return t
}

func (b *toggleImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	b.labelLoaded = b.label(b.enabled.Get())
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

func (b *toggleImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	hook = append(hook, func(opt *ebiten.DrawImageOptions) {
		if b.isPressed {
			opt.ColorScale.ScaleAlpha(0.5)
		}
	})

	b.viewCtx.draw(screen, hook...)
	b.labelLoaded.draw(screen, hook...)
}

func (t *toggleImpl) HandleTouchEvent(event touchEvent) {
	if t.viewCtxEnv.disabled.Get() {
		return
	}

	switch event.Phase {
	case touchPhaseBegan:
		if event.Position.In(t.labelLoaded.systemSetBounds()) {
			t.isPressed = true
		}
	case touchPhaseEnded, touchPhaseCancelled:
		if t.isPressed {
			t.isPressed = false
		}

		if event.Position.In(t.labelLoaded.systemSetBounds()) {
			t.enabled.Set(!t.enabled.Get())
		}
	}
}

func (t *toggleImpl) HandleKeyEvent(event keyEvent) {}

func (t *toggleImpl) HandleInputEvent(event inputEvent) {}
