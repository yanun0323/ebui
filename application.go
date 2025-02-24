package ebui

import (
	"image/color"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	_defaultWidth  = 500
	_defaultHeight = 500
)

var (
	resourceDir atomic.Value
)

type Application struct {
	stateManager    *StateManager
	eventManager    *EventManager
	animManager     *AnimationManager
	rootView        SomeView
	backgroundColor color.Color
	bounds          CGRect
}

func NewApplication(root View) *Application {
	app := &Application{
		stateManager: globalStateManager,
		eventManager: globalEventManager,
		animManager:  globalAnimationManager,
		rootView:     ZStack(root.Body()),
		bounds:       rect(0, 0, _defaultWidth, _defaultHeight),
	}

	app.SetBounds(_defaultWidth, _defaultHeight)
	return app
}

func (app *Application) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	app.SetBounds(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (app *Application) Update() error {
	// 1. 更新動畫
	app.animManager.Update()

	// 2. 處理狀態更新
	if app.stateManager.isDirty() {
		app.reLayout()
		app.stateManager.clearDirty()
	}

	// 3. 處理輸入事件
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		pos := pt(float64(x), float64(y))

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
			Position: pt(float64(x), float64(y)),
		})
	}

	return nil
}

func (app *Application) Draw(screen *ebiten.Image) {
	baseHiDPI := ebiten.NewImage(int(app.bounds.Dx()), int(app.bounds.Dy()))
	if app.backgroundColor != nil {
		baseHiDPI.Fill(app.backgroundColor)
	}
	app.rootView.draw(baseHiDPI)
	screen.DrawImage(baseHiDPI, nil)
}

// 設置背景顏色
func (app *Application) SetBackgroundColor(color color.Color) {
	app.backgroundColor = color
}

// 設置視窗大小
func (app *Application) SetBounds(width, height int) {
	bounds := rect(0, 0, float64(width), float64(height))
	app.bounds = bounds
	app.stateManager.SetBounds(bounds)
	app.reLayout()
}

func (app *Application) SetResourceFolder(folder string) {
	resourceDir.Store(folder)
}

func (app *Application) reLayout() {
	_, _, layoutFn := app.rootView.preload()
	_ = layoutFn(ptZero, app.bounds.Size())
}
