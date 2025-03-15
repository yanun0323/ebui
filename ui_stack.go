package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func VStack(views ...View) SomeView {
	return stack(stackTypeVStack, false, views...)
}

func HStack(views ...View) SomeView {
	return stack(stackTypeHStack, false, views...)
}

func ZStack(views ...View) SomeView {
	return stack(stackTypeZStack, false, views...)
}

func stack(types stackType, flexibleStack bool, views ...View) SomeView {
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

	types         stackType
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

func (s *stackImpl) preload(parent *viewCtxEnv, types ...stackType) (preloadData, layoutFunc) {
	stackFormula := &stackPreloader{
		types:                 s.types,
		stackCtx:              s.viewCtx,
		children:              s.children,
		preloadStackOnlyFrame: s.flexibleStack,
	}

	return stackFormula.preload(parent, types...)
}
func (s *stackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	s.viewCtx.draw(screen, hook...)
	for _, child := range s.children {
		child.draw(screen, hook...)
	}
}
