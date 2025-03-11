package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Spacer 元件
type spacerImpl struct {
	*viewCtx
}

func Spacer() SomeView {
	sp := &spacerImpl{}
	sp.viewCtx = newViewContext(sp)
	return sp
}

func (*spacerImpl) isSpacer() bool {
	return true
}

func (*spacerImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	return preloadData{}, nil
}

func (*spacerImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	// Spacer 是空白元件，不需要繪製任何內容
}
