package main

import (
	"log"
	"math/rand"

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
				Frame(Bind(NewSize(100, 500))).
				Disabled(dis),

			Toggle(dis),
		).Spacing().Align(Bind(layout.AlignCenter)),
	).Frame(Bind(NewSize(500, 500))) //.Debug()
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
