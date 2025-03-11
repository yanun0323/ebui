package ebui

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type windowResizingMode int

var (
	Terminate                                                  = ebiten.Termination
	WindowResizingModeDisabled              windowResizingMode = windowResizingMode(ebiten.WindowResizingModeDisabled)
	WindowResizingModeOnlyFullscreenEnabled windowResizingMode = windowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)
	WindowResizingModeEnabled               windowResizingMode = windowResizingMode(ebiten.WindowResizingModeEnabled)

	resourceDir *Binding[string] = Bind("")
)

type application struct {
	isDebug         bool
	debugInfo       string
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
func (app *application) SetWindowBackgroundColor(color CGColor) {
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
	resourceDir.Set(folder)
	println("set resource folder:", folder)
}

func (app *application) Debug() {
	app.isDebug = true
}

func (app *application) VSyncEnabled(enabled bool) {
	ebiten.SetVsyncEnabled(enabled)
}

func (app *application) Run(title string) error {
	ebiten.SetWindowTitle(title)

	if err := ebiten.RunGame(&game{application: app}); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return Terminate
		}

		return err
	}

	return nil
}

type game struct {
	*application
}

func (app *game) Update() error {
	EbitenUpdate(app.rootView)
	if app.isDebug {
		count := app.rootView.count()
		app.debugInfo = fmt.Sprintf("TPS: %.2f, FPS: %.2f, ViewCount: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), count)
	}

	return nil
}

func (app *game) Draw(screen *ebiten.Image) {
	if app.backgroundColor != nil {
		screen.Fill(app.backgroundColor)
	}

	EbitenDraw(screen, app.rootView)

	if app.isDebug {
		ebitenutil.DebugPrint(screen, app.debugInfo)
	}
}

func (app *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	EbitenLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
