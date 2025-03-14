package main

import (
	"log"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/examples/view"
)

var (
	gray   = NewColor(128)
	white  = NewColor(255)
	black  = NewColor(0)
	red    = NewColor(128, 0, 0)
	green  = NewColor(0, 128, 0)
	blue   = NewColor(0, 0, 128)
	yellow = NewColor(128, 128, 0)
)

func NewContentView() View {
	return &ContentView{}
}

type ContentView struct{}

func (v *ContentView) Body() SomeView {
	return ZStack(
		view.LayoutView2(),
	).
		Padding(Bind(NewInset(30))).
		Frame(Bind(NewSize(1000, 600))).
		Border(Bind(NewInset(1.5)), Bind(white)).
		RoundCorner(Bind(15.0)).
		Center().
		FontKerning(Bind(1.0))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(64, 16, 64, 64))
	app.SetWindowSize(1200, 800)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(true)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
