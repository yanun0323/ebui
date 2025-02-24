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

type Heap struct {
	data []int
}

func (h *Heap) Len() int {
	return len(h.data)
}
func (h *Heap) Less(i, j int) bool {
	return h.data[i] < h.data[j]
}
func (h *Heap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}
func (h *Heap) Push(x interface{}) {
	h.data = append(h.data, x.(int))
}
func (h *Heap) Pop() interface{} {
	x := h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	return x
}
