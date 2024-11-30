package ebui

import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yanun0323/ebui/sync"
	"github.com/yanun0323/pkg/logs"
)

var (
	rootViewCache = atomic.Value{}
	once          sync.Once
)

/* Check Interface Implementation */
var _ ebiten.Game = (*app)(nil)

type app struct {
	contentView SomeView
	debug       bool
	memMonitor  bool
	memStatus   sync.Value[memStats]
}

func Run(title string, contentView SomeView, options ...RunOption) error {
	logs.Info("initializing app...")
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	app := &app{
		contentView: contentView,
	}

	for _, option := range options {
		option(app)
	}

	logs.Info("starting app...")

	if app.memMonitor {
		go app.monitorMem()
	}

	if err := ebiten.RunGame(app); err != nil {
		if errors.Is(err, ebiten.Termination) {
			return nil
		}

		return fmt.Errorf("run app, err: %w", err)
	}

	return nil
}

func SetWindowSize(w, h int) {
	ebiten.SetWindowSize(w, h)
}

func (a *app) Update() error {

	EbitenUpdate(a.contentView)
	runtime.GC()

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	EbitenDraw(a.contentView, screen)

	if a.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("[Frame Rate] TPS: %.1f, FPS: %.1f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), a.memStatus.Load().string()))
	}
}

func (a *app) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func EbitenUpdate(v SomeView) {
	if v != nil {
		// 		w, h := ebiten.WindowSize()
		// 		println(w, h)
		// 		v := newView(typesNone, nil, sv.Body())
		// 		v.deepReset()
		// 		v.setSize(size{w, h})
		// 		layout(v, point{}, size{w, h})

		// 		v.deepUpdateAction()
		// 		v.deepUpdateEnvironment()
	}

	// tickTock()
}

func EbitenDraw(v SomeView, screen *ebiten.Image) {
	if v != nil {
		v.draw(screen)
	}
}

func (a *app) monitorMem() {
	time.Sleep(time.Second)
	for {
		mStats := runtime.MemStats{}
		runtime.ReadMemStats(&mStats)
		stats := a.memStatus.Load()
		stats.update(mStats)
		a.memStatus.Store(stats)

		time.Sleep(time.Second)
	}
}

var (
	unit     float64 = 1024
	unitStep         = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
)

type memStats struct {
	isInit      bool
	HeapSys     uint64
	HeapAlloc   uint64
	HeapInuse   uint64
	HeapObjects uint64
	StackSys    uint64
	StackInuse  uint64
}

func (m *memStats) update(s runtime.MemStats) {
	if !m.isInit {
		m.isInit = true

		m.HeapSys = s.HeapSys
		m.HeapAlloc = s.HeapAlloc
		m.HeapInuse = s.HeapInuse
		m.HeapObjects = s.HeapObjects
		m.StackSys = s.StackSys
		m.StackInuse = s.StackInuse

		return
	}

	m.HeapSys = (m.HeapSys + s.HeapSys) / 2
	m.HeapAlloc = (m.HeapAlloc + s.HeapAlloc) / 2
	m.HeapInuse = (m.HeapInuse + s.HeapInuse) / 2
	m.HeapObjects = (m.HeapObjects + s.HeapObjects) / 2
	m.StackSys = (m.StackSys + s.StackSys) / 2
	m.StackInuse = (m.StackInuse + s.StackInuse) / 2
}

func (m memStats) string() string {
	if !m.isInit {
		return ""
	}

	return fmt.Sprintf("[Heap] sys: %s, alloc: %s, inuse: %s, objects: %s \n[Stack] sys: %s, inuse: %s",
		m.truncate(m.HeapSys), m.truncate(m.HeapAlloc), m.truncate(m.HeapInuse), m.truncate(m.HeapObjects), m.truncate(m.StackSys), m.truncate(m.StackInuse))
}

func (memStats) truncate(bytes uint64) string {
	unitStepIndex := 0
	bt := float64(bytes)
	for bt > unit {
		bt /= unit
		unitStepIndex++
	}

	return fmt.Sprintf("%.1f %s", bt, unitStep[unitStepIndex])
}
