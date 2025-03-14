package ebui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yanun0323/ebui/internal/helper"
)

var (
	popupViews []SomeView
)

// EbitenUpdate updates the application state
//
// should be called in ebiten.Update
func EbitenUpdate(contentView SomeView) {
	m := helper.NewMetric()

	// 1. update animations
	globalAnimationManager.Update()

	mAnim := m.ElapsedAndReset()

	// 2. handle state updates
	if globalStateManager.isDirty() {
		resetLayout(contentView)
	}

	mLayout := m.ElapsedAndReset()

	// 3. handle wheel events
	x, y := ebiten.Wheel()
	if x != 0 || y != 0 {
		logf("scrolling, x: %.2f, y: %.2f", x, y)
		speed := DefaultScrollSpeed.Value()
		globalEventManager.DispatchWheelEvent(wheelEvent{
			Delta: NewPoint(x*speed, y*speed),
		})
	}

	// 4. handle touch events
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		pos := NewPoint(float64(x), float64(y))

		if !globalEventManager.isTracking {
			globalEventManager.DispatchTouchEvent(touchEvent{
				Phase:    touchPhaseBegan,
				Position: pos,
			})
		} else {
			globalEventManager.DispatchTouchEvent(touchEvent{
				Phase:    touchPhaseMoved,
				Position: pos,
			})
		}
	} else if globalEventManager.isTracking {
		x, y := ebiten.CursorPosition()
		globalEventManager.DispatchTouchEvent(touchEvent{
			Phase:    touchPhaseEnded,
			Position: NewPoint(float64(x), float64(y)),
		})
	}

	mMouse := m.ElapsedAndReset()

	// 5. handle keyboard events
	altPressing := ebiten.IsKeyPressed(ebiten.KeyAltLeft) || ebiten.IsKeyPressed(ebiten.KeyAltRight) || ebiten.IsKeyPressed(ebiten.KeyAlt)
	shiftPressing := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) || ebiten.IsKeyPressed(ebiten.KeyShift)
	controlPressing := ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight) || ebiten.IsKeyPressed(ebiten.KeyControl)
	metaPressing := ebiten.IsKeyPressed(ebiten.KeyMeta) || ebiten.IsKeyPressed(ebiten.KeyMetaLeft) || ebiten.IsKeyPressed(ebiten.KeyMetaRight)

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		globalEventManager.DispatchKeyEvent(keyEvent{
			Key:     key,
			Phase:   keyPhaseJustPressed,
			Shift:   shiftPressing,
			Control: controlPressing,
			Alt:     altPressing,
			Meta:    metaPressing,
		})
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		globalEventManager.DispatchKeyEvent(keyEvent{
			Key:     key,
			Phase:   keyPhasePressing,
			Shift:   shiftPressing,
			Control: controlPressing,
			Alt:     altPressing,
			Meta:    metaPressing,
		})
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		globalEventManager.DispatchKeyEvent(keyEvent{
			Key:     key,
			Phase:   keyPhaseJustReleased,
			Shift:   shiftPressing,
			Control: controlPressing,
			Alt:     altPressing,
			Meta:    metaPressing,
		})
	}

	mKeyboard := m.ElapsedAndReset()

	// 6. handle input events
	input := ebiten.AppendInputChars(nil)
	for _, char := range input {
		globalEventManager.DispatchInputEvent(inputEvent{
			Char: char,
		})
	}

	mInput := m.ElapsedAndReset()

	m.Reset()

	if false {
		fmt.Printf("\x1b[34m[UPD]\x1b[0m\t anim: %.4f, layout: %.4f, mouse: %.4f, keyboard: %.4f, input: %.4f\n",
			mAnim.Seconds(), mLayout.Seconds(), mMouse.Seconds(), mKeyboard.Seconds(), mInput.Seconds())
	}
}

func resetLayout(contentView SomeView) {
	bounds := globalStateManager.GetBounds()
	_, layoutFn := contentView.preload(nil)
	_, _ = layoutFn(bounds.Start, bounds.Size())
	globalStateManager.clearDirty()
}

// EbitenDraw draws the application state
//
// should be called in ebiten.Draw
func EbitenDraw(screen *ebiten.Image, contentView SomeView) {
	m := helper.NewMetric()

	contentView.draw(screen)

	mDraw := m.ElapsedAndReset()

	for i := range popupViews {
		popupViews[i].draw(screen)
	}

	mDrawPopup := m.ElapsedAndReset()

	if false {
		fmt.Printf("\x1b[36m[DRAW]\x1b[0m\t draw: %.4f, draw popup: %.4f\n",
			mDraw.Seconds(), mDrawPopup.Seconds())
	}
}

// EbitenLayout sets the application bounds
//
// should be called in ebiten.Layout
func EbitenLayout(outsideWidth, outsideHeight int) {
	globalStateManager.SetBounds(NewRect(0, 0, outsideWidth, outsideHeight))
}
