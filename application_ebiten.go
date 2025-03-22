package ebui

import (
	"fmt"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yanun0323/ebui/input"
	"github.com/yanun0323/ebui/internal/helper"
)

var (
	popupViews            []SomeView
	isLeftButtonTracking atomic.Bool
)

// CursorPosition returns a position of a mouse cursor relative to the game screen (window). The cursor position is
// 'logical' position and this considers the scale of the screen.
//
// CursorPosition returns (0, 0) before the main loop on desktops and browsers.
//
// CursorPosition always returns (0, 0) on mobile native applications.
//
// CursorPosition is concurrent-safe.
func CursorPosition() (x, y int) {
	return ebiten.CursorPosition()
}

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

	// 4. handle wheel events

	x, y := ebiten.CursorPosition()
	cursor := newVector(x, y)

	contentView.onHoverEvent(cursor)

	dx, dy := ebiten.Wheel()
	speed := DefaultScrollSpeed.Value()

	contentView.onScrollEvent(cursor, input.ScrollEvent{
		Delta: newVector(-dx*speed, -dy*speed),
	})

	// 5. handle touch events

	leftButtonTracking := isLeftButtonTracking.Load()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if leftButtonTracking {
			contentView.onMouseEvent(input.MouseEvent{
				Phase:    input.MousePhaseMoved,
				Position: cursor,
			})
		} else {
			contentView.onMouseEvent(input.MouseEvent{
				Phase:    input.MousePhaseBegan,
				Position: cursor,
			})
			isLeftButtonTracking.Store(true)
		}
	} else if leftButtonTracking {
		contentView.onMouseEvent(input.MouseEvent{
			Phase:    input.MousePhaseEnded,
			Position: cursor,
		})
		isLeftButtonTracking.Store(false)
	}

	mMouse := m.ElapsedAndReset()

	// 6. handle keyboard events
	altPressing := ebiten.IsKeyPressed(ebiten.KeyAltLeft) || ebiten.IsKeyPressed(ebiten.KeyAltRight) || ebiten.IsKeyPressed(ebiten.KeyAlt)
	shiftPressing := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) || ebiten.IsKeyPressed(ebiten.KeyShift)
	controlPressing := ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight) || ebiten.IsKeyPressed(ebiten.KeyControl)
	metaPressing := ebiten.IsKeyPressed(ebiten.KeyMeta) || ebiten.IsKeyPressed(ebiten.KeyMetaLeft) || ebiten.IsKeyPressed(ebiten.KeyMetaRight)

	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		contentView.onKeyEvent(input.KeyEvent{
			Key:     input.Key(key),
			Phase:   input.KeyPhaseJustPressed,
			Shift:   shiftPressing,
			Control: controlPressing,
			Alt:     altPressing,
			Meta:    metaPressing,
		})
	}

	keys = inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		contentView.onKeyEvent(input.KeyEvent{
			Key:     input.Key(key),
			Phase:   input.KeyPhasePressing,
			Shift:   shiftPressing,
			Control: controlPressing,
			Alt:     altPressing,
			Meta:    metaPressing,
		})
	}

	keys = inpututil.AppendJustReleasedKeys(nil)
	for _, key := range keys {
		contentView.onKeyEvent(input.KeyEvent{
			Key:     input.Key(key),
			Phase:   input.KeyPhaseJustReleased,
			Shift:   shiftPressing,
			Control: controlPressing,
			Alt:     altPressing,
			Meta:    metaPressing,
		})
	}

	mKeyboard := m.ElapsedAndReset()

	// 7. handle input events
	inputs := ebiten.AppendInputChars(nil)
	for _, char := range inputs {
		contentView.onTypeEvent(input.TypeEvent{
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
