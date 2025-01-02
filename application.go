package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Application struct {
	stateManager    *StateManager
	eventManager    *EventManager
	animManager     *AnimationManager
	rootView        View
	backgroundColor color.Color
	bounds          image.Rectangle
}

func NewApplication(root SomeView) *Application {
	bounds := image.Rect(0, 0, 400, 300)
	defaultStateManager.SetBounds(bounds)

	app := &Application{
		stateManager:    defaultStateManager,
		eventManager:    defaultEventManager,
		animManager:     NewAnimationManager(),
		backgroundColor: color.White,
		bounds:          bounds,
	}

	app.rootView = root.Build()

	return app
}

func (app *Application) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	bounds := image.Rect(0, 0, outsideWidth, outsideHeight)
	app.rootView.Layout(bounds)
	return outsideWidth, outsideHeight
}

func (app *Application) Update() error {
	// 1. 更新動畫
	app.animManager.Update()

	// 2. 處理狀態更新
	if app.stateManager.isDirty() {
		app.rootView = app.rootView.Build()
		app.rootView.Layout(app.bounds)
		app.stateManager.clearDirty()
	}

	// 3. 處理輸入事件
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		pos := image.Pt(x, y)

		if !app.eventManager.isTracking {
			app.eventManager.DispatchTouchEvent(TouchEvent{
				Phase:    TouchPhaseBegan,
				Position: pos,
			})
		} else {
			app.eventManager.DispatchTouchEvent(TouchEvent{
				Phase:    TouchPhaseMoved,
				Position: pos,
			})
		}
	} else if app.eventManager.isTracking {
		x, y := ebiten.CursorPosition()
		app.eventManager.DispatchTouchEvent(TouchEvent{
			Phase:    TouchPhaseEnded,
			Position: image.Pt(x, y),
		})
	}

	return nil
}

func (app *Application) Draw(screen *ebiten.Image) {
	app.rootView.Draw(screen)
}
