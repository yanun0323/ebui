package ebui

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implementation */
var _ ebiten.Game = (*app)(nil)

type app struct {
	view  SomeView
	debug bool
}

func Run(title string, contentView View, debug ...bool) {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	app := &app{
		view:  contentView.Body(),
		debug: len(debug) != 0 && debug[0],
	}

	if app.debug {
		logs.SetDefaultLevel(logs.LevelDebug)
	}

	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return
		}

		logs.Fatal(err)
	}
}

func (a *app) SetWindowSize(w, h int) {
	ebiten.SetWindowSize(w, h)
}

func (a *app) Update() error {
	EbitenUpdate(a.view)
	runtime.GC()
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	EbitenDraw(screen, a.view)

	if a.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f, FPS: %.1f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	}
}

func (a *app) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
