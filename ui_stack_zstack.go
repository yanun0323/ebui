package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ZStack(views ...View) SomeView {
	return newZStackImpl(views...)
}

func newZStackImpl(views ...View) *zstackImpl {
	zs := &zstackImpl{
		children: someViews(views...),
	}
	zs.viewCtx = newViewContext(zs)
	return zs
}

type zstackImpl struct {
	*viewCtx

	children []SomeView
}

func (z *zstackImpl) count() int {
	count := 1
	for _, child := range z.children {
		count += child.count()
	}
	return count
}

func (z *zstackImpl) preload(parent *viewCtxEnv) (flexibleSize, CGInset, layoutFunc) {
	StackFormula := &formulaStack{
		types:    formulaZStack,
		stackCtx: z.viewCtx,
		children: z.children,
	}

	return StackFormula.preload(parent)
}
func (z *zstackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := z.viewCtx.draw(screen, hook...)
	for _, child := range z.children {
		child.draw(screen, hook...)
	}

	return op
}
