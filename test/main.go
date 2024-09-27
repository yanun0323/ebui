package main

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

var (
	_blue  = color.RGBA{0, 0, 128, 128}
	_red   = color.RGBA{128, 0, 0, 128}
	_green = color.RGBA{0, 128, 0, 128}
)

func main() {
	contentView := ebui.VStack(
		ebui.HStack(
		// ebui.VStack().
		// 	Frame(50, 50).
		// 	BackgroundColor(color.RGBA{0, 128, 0, 128}),
		// ebui.VStack().
		// 	Frame(50, 50).
		// 	BackgroundColor(color.RGBA{0, 128, 128, 128}),
		).
			Frame(100, 100).
			BackgroundColor(_red),
		ebui.HStack().
			Frame(50, 50).
			BackgroundColor(_green),
	).BackgroundColor(_blue)
	ebui.Run("Windows Title", contentView, true)
}
