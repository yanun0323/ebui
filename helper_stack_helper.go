package ebui

import "github.com/hajimehoshi/ebiten/v2"

func getTypes(types ...formulaType) formulaType {
	if len(types) == 0 {
		return formulaZStack
	}
	return types[0]
}

type spacingBlockImpl struct {
	*viewCtx
	spacing *Binding[float64]
}

func spacingBlock(spacing *Binding[float64]) SomeView {
	sp := &spacingBlockImpl{
		spacing: spacing,
	}
	sp.viewCtx = newViewContext(sp)
	return sp
}

func (v *spacingBlockImpl) preload(parent *viewCtxEnv, stackTypes ...formulaType) (preloadData, layoutFunc) {
	spacing := v.spacing.Get()
	types := getTypes(stackTypes...)
	data := newPreloadData(NewSize(spacing), CGInset{}, CGInset{})
	switch types {
	case formulaVStack:
		data.FrameSize.Width = 0
		data.IsInfWidth = false
	case formulaHStack:
		data.FrameSize.Height = 0
		data.IsInfHeight = false
	}

	return data, func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect, alignFunc alignFunc) {
		if !isInf(spacing) {
			flexBoundsSize.Width = spacing
			flexBoundsSize.Height = spacing
		}

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

func (v *spacingBlockImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	// SpacingBlock is a blank component, so it does not need to draw anything
}
