package component

import (
	"image/color"

	"github.com/yanun0323/ebui"
)

type infoView struct{}

func InfoView() ebui.View {
	return &infoView{}
}

func (view *infoView) Body() ebui.SomeView {
	return ebui.Text("O").
		Frame(50, 50).
		ForegroundColor(color.White).
		BackgroundColor(color.Gray{128}).
		CornerRadius(25)
}
