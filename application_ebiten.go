package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	popupViews []SomeView
)

func EbitenUpdate(contentView SomeView) {
	// 1. 更新動畫
	globalAnimationManager.Update()

	// 2. 處理狀態更新
	if globalStateManager.isDirty() {
		resetLayout(contentView)
	}

	// 3. 處理輸入事件
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

	altPressing := ebiten.IsKeyPressed(ebiten.KeyAltLeft) || ebiten.IsKeyPressed(ebiten.KeyAltRight) || ebiten.IsKeyPressed(ebiten.KeyAlt)
	shiftPressing := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) || ebiten.IsKeyPressed(ebiten.KeyShift)
	controlPressing := ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight) || ebiten.IsKeyPressed(ebiten.KeyControl)
	metaPressing := ebiten.IsKeyPressed(ebiten.KeyMeta) || ebiten.IsKeyPressed(ebiten.KeyMetaLeft) || ebiten.IsKeyPressed(ebiten.KeyMetaRight)

	// 4. 處理鍵盤事件
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

	// 5. 處理輸入事件
	input := ebiten.AppendInputChars(nil)
	for _, char := range input {
		globalEventManager.DispatchInputEvent(inputEvent{
			Char: char,
		})
	}
}

func resetLayout(contentView SomeView) {
	bounds := globalStateManager.GetBounds()
	_, layoutFn := contentView.preload(nil)
	_, _ = layoutFn(bounds.Start, bounds.Size())
	globalStateManager.clearDirty()
}

func EbitenDraw(screen *ebiten.Image, contentView SomeView) {
	contentView.draw(screen)

	for i := range popupViews {
		popupViews[i].draw(screen)
	}
}

func EbitenLayout(outsideWidth, outsideHeight int) {
	globalStateManager.SetBounds(NewRect(0, 0, outsideWidth, outsideHeight))
}
