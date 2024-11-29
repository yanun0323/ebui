package ebui

import (
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/* Check Interface Implementation */
var _ SomeView = (*buttonView)(nil)

func Button(action func(), label View) SomeView {
	v := &buttonView{
		action: action,
		label:  label.Body(),
	}

	v.uiView = newView(typesButton, v, label)
	return v
}

type buttonView struct {
	*uiView
	label SomeView

	action     func()
	invokeTick atomic.Int64
}

func (v *buttonView) Body() SomeView {
	return v
}

func (v *buttonView) deepUpdateAction() {
	cX, cY := ebiten.CursorPosition()
	isPressing := v.label.isPress(cX, cY)
	v.isPressing = isPressing && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isPressing {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			currentTick := currentTicker()
			if v.invokeTick.Load() != currentTicker() {
				v.action()
				v.invokeTick.Store(currentTick)
			}
		}
	}

	v.uiView.deepUpdateAction()
	v.uiView.deepUpdateEnvironment()
}
