package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ZStack(views ...View) SomeView {
	zs := &zstackImpl{
		children: someViews(views...),
	}
	zs.ctx = newViewContext(zs)
	return zs
}

type zstackImpl struct {
	*ctx

	children []SomeView
}

func (z *zstackImpl) preload() (flexibleSize, CGInset, layoutFunc) {
	StackFormula := &formulaStack{
		types:    formulaZStack,
		stackCtx: z.ctx,
		children: z.children,
	}

	return StackFormula.preload()
}
func (z *zstackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := z.ctx.draw(screen, hook...)
	for _, child := range z.children {
		child.draw(screen)
	}

	return op
}
