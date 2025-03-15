package main

import (
	"log"
	"strconv"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/layout"
)

func NewContentView() View {
	return &ContentView{}
}

type ContentView struct{}

func (v *ContentView) Body() SomeView {
	println("ContentView.Body() called")

	rect := func(i int) SomeView {
		return Text(Const(strconv.Itoa(i))).
			Align(Bind(layout.AlignCenter)).
			Frame(Bind(NewSize(100))).
			BackgroundColor(Bind(NewColor(128, 0, 0)))
	}

	enum := func(count int) []View {
		res := make([]View, 0, count)
		for i := range count {
			res = append(res, rect(i))
		}
		return res
	}

	return ScrollView(
		VStack(
			enum(10)...,
		).Spacing(),
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
