package component

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

type actionView struct {
	action func()
}

func ActionView(clickedAction func()) ebui.View {
	return &actionView{
		action: clickedAction,
	}
}

func (view *actionView) Body() ebui.SomeView {
	return ebui.Button(
		view.action,
		ebui.Text("Info").
			ForegroundColor(color.White),
	).
		Frame(50, 50).
		BackgroundColor(color.Gray{128}).
		CornerRadius()
}
