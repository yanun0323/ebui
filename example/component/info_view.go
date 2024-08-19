package component

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

type infoView struct {
	action func()
}

func InfoView(clickedAction func()) ebui.View {
	return &infoView{
		action: clickedAction,
	}
}

func (view *infoView) Body() ebui.SomeView {
	return ebui.Button(
		view.action,
		ebui.Text("Info").
			ForegroundColor(color.White),
	).
		Frame(50, 50).
		BackgroundColor(color.Gray{128}).
		CornerRadius(25)
}
