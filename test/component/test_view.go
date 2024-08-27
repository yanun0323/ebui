package component

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

type testView struct {
	title   string
	content string
}

func TestView(title, content string) ebui.View {
	return &testView{
		title:   title,
		content: content,
	}
}

func (view *testView) Body() ebui.SomeView {
	return ebui.VStack(
		ebui.HStack().BackgroundColor(color.RGBA{0, 0, 128, 128}),
		ebui.HStack(
			ebui.Spacer(),
			ebui.VStack().Frame(50, 50).BackgroundColor(color.Gray{128}),
			ebui.Button(func() {
				println("Hi")
			}, ebui.Text("Hello!")).
				BackgroundColor(color.RGBA{128, 0, 0, 128}).
				Padding(15, 15, 15, 15).
				BackgroundColor(color.RGBA{64, 0, 0, 64}).
				Frame(75, 50),
			ebui.Text("Hello!").
				BackgroundColor(color.Gray{200}).
				Padding(5, 5, 5, 5).
				BackgroundColor(color.Gray{50}),
			ebui.Spacer(),
		).BackgroundColor(color.RGBA{0, 128, 0, 128}),
		ebui.ZStack().BackgroundColor(color.RGBA{128, 0, 0, 128}),
	).
		BackgroundColor(color.White).
		// Padding(15, 15, 15, 15).
		Frame(400, 400).
		ForegroundColor(color.Black)
}
