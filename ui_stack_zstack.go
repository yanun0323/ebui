package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ZStack(views ...View) SomeView {
	zs := &zstackImpl{
		children: someViews(views...),
	}
	zs.ctx = newViewContext(tagZStack, zs)
	return zs
}

type zstackImpl struct {
	*ctx

	children []SomeView
}

func (z *zstackImpl) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	StackFormula := &formulaStack{
		stackCtx: z.ctx,
		children: z.children,
	}

	return StackFormula.preload()
}
func (z *zstackImpl) draw(screen *ebiten.Image, bounds ...CGRect) {
	z.ctx.draw(screen, bounds...)
	for _, child := range z.children {
		child.draw(screen)
	}
}
