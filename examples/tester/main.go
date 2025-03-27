package main

import (
	"log"

	. "github.com/yanun0323/ebui"
)

func NewContentView() View {
	return &ContentView{}
}

type ContentView struct{}

func (v *ContentView) Body() SomeView {
	return HStack(
		Spacer(),
		VStack(
			Spacer(),
			Text("Hello, World!").Debug(),
			Text("你好\n世界!\nThe third line").LineLimit(Const(1)).Debug(),
			Text("你好\n世界!\n行高: 10").FontLineHeight(Const(10.0)).Debug(),
			Text("你好\n世界!\n字距: 10").FontKerning(Const(10.0)).Debug(),
			Spacer(),
		),
		Spacer(),
	)
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
