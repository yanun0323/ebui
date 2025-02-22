package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...View) SomeView {
	vs := &vstackImpl{
		children: someViews(views...),
	}
	vs.ctx = newViewContext(tagVStack, vs)
	return vs
}

type vstackImpl struct {
	*ctx

	children []SomeView
}

func (v *vstackImpl) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	StackFormula := &formulaStack{
		stackCtx: v.ctx,
		children: v.children,
	}

	return StackFormula.preload()
}

func (v *vstackImpl) draw(screen *ebiten.Image, bounds ...CGRect) {
	v.ctx.draw(screen, bounds...)
	for _, child := range v.children {
		child.draw(screen)
	}
}
