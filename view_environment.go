package ebui

import (
	"image/color"

	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/sys"
)

type uiViewEnvironment struct {
	fColor     color.Color
	fontSize   font.Size
	fontWeight font.Weight
	isPressing bool
}

func (p *uiViewEnvironment) set(anchor uiViewEnvironment) {
	p.fColor = sys.If(p.fColor == nil, anchor.fColor, p.fColor)
	p.fontSize = sys.If(p.fontSize == 0, anchor.fontSize, p.fontSize)
	p.fontWeight = sys.If(p.fontWeight == 0, anchor.fontWeight, p.fontWeight)
	p.isPressing = sys.If(!p.isPressing, anchor.isPressing, p.isPressing)
}
