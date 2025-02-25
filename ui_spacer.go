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

func (*spacerImpl) preload(parent *viewCtxEnv) (flexibleSize, CGInset, layoutFunc) {
	return newFlexibleSize(Inf, Inf, true), NewInset(0, 0, 0, 0), nil
}

func (*spacerImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	// Spacer 是空白元件，不需要繪製任何內容
	return &ebiten.DrawImageOptions{}
}
