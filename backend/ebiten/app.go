package ebiten

import (
	"errors"

	ebitenlib "github.com/hajimehoshi/ebiten/v2"
	ui "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/internal/runtime"
)

type Config struct {
	Title      string
	Width      int
	Height     int
	Background ui.Color
	Floating   bool
	FontPath   string
}

type Option func(*Config)

func WithWindowTitle(title string) Option {
	return func(config *Config) {
		config.Title = title
	}
}

func WithWindowSize(width, height int) Option {
	return func(config *Config) {
		config.Width = width
		config.Height = height
	}
}

func WithBackground(color ui.Color) Option {
	return func(config *Config) {
		config.Background = color
	}
}

func WithFloating(enabled bool) Option {
	return func(config *Config) {
		config.Floating = enabled
	}
}

func WithFontPath(path string) Option {
	return func(config *Config) {
		config.FontPath = path
	}
}

type App struct {
	config       Config
	runtime      *runtime.Runtime
	renderer     *Renderer
	mousePressed bool
	width        int
	height       int
}

func NewApp(root ui.View, options ...Option) (*App, error) {
	if root == nil {
		return nil, errors.New("root view is required")
	}

	config := Config{
		Title:      "EBUI v2",
		Width:      960,
		Height:     640,
		Background: ui.RGBA(240, 243, 248, 255),
	}
	for _, option := range options {
		option(&config)
	}

	renderer, err := NewRenderer(config.FontPath)
	if err != nil {
		return nil, err
	}

	app := &App{
		config:   config,
		renderer: renderer,
		runtime:  runtime.New(root, renderer),
		width:    config.Width,
		height:   config.Height,
	}
	app.runtime.Resize(ui.S(float64(config.Width), float64(config.Height)))
	return app, nil
}

func (app *App) Run() error {
	ebitenlib.SetWindowTitle(app.config.Title)
	ebitenlib.SetWindowFloating(app.config.Floating)
	ebitenlib.SetWindowSize(app.config.Width, app.config.Height)
	return ebitenlib.RunGame(app)
}

func (app *App) Update() error {
	x, y := ebitenlib.CursorPosition()
	pointer := ui.P(float64(x), float64(y))
	events := []runtime.Event{{Kind: runtime.EventPointerMove, Point: pointer}}

	pressed := ebitenlib.IsMouseButtonPressed(ebitenlib.MouseButtonLeft)
	if pressed && !app.mousePressed {
		events = append(events, runtime.Event{Kind: runtime.EventPointerDown, Point: pointer})
	}
	if !pressed && app.mousePressed {
		events = append(events, runtime.Event{Kind: runtime.EventPointerUp, Point: pointer})
	}
	app.mousePressed = pressed

	wx, wy := ebitenlib.Wheel()
	if wx != 0 || wy != 0 {
		events = append(events, runtime.Event{
			Kind:   runtime.EventWheel,
			Point:  pointer,
			Scroll: ui.P(wx, wy),
		})
	}

	app.runtime.Update(events)
	return nil
}

func (app *App) Draw(screen *ebitenlib.Image) {
	screen.Fill(app.config.Background.NRGBA())
	app.renderer.Draw(screen, app.runtime.Frame())
}

func (app *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	app.width = outsideWidth
	app.height = outsideHeight
	app.runtime.Resize(ui.S(float64(outsideWidth), float64(outsideHeight)))
	return outsideWidth, outsideHeight
}
