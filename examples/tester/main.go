package main

import (
	"log"
	"math/rand"
	"strconv"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/layout"
)

func NewContentView() View {
	return &ContentView{}
}

type ContentView struct{}

func (v *ContentView) Body() SomeView {
	clr := Bind(NewColor(128))
	dis := Bind(false)

	rect := func(i int, clr *Binding[CGColor]) SomeView {
		return Text(strconv.Itoa(i)).
			Fill(clr).
			Frame(Bind(NewSize(50, 50))).
			Center().
			Disabled(dis)
	}

	red := Bind(NewColor(255, 0, 0))
	blue := Bind(NewColor(0, 0, 255))

	return ScrollView(
		VStack(
			Button("Test Button", func() {
				r := rand.Intn(256)
				g := rand.Intn(256)
				b := rand.Intn(256)
				clr.Set(NewColor(r, g, b))
			}),

			Rectangle().
				Fill(clr).
				Frame(Bind(NewSize(100, 100))).
				Disabled(dis),

			ScrollView(
				VStack(
					rect(1, red),
					rect(2, red),
					rect(3, red),
					rect(4, red),
					rect(5, red),
				).Debug(),
			).Frame(Bind(NewSize(200, 200))).
				ScrollViewDirection(Const(layout.DirectionVertical)).
				Debug(),

			ScrollView(
				HStack(
					rect(1, blue),
					rect(2, blue),
					rect(3, blue),
					rect(4, blue),
					rect(5, blue),
				).Debug(),
			).Frame(Bind(NewSize(200, 200))).
				ScrollViewDirection(Const(layout.DirectionHorizontal)).
				Debug(),

			Toggle(dis),
		).Spacing(),
	).Frame(Bind(NewSize(500, 500))).Align(Bind(layout.AlignLeading | layout.AlignTrailing)) //.Debug()
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(64, 16, 64, 64))
	app.SetWindowSize(800, 600)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(false)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}
}
