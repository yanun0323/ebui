package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type spacerImpl struct {
	*viewCtx
}

// Spacer creates a blank component, so it does not need to draw anything.
func Spacer() SomeView {
	sp := &spacerImpl{}
	sp.viewCtx = newViewContext(sp)
	return sp
}

func (sp *spacerImpl) preload(parent *viewCtxEnv, stackTypes ...formulaType) (preloadData, layoutFunc) {
	types := getTypes(stackTypes...)

	sz := NewSize(Inf)
	switch types {
	case formulaVStack:
		sz.Width = 0
	case formulaHStack:
		sz.Height = 0
	}

	return newPreloadData(sz, CGInset{}, CGInset{}), func(start CGPoint, flexBoundsSize CGSize) (CGRect, alignFunc) {
		switch types {
		case formulaVStack:
			return CGRect{start, NewPoint(start.X, start.Y+flexBoundsSize.Height)}, func(CGPoint) {}
		case formulaHStack:
			return CGRect{start, NewPoint(start.X+flexBoundsSize.Width, start.Y)}, func(CGPoint) {}
		default:
			return CGRect{start, start}, func(CGPoint) {}
		}
	}
}

func (*spacerImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	// Spacer is a blank component, so it does not need to draw anything
}
