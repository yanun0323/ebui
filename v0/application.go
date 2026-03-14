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
	isDebug           bool
	debugInfo         string
	rootView          SomeView
	backgroundColor   color.Color
	initUnfocused     bool
	screenTransparent bool
	skipTaskbar       bool
	singleThread      bool
	layoutHook        func()
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
	ebiten.SetWindowSize(width, height)
	EbitenLayout(width, height)
}

func (app *application) SetWindowPosition(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (app *application) SetWindowResizingMode(mode windowResizingMode) {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeType(mode))
}

func (app *application) SetLayoutHook(hook func()) {
	app.layoutHook = hook
}

// SetRunWithoutFocus indicates whether the window is unfocused or not on launching.
// SetRunWithoutFocus is valid on desktops and browsers.
//
// The default (zero) value is false, which means that the window is focused.
func (app *application) SetRunWithoutFocus(initUnfocused bool) {
	ebiten.SetRunnableOnUnfocused(initUnfocused)
	app.initUnfocused = initUnfocused
}

// SetScreenTransparent indicates whether the window is transparent or not.
// SetScreenTransparent is valid on desktops and browsers.
//
// The default (zero) value is false, which means that the window is not transparent.
func (app *application) SetScreenTransparent(screenTransparent bool) {
	app.screenTransparent = screenTransparent
}

// SetSkipTaskbar indicates whether an application icon is shown on a taskbar or not.
// SetSkipTaskbar is valid only on Windows.
//
// The default (zero) value is false, which means that an icon is shown on a taskbar.
func (app *application) SetSkipTaskbar(skipTaskbar bool) {
	app.skipTaskbar = skipTaskbar
}

// SetSingleThread indicates whether the single thread mode is used explicitly or not.
//
// The single thread mode disables Ebitengine's thread safety to unlock maximum performance.
// If you use this you will have to manage threads yourself.
// Functions like `SetWindowSize` will no longer be concurrent-safe with this build tag.
// They must be called from the main thread or the same goroutine as the given game's callback functions like Update.
//
// SetSingleThread works only with desktops and consoles.
//
// If SetSingleThread is false, and if the build tag `ebitenginesinglethread` is specified,
// the single thread mode is used.
//
// The default (zero) value is false, which means that the single thread mode is disabled.
func (app *application) SetSingleThread(singleThread bool) {
	app.singleThread = singleThread
}

func (app *application) SetWindowFloating(floating bool) {
	ebiten.SetWindowFloating(floating)
}

func (app *application) SetResourceFolder(folder string) {
	resourceDir.Set(folder)
	logf("set resource folder: %s", folder)
}

func (app *application) Debug() {
	app.isDebug = true
}

func (app *application) VSyncEnabled(enabled bool) {
	ebiten.SetVsyncEnabled(enabled)
}

func (app *application) Run(title string) error {
	ebiten.SetWindowTitle(title)

	if err := ebiten.RunGameWithOptions(&game{application: app}, &ebiten.RunGameOptions{
		InitUnfocused:     app.initUnfocused,
		ScreenTransparent: app.screenTransparent,
		SkipTaskbar:       app.skipTaskbar,
		SingleThread:      app.singleThread,
	}); err != nil {
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
	if app.layoutHook != nil {
		app.layoutHook()
	}

	EbitenLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
