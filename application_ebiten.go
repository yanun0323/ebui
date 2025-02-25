package ebui

import "github.com/hajimehoshi/ebiten/v2"

func EbitenUpdate(contentView SomeView) {
	w, h := ebiten.WindowSize()
	size := CGSize{float64(w), float64(h)}

	// 1. 更新動畫
	globalAnimationManager.Update()

	// 2. 處理狀態更新
	if globalStateManager.isDirty() {
		resetLayout(contentView, size)
		globalStateManager.clearDirty()
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
}

func resetLayout(contentView SomeView, size CGSize) {
	_, _, layoutFn := contentView.preload()
	_ = layoutFn(CGPoint{}, size)
}

func EbitenDraw(screen *ebiten.Image, contentView SomeView) {
	contentView.draw(screen)
}

func EbitenLayout(outsideWidth, outsideHeight int) {
	globalStateManager.SetBounds(NewRect(0, 0, outsideWidth, outsideHeight))
}
