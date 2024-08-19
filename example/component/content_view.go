package component

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

type contentView struct {
	title   string
	content string
}

func ContentView(title, content string) ebui.View {
	return &contentView{
		title:   title,
		content: content,
	}
}

func (view *contentView) Body() ebui.SomeView {
	return ebui.HStack(
		ebui.Spacer(),
		ebui.VStack(
			ebui.Spacer(),
			ebui.Text(view.title).
				Padding(0, 15, 0, 15).
				ForegroundColor(color.White).
				BackgroundColor(color.Gray{128}).
				Frame(200, -1),
			ebui.Text(view.content),
			ebui.Spacer(),
		),
		ebui.Spacer(),
	).
		ForegroundColor(color.RGBA{200, 200, 200, 255}).
		BackgroundColor(color.RGBA{255, 0, 0, 255}).
		Padding(5, 5, 5, 5)
}

func (view *contentView) Body2() ebui.SomeView {
	return ebui.HStack(
		ebui.Spacer(),
		ebui.VStack(
			ebui.Spacer(),
			InfoView(func() { println("Info Clicked") }),
			ebui.Spacer(),
		),
		ebui.Text(&view.title).
			Padding(0, 15, 0, 15).
			ForegroundColor(color.Gray{128}).
			Padding(10, 10, 10, 10).
			BackgroundColor(color.White).
			Padding(30, 30, 30, 30).
			BackgroundColor(color.Gray{64}).
			Frame(200, -1),
		ebui.VStack(
			ebui.Text("VStack").
				BackgroundColor(color.Gray{32}),
			ebui.VStack(
				ebui.Text("HStack 1").
					BackgroundColor(color.Gray{160}),
				ebui.Text("HStack 2").
					Padding(5, 5, 5, 5).
					BackgroundColor(color.Gray{192}),
			),
		),
		ebui.Text("Trailing").
			BackgroundColor(color.Gray{16}),
		ebui.Spacer(),
	).
		ForegroundColor(color.RGBA{200, 200, 200, 255}).
		BackgroundColor(color.RGBA{255, 0, 0, 255}).
		Padding(5, 5, 5, 5)
}
