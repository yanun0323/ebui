package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
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
}

func resetLayout(contentView SomeView) {
	bounds := globalStateManager.GetBounds()
	_, _, layoutFn := contentView.preload(nil)
	_ = layoutFn(bounds.Start, bounds.Size())
	globalStateManager.clearDirty()
}

func EbitenDraw(screen *ebiten.Image, contentView SomeView) {
	contentView.draw(screen)
}

func EbitenLayout(outsideWidth, outsideHeight int) {
	globalStateManager.SetBounds(NewRect(0, 0, outsideWidth, outsideHeight))
}
