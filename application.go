package ebui

import (
	"errors"
	"image/color"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
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
	rootView        SomeView
	backgroundColor color.Color
}

func NewApplication(root View) *application {
	app := &application{
		rootView: ZStack(root.Body()),
	}
	return app
}

// SetWindowBackgroundColor sets the background color of the application.
func (app *application) SetWindowBackgroundColor(color AnyColor) {
	app.backgroundColor = color
}

// SetWindowSize sets the size of the window.
func (app *application) SetWindowSize(width, height int) {
	EbitenLayout(width, height)
}

func (app *application) SetWindowResizingMode(mode windowResizingMode) {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeType(mode))
}

func (app *application) SetResourceFolder(folder string) {
	resourceDir.Store(folder)
}

func (app *application) Run(title string) error {
	ebiten.SetWindowTitle(title)

	if err := ebiten.RunGame(&game{rootView: app.rootView}); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return Terminate
		}

		return err
	}

	return nil
}

type game struct {
	rootView SomeView
}

func (app *game) Update() error {
	EbitenUpdate(app.rootView)
	return nil
}

func (app *game) Draw(screen *ebiten.Image) {
	EbitenDraw(screen, app.rootView)
}

func (app *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	EbitenLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
