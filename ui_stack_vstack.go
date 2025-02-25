package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...View) SomeView {
	vs := &vstackImpl{
		children: someViews(views...),
	}
	vs.ctx = newViewContext(vs)
	return vs
}

type vstackImpl struct {
	*ctx

	children []SomeView
}

func (v *vstackImpl) preload() (flexibleSize, CGInset, layoutFunc) {
	StackFormula := &formulaStack{
		types:    formulaVStack,
		stackCtx: v.ctx,
		children: v.children,
	}

	return StackFormula.preload()
}

func (v *vstackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := v.ctx.draw(screen, hook...)
	for _, child := range v.children {
		child.draw(screen)
	}

	return op
}
