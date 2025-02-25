package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Spacer 元件
type spacerImpl struct {
	*ctx
}

func Spacer() SomeView {
	sp := &spacerImpl{}
	sp.ctx = newViewContext(sp)
	return sp
}

func (*spacerImpl) preload() (flexibleSize, Inset, layoutFunc) {
	return newFlexibleSize(Inf, Inf, true), CGInset(0, 0, 0, 0), nil
}

func (*spacerImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	// Spacer 是空白元件，不需要繪製任何內容
	return &ebiten.DrawImageOptions{}
}
