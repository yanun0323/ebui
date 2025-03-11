package main

import (
	"log"

	. "github.com/yanun0323/ebui"
)

func NewContentView() View {
	return &ContentView{
		content: Bind(" "),
	}
}

type ContentView struct {
	content     *Binding[string]
	placeholder *Binding[string]
}

func (v *ContentView) Body() SomeView {
	return VStack(
		Spacer(),
		HStack(
			Spacer(),
			TextField(v.content, v.placeholder).
				Frame(Bind(NewSize(200, 50))),
			Spacer(),
		),
		Spacer(),
	)
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(32, 32, 32, 32))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(false)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
