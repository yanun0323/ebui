package ebui

import (
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/input"
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

func stack(types stackType, flexibleStack bool, views ...View) *stackImpl {
	s := &stackImpl{
		types:         types,
		flexibleStack: flexibleStack,
		children:      someViews(views...),
		baseCache:     newValue[*ebiten.Image](),
	}
	s.viewCtx = newViewContext(s)
	return s
}

type stackImpl struct {
	*viewCtx

	types         stackType
	flexibleStack bool
	children      []SomeView

	baseCache      *value[*ebiten.Image]
	childrenCached atomic.Bool
}

func (s *stackImpl) count() int {
	count := 1
	for _, child := range s.children {
		count += child.count()
	}
	return count
}

func (s *stackImpl) preload(parent *viewCtx, types ...stackType) (preloadData, layoutFunc) {
	stackFormula := &stackPreloader{
		types:                 s.types,
		stackCtx:              s.viewCtx,
		children:              s.children,
		preloadStackOnlyFrame: s.flexibleStack,
	}

	pd, lf := stackFormula.preload(parent, types...)
	return pd, func(start CGPoint, childBoundsSize CGSize) (CGRect, alignFunc, bool) {
		bounds, alignFunc, cached := lf(start, childBoundsSize)
		if !cached {
			s.childrenCached.Store(false)
		}
		return bounds, alignFunc, cached
	}
}
func (s *stackImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	s.viewCtx.draw(screen, hook...)

	var (
		sysFrame     = s.systemSetFrame()
		screenBounds = screen.Bounds()
		opt          = s.drawOption(NewRect(0), hook...)
	)

	cH := make([]func(*ebiten.DrawImageOptions), 0, len(hook)+1)
	cH = append(cH, hook...)
	cH = append(cH, func(op *ebiten.DrawImageOptions) {
		op.GeoM.Translate(-sysFrame.Start.X, -sysFrame.Start.Y)
	})

	var base *ebiten.Image

	childrenCached := s.childrenCached.Load()
	if childrenCached {
		base = s.baseCache.Load()
	} else {
		// base = ebiten.NewImage(sysFrame.Delta())
		base = ebiten.NewImage(screenBounds.Dx(), screenBounds.Dy())
		for _, child := range s.children {
			child.draw(base, cH...)
		}
	}

	if !childrenCached {
		s.baseCache.Store(base)
		s.childrenCached.Store(true)
	}

	screen.DrawImage(base, opt)
}

var _ eventHandler = &stackImpl{}

func (s *stackImpl) onAppearEvent() {
	s.viewCtx.onAppearEvent()
	for _, child := range s.children {
		child.onAppearEvent()
	}
}

func (s *stackImpl) onScrollEvent(cursor input.Vector, event input.ScrollEvent) bool {
	s.viewCtx.onScrollEvent(cursor, event)
	for _, child := range s.children {
		_ = child.onScrollEvent(cursor, event)
	}
	return false
}
func (s *stackImpl) onHoverEvent(cursor input.Vector) {
	s.viewCtx.onHoverEvent(cursor)
	for _, child := range s.children {
		child.onHoverEvent(cursor)
	}
}

func (s *stackImpl) onMouseEvent(event input.MouseEvent) {
	s.viewCtx.onMouseEvent(event)
	for _, child := range s.children {
		child.onMouseEvent(event)
	}
}

func (s *stackImpl) onKeyEvent(event input.KeyEvent) {
	s.viewCtx.onKeyEvent(event)
	for _, child := range s.children {
		child.onKeyEvent(event)
	}
}

func (s *stackImpl) onTypeEvent(event input.TypeEvent) {
	s.viewCtx.onTypeEvent(event)
	for _, child := range s.children {
		child.onTypeEvent(event)
	}
}

func (s *stackImpl) onGestureEvent(event input.GestureEvent) {
	s.viewCtx.onGestureEvent(event)
	for _, child := range s.children {
		child.onGestureEvent(event)
	}
}

func (s *stackImpl) onTouchEvent(event input.TouchEvent) {
	s.viewCtx.onTouchEvent(event)
	for _, child := range s.children {
		child.onTouchEvent(event)
	}
}

/*
	########   #######   ########   ##     ##  ##     ##  ##           ###
	##        ##     ##  ##     ##  ###   ###  ##     ##  ##          ## ##
	##        ##     ##  ##     ##  #### ####  ##     ##  ##         ##   ##
	######    ##     ##  ########   ## ### ##  ##     ##  ##        ##     ##
	##        ##     ##  ##   ##    ##     ##  ##     ##  ##        #########
	##        ##     ##  ##    ##   ##     ##  ##     ##  ##        ##     ##
	##         #######   ##     ##  ##     ##   #######   ########  ##     ##
*/

type stackType int

const (
	stackTypeZStack stackType = iota
	stackTypeVStack
	stackTypeHStack
)

type stackPreloader struct {
	types                 stackType
	stackCtx              *viewCtx
	children              []SomeView
	preloadStackOnlyFrame bool
}

func (v *stackPreloader) preload(parent *viewCtx, types ...stackType) (preloadData, layoutFunc) {
	var (
		children             = v.children
		childrenSummedBounds = CGSize{}
		childrenLayoutFns    = make([]layoutFunc, 0, len(children))
		flexCount            = NewPoint(0, 0)
	)

	v.stackCtx.inheritEnvFrom(parent)
	v.stackCtx.inheritStackParam(parent)

	t := v.types
	if t == stackTypeZStack && len(types) != 0 {
		t = types[0]
	}

	spacing := v.stackCtx.spacing.Value()
	isSpacingInf := isInf(spacing)

	for i, child := range children {
		child.inheritStackParam(v.stackCtx)
		childData, layoutFn := child.preload(v.stackCtx, t)
		{
			if childData.IsInfWidth {
				flexCount.X++
				if childData.FrameSize.IsInfWidth() {
					childData.FrameSize.Width = 0
				}
			}

			if childData.IsInfHeight {
				flexCount.Y++
				if childData.FrameSize.IsInfHeight() {
					childData.FrameSize.Height = 0
				}
			}
		}

		childBoundsSize := childData.BoundsSize()
		{ // calculate the summed size of the subviews and the minimum allowed flexible bounds
			switch v.types {
			case stackTypeVStack:
				if i != 0 {
					if isSpacingInf {
						flexCount.Y++
					} else {
						childrenSummedBounds.Height += spacing
					}
				}
				childrenSummedBounds = NewSize(max(childrenSummedBounds.Width, childBoundsSize.Width), childrenSummedBounds.Height+childBoundsSize.Height)
			case stackTypeHStack:
				if i != 0 {
					if isSpacingInf {
						flexCount.X++
					} else {
						childrenSummedBounds.Width += spacing
					}
				}
				childrenSummedBounds = NewSize(childrenSummedBounds.Width+childBoundsSize.Width, max(childrenSummedBounds.Height, childBoundsSize.Height))
			case stackTypeZStack:
				childrenSummedBounds = NewSize(max(childrenSummedBounds.Width, childBoundsSize.Width), max(childrenSummedBounds.Height, childBoundsSize.Height))
			}
		}

		childrenLayoutFns = append(childrenLayoutFns, layoutFn)
	}

	sData, sLayoutFn := v.stackCtx.preload(v.stackCtx)
	{
		// if the Stack itself has no size set
		// 		-> has flexible subviews: use infinite size
		// 		-> no flexible subviews: use the summed size of the subviews
		// if the Stack itself has a size set, use the Stack's size
		if v.preloadStackOnlyFrame {
			// do nothing
		} else {
			if sData.IsInfWidth {
				sData.FrameSize.Width = childrenSummedBounds.Width
				sData.IsInfWidth = flexCount.X > 0
			}

			if sData.IsInfHeight {
				sData.FrameSize.Height = childrenSummedBounds.Height
				sData.IsInfHeight = flexCount.Y > 0
			}
		}
	}

	return sData, func(start CGPoint, flexBoundsSize CGSize) (CGRect, alignFunc, bool) {
		childrenCached := true
		flexFrameSize := flexBoundsSize.Shrink(sData.Padding).Shrink(sData.Border)
		perFlexFrameSize := CGSize{}

		ensureWidthMinimum := func() {
			if !sData.IsInfWidth {
				perFlexFrameSize.Width = max(perFlexFrameSize.Width, sData.FrameSize.Width)
			}
		}

		ensureHeightMinimum := func() {
			if !sData.IsInfHeight {
				perFlexFrameSize.Height = max(perFlexFrameSize.Height, sData.FrameSize.Height)
			}
		}

		_, _ = ensureWidthMinimum, ensureHeightMinimum

		var (
			flexCountX                  = max(flexCount.X, 1)
			flexCountY                  = max(flexCount.Y, 1)
			isAlignHeaderVStackRequired = flexCount.Y == 0
			isAlignHeaderHStackRequired = flexCount.X == 0
		)

		{ // calculate the flexible size formula
			switch v.types {
			case stackTypeVStack:
				hFlexSize := sData.FrameSize.Height - childrenSummedBounds.Height
				if sData.IsInfHeight {
					hFlexSize = max(flexFrameSize.Height-childrenSummedBounds.Height, 0)
				}

				perFlexFrameSize = NewSize(flexFrameSize.Width, hFlexSize/flexCountY)
				ensureWidthMinimum()
			case stackTypeHStack:
				wFlexSize := sData.FrameSize.Width - childrenSummedBounds.Width
				if sData.IsInfWidth {
					wFlexSize = max(flexFrameSize.Width-childrenSummedBounds.Width, 0)
				}

				perFlexFrameSize = NewSize(wFlexSize/flexCountX, flexFrameSize.Height)
				ensureHeightMinimum()
			case stackTypeZStack:
				perFlexFrameSize = flexFrameSize
				ensureWidthMinimum()
				ensureHeightMinimum()
			}

			perFlexFrameSize.Width = max(perFlexFrameSize.Width, 0)
			perFlexFrameSize.Height = max(perFlexFrameSize.Height, 0)
		}

		if isSpacingInf {
			switch v.types {
			case stackTypeVStack:
				spacing = perFlexFrameSize.Height
			case stackTypeHStack:
				spacing = perFlexFrameSize.Width
			}
		}

		anchor := start.
			Add(NewPoint(sData.Padding.Left, sData.Padding.Top)).
			Add(NewPoint(sData.Border.Left, sData.Border.Top))
		summedBoundsSize := NewSize(0, 0)

		alignFuncs := make([]func(CGPoint), 0, len(childrenLayoutFns))
		aligners := make([]func(CGSize), 0, len(childrenLayoutFns))

		{ // align header
			a := v.stackCtx.transitionAlign.Value()
			switch v.types {
			case stackTypeVStack:
				if isAlignHeaderVStackRequired {
					delta := a.Y * perFlexFrameSize.Height
					anchor.Y += delta
					summedBoundsSize.Height += delta
				}
			case stackTypeHStack:
				if isAlignHeaderHStackRequired {
					delta := a.X * perFlexFrameSize.Width
					anchor.X += delta
					summedBoundsSize.Width += delta
				}
			}
		}

		for i, childLayoutFn := range childrenLayoutFns {
			if childLayoutFn == nil {
				continue
			}

			if i != 0 {
				switch v.types {
				case stackTypeVStack:
					anchor.Y += spacing
				case stackTypeHStack:
					anchor.X += spacing
				}
			}

			childBounds, alignChild, childCached := childLayoutFn(anchor, perFlexFrameSize)
			childrenCached = childrenCached && childCached

			alignFuncs = append(alignFuncs, alignChild)
			aligners = append(aligners, v.newAligner(childBounds, alignChild))
			childBoundsSize := childBounds.Size()
			{ // calculate the final position of the child view's layout
				switch v.types {
				case stackTypeVStack:
					anchor = NewPoint(anchor.X, childBounds.End.Y)
					summedBoundsSize = NewSize(
						max(summedBoundsSize.Width, childBoundsSize.Width),
						summedBoundsSize.Height+childBoundsSize.Height,
					)
				case stackTypeHStack:
					anchor = NewPoint(childBounds.End.X, anchor.Y)
					summedBoundsSize = NewSize(
						summedBoundsSize.Width+childBoundsSize.Width,
						max(summedBoundsSize.Height, childBoundsSize.Height),
					)
				case stackTypeZStack:
					summedBoundsSize = NewSize(
						max(summedBoundsSize.Width, childBoundsSize.Width),
						max(summedBoundsSize.Height, childBoundsSize.Height),
					)
				}
			}
		}

		var (
			sBounds    CGRect
			sAlignFunc alignFunc
		)

		if v.preloadStackOnlyFrame {
			sBounds, sAlignFunc, _ = v.layoutStack(sData, start, flexBoundsSize, sLayoutFn)
		} else {
			sBounds, sAlignFunc, _ = v.layoutStack(sData, start, summedBoundsSize, sLayoutFn)
		}
		sFrameSize := sBounds.Size().Shrink(sData.Padding).Shrink(sData.Border)
		for _, aligner := range aligners {
			aligner(sFrameSize)
		}

		return sBounds, func(offset CGPoint) {
			sAlignFunc(offset)
			for _, af := range alignFuncs {
				af(offset)
			}
		}, childrenCached
	}
}

func (v *stackPreloader) layoutStack(data preloadData, start CGPoint, finalFrameSize CGSize, layoutFn layoutFunc) (CGRect, alignFunc, bool) {
	finalFrame := data.FrameSize
	{
		if data.IsInfWidth || v.preloadStackOnlyFrame {
			finalFrame.Width = finalFrameSize.Width
		}

		if data.IsInfHeight || v.preloadStackOnlyFrame {
			finalFrame.Height = finalFrameSize.Height
		}
	}

	finalBounds := finalFrame.Expand(data.Padding).Expand(data.Border)

	return layoutFn(start, finalBounds)
}

func (v *stackPreloader) newAligner(childBounds CGRect, alignFunc alignFunc) func(stackFrame CGSize) {
	a := v.stackCtx.transitionAlign.Value()
	return func(stackFrame CGSize) {
		dw, dh := stackFrame.Width-childBounds.Dx(), stackFrame.Height-childBounds.Dy()
		offset := CGPoint{}
		offset.X = a.X * dw
		offset.Y = a.Y * dh

		switch v.types {
		case stackTypeVStack:
			offset.Y = 0
		case stackTypeHStack:
			offset.X = 0
		}

		alignFunc(offset)
	}
}
