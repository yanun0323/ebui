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

func (e *emptyImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	return preloadData{}, func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect) {
		return CGRect{start, start}
	}
}

func (e *emptyImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	// EmptyView 是空白元件，不需要繪製任何內容
}
