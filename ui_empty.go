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

func (e *emptyImpl) preload(parent *viewCtxEnv) (flexibleSize, CGInset, layoutFunc) {
	return flexibleSize{}, CGInset{}, nil
}

func (e *emptyImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{}
}
