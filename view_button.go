package ebui

import (
	"image"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implement */
var _ SomeView = (*buttonView)(nil)

func Button(action func(), label View) *buttonView {
	v := &buttonView{
		action: action,
	}

	v.viewOption = newViewOption(v, label)
	return v
}

type buttonView struct {
	viewOption

	action     func()
	invokeTick atomic.Int64
}

func (v *buttonView) Body() SomeView {
	return v
}

func (v *buttonView) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	current := parent.calculateViewOption(v.viewOption)
	logs.Debugf("ButtonView: %+v", current)
	current.Draw(screen, func(screen *ebiten.Image) {
		cX, cY := ebiten.CursorPosition()
		current.isPressing = current.Contain(cX, cY) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		if current.Contain(cX, cY) {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				currentTick := currentTicker()
				previousTick := v.invokeTick.Load()
				if previousTick != currentTicker() {
					v.action()
					v.invokeTick.Store(currentTick)
				}
			}
		}

		current.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, &current)
			if v != nil {
				_ = v.draw(screen, current.CreateChild(v, current.XX(), current.YY(), current.XX(), current.YY(), current.Width(), current.Height()))
			}
		})
	})

	return current.DrawnArea()
}

func (v *buttonView) bounds() (int, int) {
	return v.w, v.h
}
