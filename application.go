package ebui

import (
	"errors"
	"image/color"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	_defaultWidth  = 500
	_defaultHeight = 500
)

type windowResizingMode int

var (
	Terminate                                                  = ebiten.Termination
	WindowResizingModeDisabled              windowResizingMode = windowResizingMode(ebiten.WindowResizingModeDisabled)
	WindowResizingModeOnlyFullscreenEnabled windowResizingMode = windowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)
	WindowResizingModeEnabled               windowResizingMode = windowResizingMode(ebiten.WindowResizingModeEnabled)

	resourceDir atomic.Value
)

type application struct {
	stateManager    *stateManager
	eventManager    *eventManager
	animManager     *animationManager
	rootView        SomeView
	backgroundColor color.Color
	bounds          CGRect
}

func NewApplication(root View) *application {
	app := &application{
		stateManager: globalStateManager,
		eventManager: globalEventManager,
		animManager:  globalAnimationManager,
		rootView:     ZStack(root.Body()),
		bounds:       rect(0, 0, _defaultWidth, _defaultHeight),
	}

	app.SetWindowSize(_defaultWidth, _defaultHeight)
	return app
}

func (app *application) Run() error {
	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return Terminate
		}

		return err
	}

	return nil
}

func (app *application) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	app.SetWindowSize(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (app *application) Update() error {
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
			app.eventManager.DispatchTouchEvent(touchEvent{
				Phase:    touchPhaseBegan,
				Position: pos,
			})
		} else {
			app.eventManager.DispatchTouchEvent(touchEvent{
				Phase:    touchPhaseMoved,
				Position: pos,
			})
		}
	} else if app.eventManager.isTracking {
		x, y := ebiten.CursorPosition()
		app.eventManager.DispatchTouchEvent(touchEvent{
			Phase:    touchPhaseEnded,
			Position: pt(float64(x), float64(y)),
		})
	}

	return nil
}

func (app *application) Draw(screen *ebiten.Image) {
	baseHiDPI := ebiten.NewImage(int(app.bounds.Dx()), int(app.bounds.Dy()))
	if app.backgroundColor != nil {
		baseHiDPI.Fill(app.backgroundColor)
	}
	app.rootView.draw(baseHiDPI)
	screen.DrawImage(baseHiDPI, nil)
}

// SetWindowBackgroundColor sets the background color of the application.
func (app *application) SetWindowBackgroundColor(color color.Color) {
	app.backgroundColor = color
}

// SetWindowSize sets the size of the window.
func (app *application) SetWindowSize(width, height int) {
	bounds := rect(0, 0, float64(width), float64(height))
	app.bounds = bounds
	app.stateManager.SetBounds(bounds)
	app.reLayout()
}

func (app *application) SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func (app *application) SetWindowResizingMode(mode windowResizingMode) {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeType(mode))
}

func (app *application) SetResourceFolder(folder string) {
	resourceDir.Store(folder)
}

func (app *application) reLayout() {
	_, _, layoutFn := app.rootView.preload()
	_ = layoutFn(ptZero, app.bounds.Size())
}
