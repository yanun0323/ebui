package main

import (
	"context"
	"log"
	"sync/atomic"

	. "github.com/yanun0323/ebui"
)

var (
	white  = NewColor(255, 255, 255, 255)
	red    = NewColor(255, 0, 0, 255)
	blue   = NewColor(0, 0, 255, 255)
	green  = NewColor(0, 255, 0, 255)
	black  = NewColor(16, 16, 16, 255)
	yellow = NewColor(255, 255, 0, 255)

	borderInset = NewInset(1, 1, 1, 1)
)

func NewContentView() View {
	pauseString := "Random Color Pause"
	startString := "Random Color Start "

	return &ContentView{
		pauseString: pauseString,
		startString: startString,
		isRandom:    &atomic.Bool{},
		color:       Bind(red),
		content:     Bind(pauseString),
	}
}

type ContentView struct {
	pauseString string
	startString string

	isRandom  *atomic.Bool
	pauseFunc context.CancelFunc
	color     *Binding[CGColor]
	content   *Binding[string]
}

func (v *ContentView) Body() SomeView {
	println("ContentView.Body() called")

	return HStack(
		VStack(
			Rectangle().
				Frame(Bind(NewSize(100, 50))).
				RoundCorner(Bind(15.0)).
				BackgroundColor(Bind(red)).
				Border(Bind(borderInset), Bind(white)).
				Padding(Bind(NewInset(5, 5, 5, 5))),
			Rectangle().
				Frame(Bind(NewSize(100, 50))).
				Border(Bind(borderInset), Bind(white)).
				RoundCorner(Bind(15.0)).
				BackgroundColor(Bind(green)).
				Padding(Bind(NewInset(5, 5, 5, 5))),
			Rectangle().
				Frame(Bind(NewSize(100, 50))).
				BackgroundColor(Bind(blue)).
				RoundCorner(Bind(15.0)).
				Border(Bind(borderInset), Bind(white)).
				Padding(Bind(NewInset(5, 5, 5, 5))),
			Rectangle().
				Frame(Bind(NewSize(100, 50))).
				Border(Bind(borderInset), Bind(white)).
				Padding(Bind(NewInset(5, 5, 5, 5))),
			Rectangle().
				Frame(Bind(NewSize(100, 50))).
				Border(Bind(borderInset), Bind(white)).
				BackgroundColor(Bind(yellow)).
				RoundCorner(Bind(15.0)).
				Padding(Bind(NewInset(5, 5, 5, 5))),

			HStack(
				Spacer(),
				Circle().Frame(Bind(NewSize(20, 20))).BackgroundColor(Bind(red)),
			).Frame(Bind(NewSize(100, 20))).BackgroundColor(Bind(green)),
		),

		VStack(
			Spacer(),
			Circle().Frame(Bind(NewSize(20, 20))).BackgroundColor(Bind(blue)),
		).Frame(Bind(NewSize(20, 100))).BackgroundColor(Bind(green)),
	)
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(255, 32, 255, 32))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(false)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
