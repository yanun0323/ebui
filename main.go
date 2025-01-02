package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/your-project/ebui"
)

type ContentView struct {
	count ebui.Binding[int]
}

func (v *ContentView) Body() ebui.SomeView {
	return ebui.VStack(
		ebui.Text("Hello, Ebiten UI!"),
		ebui.Text(fmt.Sprintf("Count: %d", v.count.Get())),
		ebui.Button(func() {
			v.count.Set(v.count.Get() + 1)
		}, ebui.Text("Increment")),
	)
}

func main() {
	view := &ContentView{
		count: ebui.Binding[int]{},
	}

	app := &ebui.Application{
		StateManager: ebui.NewStateManager(image.Rect(0, 0, 640, 480)),
		EventManager: ebui.NewEventManager(),
		AnimManager:  ebui.NewAnimationManager(),
		RootView:     view.Body().Build(),
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SwiftUI-like Framework")
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
