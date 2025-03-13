package ebui

import "github.com/hajimehoshi/ebiten/v2"

type emptyImpl struct {
	*viewCtx
}

func EmptyView() SomeView {
	empty := &emptyImpl{}
	empty.viewCtx = newViewContext(empty)
	return empty
}

func (e *emptyImpl) preload(parent *viewCtxEnv, _ ...formulaType) (preloadData, layoutFunc) {
	return preloadData{}, func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect, alignFunc alignFunc) {
		return CGRect{start, start}, func(CGPoint) {}
	}
}

func (e *emptyImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	// EmptyView is a blank component, so it does not need to draw anything
}
