package main

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/font"
)

var (
	_blue   = color.RGBA{0, 0, 128, 128}
	_red    = color.RGBA{128, 0, 0, 128}
	_green  = color.RGBA{0, 128, 0, 128}
	_yellow = color.RGBA{128, 128, 0, 128}
)

//go:embed *
var resource embed.FS

func main() {
	count := 0
	countStr := fmt.Sprintf("count: %d", count)

	leftView := ebui.VStack(
		ebui.HStack(
			ebui.VStack().
				Frame(25, 25).
				BackgroundColor(_green),
			ebui.VStack().
				Frame(25, 25).
				BackgroundColor(color.RGBA{0, 128, 128, 128}),
		).
			Frame(70, 70).
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
		ebui.Text("Hello, World!!!!!!").
			FontSize(font.Body).
			BackgroundColor(_green),
		ebui.Image("icon", resource).
			Resizable().
			AspectRatio().
			Frame(60, 60),
		ebui.Image("./test/icon").
			Resizable().
			Frame(80, 40),
		ebui.Text(&countStr).
			FontSize(font.Body).
			BackgroundColor(_red),
		ebui.Button(func() {
			count++
			countStr = fmt.Sprintf("count: %d", count)
		}, ebui.Text("Add Count").
			BackgroundColor(_green)).
			Frame(100, 40),
		ebui.Spacer(),
		ebui.HStack().
			Frame(25, 25).
			BackgroundColor(_yellow),
		ebui.ZStack(
			ebui.VStack().Frame(50, 50).BackgroundColor(_red),
			ebui.VStack().Frame(80, 80).BackgroundColor(_blue),
		),
	).BackgroundColor(_blue)

	rightView := ebui.VStack(
		ebui.Circle().ForegroundColor(_green),
	)

	contentView := ebui.HStack(
		leftView.Border(_red, 2),
		rightView.Frame(80, 80).Border(_blue, 2),
	)

	ebui.Run("Windows Title", contentView,
		ebui.WithDebug(),
		ebui.WithMemMonitor(),
	)
}
