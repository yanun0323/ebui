package main

import (
	"log"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/animation"
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
	bgColor := Bind(NewColor(255))
	toggle := Bind(false).AddListener(func(oldVal, newVal bool, animStyle ...animation.Style) {
		if newVal {
			bgColor.Set(NewColor(16), animStyle...)
		} else {
			bgColor.Set(NewColor(255), animStyle...)
		}
	})
	return VStack(
		Toggle(toggle),
		view.PageScrollView(),
	).
		Padding(Bind(NewInset(30))).
		Frame(Bind(NewSize(1000, 600))).
		Border(Bind(NewInset(1.5)), Bind(white)).
		RoundCorner(Bind(15.0)).
		Center().
		BackgroundColor(bgColor).
		FontKerning(Bind(1.0))

}

func main() {
	app := NewApplication(NewContentView())
	// app.SetWindowBackgroundColor(NewColor(0))
	app.SetWindowSize(1200, 800)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(true)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
