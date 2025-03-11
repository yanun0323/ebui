package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...View) SomeView {
	vs := &vstackImpl{
		children: someViews(views...),
	}
	vs.viewCtx = newViewContext(vs)
	return vs
}

type vstackImpl struct {
	*viewCtx

	children []SomeView
}

func (v *vstackImpl) count() int {
	count := 1
	for _, child := range v.children {
		count += child.count()
	}
	return count
}

func (v *vstackImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	stackFormula := &formulaStack{
		types:    formulaVStack,
		stackCtx: v.viewCtx,
		children: v.children,
	}

	return stackFormula.preload(parent)
}

func (v *vstackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	v.viewCtx.draw(screen, hook...)
	for _, child := range v.children {
		child.draw(screen, hook...)
	}
}
