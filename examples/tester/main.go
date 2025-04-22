package main

import (
	"log"

	. "github.com/yanun0323/ebui"
)

var (
	gray   = NewColor(128)
	white  = NewColor(239)
	black  = NewColor(16)
	red    = NewColor(128, 0, 0)
	green  = NewColor(0, 128, 0)
	blue   = NewColor(0, 0, 128)
	yellow = NewColor(128, 128, 0)
)

func NewContentView() View {
	return &ContentView{
		lightMode: Bind(true).Animated(),
		length:    Bind(5.0).Animated(),
	}
}

type ContentView struct {
	lightMode *Binding[bool]
	length    *Binding[float64]
}

func (v *ContentView) Body() SomeView {
	return VStack(
		Rectangle().Fill(Const(green)),
		Rectangle().Fill(Const(red)),
	).Spacing().Padding()
}

func (v *ContentView) CurrentBackgroundColorText() SomeView {
	return Text(BindOneWay(v.lightMode, func(lightMode bool) string {
		if lightMode {
			return "White"
		}
		return "Black"
	}))
}

func (v *ContentView) ChangeBackgroundColor() {
	v.lightMode.Set(!v.lightMode.Get())
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowSize(300, 600)
	app.SetWindowBackgroundColor(black)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(true)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
