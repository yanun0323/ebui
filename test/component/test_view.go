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
			ebui.Text("Hello!").Padding(5, 5, 5, 5).BackgroundColor(color.Gray{200}),
			ebui.Spacer(),
		).BackgroundColor(color.RGBA{0, 128, 0, 128}),
		ebui.ZStack().BackgroundColor(color.RGBA{128, 0, 0, 128}),
	).
		BackgroundColor(color.White).
		// Padding(15, 15, 15, 15).
		Frame(400, 400).
		ForegroundColor(color.Black)
}
