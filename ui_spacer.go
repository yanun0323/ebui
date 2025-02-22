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
	sp.ctx = newViewContext(tagSpacer, sp)
	return sp
}

func (s *spacerImpl) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	return sz(Inf, Inf), ins(0, 0, 0, 0), func(start CGPoint, flexSize CGSize) CGRect {
		logf("[SPACER] %s start: %+v, flexSize: %+v\n", s.debug(), start, flexSize)
		return rect(start.X, start.Y, start.X+flexSize.Width, start.Y+flexSize.Height)
	}
}

func (s *spacerImpl) draw(screen *ebiten.Image, bounds ...CGRect) {
	// Spacer 是空白元件，不需要繪製任何內容
}
