package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// HStack 水平排列
func HStack(views ...View) SomeView {
	hs := &hstackImpl{
		children: someViews(views...),
	}
	hs.ctx = newViewContext(hs)
	return hs
}

type hstackImpl struct {
	*ctx

	children []SomeView
}

func (h *hstackImpl) preload() (flexibleCGSize, Inset, layoutFunc) {
	StackFormula := &formulaStack{
		types:    formulaHStack,
		stackCtx: h.ctx,
		children: h.children,
	}

	return StackFormula.preload()
}
func (h *hstackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := h.ctx.draw(screen, hook...)
	for _, child := range h.children {
		child.draw(screen)
	}

	return op
}
