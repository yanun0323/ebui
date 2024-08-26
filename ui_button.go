package ebui

import (
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/* Check Interface Implementation */
var _ SomeView = (*buttonView)(nil)

func Button(action func(), label View) *buttonView {
	v := &buttonView{
		action: action,
	}

	v.uiView = newUIView(typeButton, v, label)
	return v
}

type buttonView struct {
	*uiView

	action     func()
	invokeTick atomic.Int64
}

func (v *buttonView) Body() SomeView {
	return v
}

func (v *buttonView) draw(screen *ebiten.Image) {
	cache := v.Copy()
	cache.Draw(screen, func(screen *ebiten.Image) {
		// FIXME: Fix me
		cX, cY := ebiten.CursorPosition()
		cache.isPressing = cache.Contain(cX, cY) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		if cache.Contain(cX, cY) {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				currentTick := currentTicker()
				previousTick := v.invokeTick.Load()
				if previousTick != currentTicker() {
					v.action()
					v.invokeTick.Store(currentTick)
				}
			}
		}

		cache.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, cache)
			if v != nil {
				v.draw(screen)
			}
		})
	})
}
