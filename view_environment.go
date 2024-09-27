package ebui

import (
	"image/color"

	"github.com/yanun0323/ebui/font"
)

type uiViewEnvironment struct {
	fColor     color.Color
	fontSize   font.Size
	fontWeight font.Weight
	isPressing bool
}

func (p *uiViewEnvironment) set(anchor uiViewEnvironment) {
	p.fColor = rpZero(p.fColor, anchor.fColor)
	p.fontSize = rpZero(p.fontSize, anchor.fontSize)
	p.fontWeight = rpZero(p.fontWeight, anchor.fontWeight)
	p.isPressing = rpZero(p.isPressing, anchor.isPressing)
}
