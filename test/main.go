package main

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

var (
	_blue   = color.RGBA{0, 0, 128, 128}
	_red    = color.RGBA{128, 0, 0, 128}
	_green  = color.RGBA{0, 128, 0, 128}
	_yellow = color.RGBA{128, 128, 0, 128}
)

func main() {
	contentView := ebui.VStack(
		ebui.HStack(
			ebui.VStack().
				Frame(25, 25).
				BackgroundColor(color.RGBA{0, 128, 0, 128}),
			ebui.VStack().
				Frame(25, 25).
				BackgroundColor(color.RGBA{0, 128, 128, 128}),
		).
			Frame(100, 100).
			BackgroundColor(_red),
		ebui.HStack(
			ebui.VStack().
				Frame(10, 10).
				BackgroundColor(_green),
			ebui.VStack().
				Frame(10, 10).
				BackgroundColor(_yellow),
		).
			Frame(50, 50).
			BackgroundColor(_green),
		ebui.HStack().
			Frame(25, 25).
			BackgroundColor(_yellow),
	).Frame(200, 300).BackgroundColor(_blue)
	ebui.Run("Windows Title", contentView, true)
}
