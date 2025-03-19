package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/input"
	layout "github.com/yanun0323/ebui/layout"
)

const (
	_defaultToggleSize    = 30
	_defaultTogglePadding = 2
)

var (
	defaultToggleOnColor            = ivory
	defaultToggleOffColor           = ivory
	defaultToggleOnBackgroundColor  = NewColor(64, 191, 64)
	defaultToggleOffBackgroundColor = NewColor(64, 64, 64)
	defaultToggleOffset             = NewPoint(_defaultToggleSize-(2*_defaultTogglePadding), 0)
)

type toggleImpl struct {
	*viewCtx

	label       func() SomeView
	labelLoaded SomeView
	enabled     *Binding[bool]
	isPressed   bool

	defaultToggleOffset          *Binding[CGPoint]
	defaultToggleColor           *Binding[CGColor]
	defaultToggleBackgroundColor *Binding[CGColor]
}

func Toggle(enabled *Binding[bool], label ...func() SomeView) SomeView {
	t := &toggleImpl{
		enabled: enabled,
	}

	if len(label) != 0 && label[0] != nil {
		t.label = label[0]
	} else {
		t.defaultToggleOffset = BindOneWay(t.enabled, func(enabled bool) CGPoint {
			if enabled {
				return defaultToggleOffset
			}
			return CGPoint{}
		})
		t.defaultToggleColor = BindOneWay(t.enabled, func(enabled bool) CGColor {
			if enabled {
				return defaultToggleOnColor
			}
			return defaultToggleOffColor
		})
		t.defaultToggleBackgroundColor = BindOneWay(t.enabled, func(enabled bool) CGColor {
			if enabled {
				return defaultToggleOnBackgroundColor
			}
			return defaultToggleOffBackgroundColor
		})
		// if t.enabled.Get() {
		// 	t.defaultToggleOffset.Set(defaultToggleOffset)
		// 	t.defaultToggleColor.Set(defaultToggleOnColor)
		// 	t.defaultToggleBackgroundColor.Set(defaultToggleOnBackgroundColor)
		// }
		// enabled.AddListener(func(_, newVal bool, animStyle ...animation.Style) {
		// 	if newVal {
		// 		t.defaultToggleOffset.Set(defaultToggleOffset, animStyle...)
		// 		t.defaultToggleColor.Set(defaultToggleOnColor, animStyle...)
		// 		t.defaultToggleBackgroundColor.Set(defaultToggleOnBackgroundColor, animStyle...)
		// 	} else {
		// 		t.defaultToggleOffset.Set(CGPoint{}, animStyle...)
		// 		t.defaultToggleColor.Set(defaultToggleOffColor, animStyle...)
		// 		t.defaultToggleBackgroundColor.Set(defaultToggleOffBackgroundColor, animStyle...)
		// 	}
		// })

		t.labelLoaded = t.defaultLabel()
	}

	t.viewCtx = newViewContext(t)

	return t
}

func (b *toggleImpl) preload(parent *viewCtxEnv, _ ...stackType) (preloadData, layoutFunc) {
	if b.label != nil {
		b.labelLoaded = b.label()
	}
	if b.labelLoaded == nil {
		panic("empty view from button label")
	}

	formulaStack := &stackPreloader{
		types:    stackTypeZStack,
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

func (b *toggleImpl) defaultLabel() SomeView {
	return HStack(
		Circle().
			Frame(Const(NewSize(_defaultToggleSize, _defaultToggleSize))).
			Offset(b.defaultToggleOffset).
			BackgroundColor(b.defaultToggleColor).
			Padding(Const(NewInset(_defaultTogglePadding, _defaultTogglePadding, _defaultTogglePadding, _defaultTogglePadding))),
	).
		Frame(Const(NewSize(60, _defaultToggleSize+_defaultTogglePadding*2))).
		BackgroundColor(b.defaultToggleBackgroundColor).
		RoundCorner(Const(float64(_defaultToggleSize / 2))).
		Padding(Const(NewInset(5, 5, 5, 5))).
		Align(Bind(layout.AlignLeading | layout.AlignTop))
}

func (t *toggleImpl) onMouseEvent(event input.MouseEvent) {
	defer t.viewCtx.onMouseEvent(event)

	if t.viewCtxEnv.disabled.Get() {
		return
	}

	switch event.Phase {
	case input.MousePhaseBegan:
		if t.labelLoaded.systemSetBounds().Contains(event.Position) {
			t.isPressed = true
		}
	case input.MousePhaseEnded, input.MousePhaseCancelled:
		if t.isPressed {
			t.isPressed = false
		}

		if t.labelLoaded.systemSetBounds().Contains(event.Position) {
			t.enabled.Set(!t.enabled.Get())
		}
	}
}

func (t *toggleImpl) processable() bool {
	return true
}
