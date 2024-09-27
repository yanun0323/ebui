package component

import (
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
	return ebui.VStack()
	// return ebui.HStack(
	// 	ebui.Spacer(),
	// 	ebui.VStack(
	// 		ebui.Spacer(),
	// 		ebui.Text(view.title).
	// 			Padding(0, 15, 0, 15).
	// 			ForegroundColor(color.White).
	// 			BackgroundColor(color.Gray{128}),
	// 		ebui.Text(view.content),
	// 		ebui.Spacer(),
	// 	).Frame(200, -1),
	// 	ebui.Spacer(),
	// ).
	// 	ForegroundColor(color.RGBA{200, 200, 200, 255}).
	// 	BackgroundColor(color.RGBA{255, 0, 0, 255}).
	// 	Padding(5, 5, 5, 5)
}
