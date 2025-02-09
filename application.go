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
	rootView        SomeView
	backgroundColor color.Color
	bounds          image.Rectangle
}

func NewApplication(root View) *Application {
	bounds := image.Rect(0, 0, 400, 300)
	defaultStateManager.SetBounds(bounds)

	app := &Application{
		stateManager: defaultStateManager,
		eventManager: defaultEventManager,
		animManager:  defaultAnimationManager,
		rootView:     root.Body(),
		bounds:       bounds,
	}

	return app
}

func (app *Application) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	bounds := image.Rect(0, 0, outsideWidth, outsideHeight)
	app.bounds = bounds
	app.rootView.layout(bounds)
	return outsideWidth, outsideHeight
}

func (app *Application) Update() error {
	// 1. 更新動畫
	app.animManager.Update()

	// 2. 處理狀態更新
	if app.stateManager.isDirty() {
		app.rootView.layout(app.bounds)
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
	if app.backgroundColor != nil {
		screen.Fill(app.backgroundColor)
	}
	app.rootView.Body().draw(screen)
}

// 設置背景顏色
func (app *Application) SetBackgroundColor(color color.Color) {
	app.backgroundColor = color
}

// 設置視窗大小
func (app *Application) SetBounds(bounds image.Rectangle) {
	app.bounds = bounds
	app.stateManager.SetBounds(bounds)
	app.rootView.Body().layout(bounds)
}
