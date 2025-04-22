package main

import (
	"github.com/yanun0323/ebui"
	preview "github.com/yanun0323/ebui"
)

func main() {
	view := preview.Preview_Text()
	app := ebui.NewApplication(view)
	app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	app.Run("preview")
}
