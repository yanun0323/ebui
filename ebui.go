package ebui

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yanun0323/pkg/logs"
)

var (
	_rootViewCache = atomic.Value{}
)

/* Check Interface Implementation */
var _ ebiten.Game = (*app)(nil)

type app struct {
	contentView SomeView
	debug       bool
}

func Run(title string, contentView View, debug ...bool) error {
	logs.Info("initializing app...")
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	app := &app{
		contentView: contentView.Body(),
		debug:       len(debug) != 0 && debug[0],
	}

	logs.Info("starting app...")
	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return nil
		}

		return fmt.Errorf("run app, err: %w", err)
	}

	return nil
}

func (a *app) SetWindowSize(w, h int) {
	ebiten.SetWindowSize(w, h)
}

func (a *app) Update() error {
	sync.OnceFunc(func() {
		logs.Info("starting app...")
	})

	EbitenUpdate(a.contentView)
	runtime.GC()

	sync.OnceFunc(func() {
		logs.Info("app started successfully!")
	})
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	EbitenDraw(screen)

	if a.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f, FPS: %.1f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	}
}

func (a *app) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func EbitenUpdate(sv SomeView) {
	if sv != nil {
		w, h := ebiten.WindowSize()
		flexSize := size{w, h}
		pos := point{}
		v := newView(typesNone, nil, sv.Body())

		logs.Info("=== preloading...")
		v.preloadSize()
		v.size = flexSize

		logs.Info("=== setting size position...")
		v.setSizePosition(flexSize, &pos)
		_rootViewCache.Store(sv.Body())
	}

	tickTock()
}

func EbitenDraw(screen *ebiten.Image) {
	if r, ok := _rootViewCache.Load().(uiViewDelegator); ok {
		r.UIView().Draw(screen)
	}
}

func invokeSomeView(sv SomeView, fn func(*uiView)) {
	if sv, ok := sv.(uiViewDelegator); ok {
		fn(sv.UIView())
	}
}
