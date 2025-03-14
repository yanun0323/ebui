package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...View) SomeView {
	return stack(formulaVStack, false, views...)
}

func HStack(views ...View) SomeView {
	return stack(formulaHStack, false, views...)
}

func ZStack(views ...View) SomeView {
	return stack(formulaZStack, false, views...)
}

func stack(types formulaType, flexibleStack bool, views ...View) SomeView {
	s := &stackImpl{
		types:         types,
		flexibleStack: flexibleStack,
		children:      someViews(views...),
	}
	s.viewCtx = newViewContext(s)
	return s
}

type stackImpl struct {
	*viewCtx

	types         formulaType
	flexibleStack bool
	children      []SomeView
}

func (s *stackImpl) count() int {
	count := 1
	for _, child := range s.children {
		count += child.count()
	}
	return count
}

func (s *stackImpl) preload(parent *viewCtxEnv, types ...formulaType) (preloadData, layoutFunc) {
	stackFormula := &formulaStack{
		types:                           s.types,
		stackCtx:                        s.viewCtx,
		children:                        s.children,
		ignorePreloadingChildSummedSize: s.flexibleStack,
	}

	return stackFormula.preload(parent, types...)
}
func (s *stackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	s.viewCtx.draw(screen, hook...)
	for _, child := range s.children {
		child.draw(screen, hook...)
	}
}
