package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// HStack 水平排列
func HStack(views ...View) SomeView {
	hs := &hstackImpl{
		children: someViews(views...),
	}
	hs.viewCtx = newViewContext(hs)
	return hs
}

type hstackImpl struct {
	*viewCtx

	children []SomeView
}

func (h *hstackImpl) count() int {
	count := 1
	for _, child := range h.children {
		count += child.count()
	}
	return count
}

func (h *hstackImpl) preload(parent *viewCtxEnv) (flexibleSize, CGInset, layoutFunc) {
	StackFormula := &formulaStack{
		types:    formulaHStack,
		stackCtx: h.viewCtx,
		children: h.children,
	}

	return StackFormula.preload(parent)
}
func (h *hstackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := h.viewCtx.draw(screen, hook...)
	for _, child := range h.children {
		child.draw(screen, hook...)
	}

	return op
}
