package main

import (
	"log"
	"sync/atomic"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/layout"
)

var (
	white  = NewColor(255)
	red    = NewColor(255, 0, 0)
	blue   = NewColor(0, 0, 255)
	green  = NewColor(0, 255, 0)
	yellow = NewColor(255, 255, 0)

	borderInset = NewInset(1)
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

	isRandom *atomic.Bool
	color    *Binding[CGColor]
	content  *Binding[string]
}

func (v *ContentView) Body() SomeView {
	println("ContentView.Body() called")

	return VStack(
		HStack(
			VStack(
				Rectangle().
					Frame(Bind(NewSize(80, 40))).
					RoundCorner(Bind(15.0)).
					BackgroundColor(Bind(red)).
					Border(Bind(borderInset), Bind(white)).
					Padding(Bind(NewInset(5))),
				Rectangle().
					Frame(Bind(NewSize(100, 50))).
					Border(Bind(borderInset), Bind(white)).
					RoundCorner(Bind(15.0)).
					BackgroundColor(Bind(green)).
					Padding(Bind(NewInset(5))),
			).Modify(debugFunc).
				Align(Bind(layout.AlignLeading)),

			VStack(
				Rectangle().
					Frame(Bind(NewSize(100, 80))).
					BackgroundColor(Bind(blue)).
					RoundCorner(Bind(15.0)).
					Border(Bind(borderInset), Bind(white)).
					Padding(Bind(NewInset(5))),
				Rectangle().
					Frame(Bind(NewSize(80, 40))).
					Border(Bind(borderInset), Bind(white)).
					Padding(Bind(NewInset(5))),
			).Modify(debugFunc).
				Align(Bind(layout.AlignCenter)),

			VStack(
				Rectangle().
					Frame(Bind(NewSize(120, 50))).
					Border(Bind(borderInset), Bind(white)).
					BackgroundColor(Bind(yellow)).
					RoundCorner(Bind(15.0)).
					Padding(Bind(NewInset(5))),
				Rectangle().
					Frame(Bind(NewSize(100, 50))).
					Border(Bind(CGInset{}), Bind(white)).
					BackgroundColor(Bind(green)).
					RoundCorner(Bind(15.0)).
					Padding(Bind(NewInset(5))),
			).Modify(debugFunc).
				Align(Bind(layout.AlignTrailing)),
		),

		HStack(
			rect(50, 20, red),
			rect(50, 40, green),
			rect(50, 60, blue),
		).Modify(debugFunc).
			Align(Bind(layout.AlignTop)),

		HStack(
			rect(50, 20, red),
			rect(50, 40, green),
			rect(50, 60, blue),
		).Modify(debugFunc).
			Align(Bind(layout.AlignCenter)),

		HStack(
			rect(50, 20, red),
			rect(50, 40, green),
			rect(50, 60, blue),
		).Modify(debugFunc).
			Align(Bind(layout.AlignBottom)),
	)
}

func rect(w, h int, color CGColor) SomeView {
	sli := []int{20, 40, 60, 80, 100}
	sz := Bind(NewSize(w, h))
	return Button("rect", func() {
		WithAnimation(func() {
			h := int(sz.Get().Height)
			idx := 0
			for i, s := range sli {
				if h == s {
					idx = i
					break
				}
			}
			idx = (idx + 1) % len(sli)
			sz.Set(NewSize(w, sli[idx]))
		})
	}, func() SomeView {
		return Rectangle().
			Frame(sz).
			BackgroundColor(Bind(color)).
			RoundCorner(Bind(15.0)).
			Padding(Bind(NewInset(5)))
	})
}

func debugFunc(view SomeView) SomeView {
	return view.
		Padding(Bind(NewInset(5))).
		Border(Bind(NewInset(1)), Bind(red)).
		Padding(Bind(NewInset(5)))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(64, 16, 64, 64))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(false)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
