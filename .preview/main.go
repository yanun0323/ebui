package main

import (
	preview "github.com/yanun0323/ebui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui"
)

func main() {

	ebiten.SetRunnableOnUnfocused(true)

	view := preview.Preview_Text()
	app := ebui.NewApplication(view)
	app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	app.Run("preview")
}
