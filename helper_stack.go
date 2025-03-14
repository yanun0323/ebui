package ebui

type formulaType int

const (
	formulaVStack formulaType = iota
	formulaHStack
	formulaZStack
)

type formulaStack struct {
	types    formulaType
	stackCtx *viewCtx
	children []SomeView
}

func (v *formulaStack) childrenWithSpacing() []SomeView {
	spacing := v.stackCtx.spacing.Value()
	if spacing == 0 {
		return v.children
	}

	children := make([]SomeView, len(v.children)*2-1)
	for i, child := range v.children {
		children[i*2] = child
		if i != len(v.children)-1 {
			children[i*2+1] = spacingBlock(v.stackCtx.spacing)
		}
	}

	return children
}

const (
	_alignHeaderFooterCount = 1
)

func (v *formulaStack) preload(parent *viewCtxEnv, types ...formulaType) (preloadData, layoutFunc) {
	var (
		children             = v.children
		stackEnv             = v.stackCtx.inheritFrom(parent)
		childrenSummedBounds = CGSize{}
		childrenLayoutFns    = make([]layoutFunc, 0, len(children))
		flexCount            = NewPoint(0, 0)
	)

	t := v.types
	if t == formulaZStack && len(types) != 0 {
		t = types[0]
	}

	spacing := v.stackCtx.spacing.Value()
	isSpacingInf := isInf(spacing)

	for i, child := range children {
		childData, layoutFn := child.preload(stackEnv, t)

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
			case formulaVStack:
				if i != 0 {
					if isSpacingInf {
						flexCount.Y++
					} else {
						childrenSummedBounds.Height += spacing
					}
				}
				childrenSummedBounds = NewSize(max(childrenSummedBounds.Width, childBoundsSize.Width), childrenSummedBounds.Height+childBoundsSize.Height)
			case formulaHStack:
				if i != 0 {
					if isSpacingInf {
						flexCount.X++
					} else {
						childrenSummedBounds.Width += spacing
					}
				}
				childrenSummedBounds = NewSize(childrenSummedBounds.Width+childBoundsSize.Width, max(childrenSummedBounds.Height, childBoundsSize.Height))
			case formulaZStack:
				childrenSummedBounds = NewSize(max(childrenSummedBounds.Width, childBoundsSize.Width), max(childrenSummedBounds.Height, childBoundsSize.Height))
			}
		}

		childrenLayoutFns = append(childrenLayoutFns, layoutFn)
	}

	sData, sLayoutFn := v.stackCtx.preload(stackEnv)
	{
		// if the Stack itself has no size set
		// 		-> has flexible subviews: use infinite size
		// 		-> no flexible subviews: use the summed size of the subviews
		// if the Stack itself has a size set, use the Stack's size
		if sData.IsInfWidth {
			sData.FrameSize.Width = childrenSummedBounds.Width
			sData.IsInfWidth = flexCount.X > 0
		}

		if sData.IsInfHeight {
			sData.FrameSize.Height = childrenSummedBounds.Height
			sData.IsInfHeight = flexCount.Y > 0
		}
	}

	return sData, func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect, alignFunc alignFunc) {
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
			case formulaVStack:
				hFlexSize := sData.FrameSize.Height - childrenSummedBounds.Height
				if sData.IsInfHeight {
					hFlexSize = max(flexFrameSize.Height-childrenSummedBounds.Height, 0)
				}

				perFlexFrameSize = NewSize(flexFrameSize.Width, hFlexSize/flexCountY)
				ensureWidthMinimum()
			case formulaHStack:
				wFlexSize := sData.FrameSize.Width - childrenSummedBounds.Width
				if sData.IsInfWidth {
					wFlexSize = max(flexFrameSize.Width-childrenSummedBounds.Width, 0)
				}

				perFlexFrameSize = NewSize(wFlexSize/flexCountX, flexFrameSize.Height)
				ensureHeightMinimum()
			case formulaZStack:
				perFlexFrameSize = flexFrameSize
				ensureWidthMinimum()
				ensureHeightMinimum()
			}

			perFlexFrameSize.Width = max(perFlexFrameSize.Width, 0)
			perFlexFrameSize.Height = max(perFlexFrameSize.Height, 0)
		}

		if isSpacingInf {
			switch v.types {
			case formulaVStack:
				spacing = perFlexFrameSize.Height
			case formulaHStack:
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
			case formulaVStack:
				if isAlignHeaderVStackRequired {
					delta := a.Y * perFlexFrameSize.Height
					anchor.Y += delta
					summedBoundsSize.Height += delta
				}
			case formulaHStack:
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
				case formulaVStack:
					anchor.Y += spacing
				case formulaHStack:
					anchor.X += spacing
				}
			}

			childBounds, alignChild := childLayoutFn(anchor, perFlexFrameSize)
			alignFuncs = append(alignFuncs, alignChild)
			aligners = append(aligners, v.newAligner(childBounds, alignChild))
			childBoundsSize := childBounds.Size()
			{ // calculate the final position of the child view's layout
				switch v.types {
				case formulaVStack:
					anchor = NewPoint(anchor.X, childBounds.End.Y)
					summedBoundsSize = NewSize(
						max(summedBoundsSize.Width, childBoundsSize.Width),
						summedBoundsSize.Height+childBoundsSize.Height,
					)
				case formulaHStack:
					anchor = NewPoint(childBounds.End.X, anchor.Y)
					summedBoundsSize = NewSize(
						summedBoundsSize.Width+childBoundsSize.Width,
						max(summedBoundsSize.Height, childBoundsSize.Height),
					)
				case formulaZStack:
					summedBoundsSize = NewSize(
						max(summedBoundsSize.Width, childBoundsSize.Width),
						max(summedBoundsSize.Height, childBoundsSize.Height),
					)
				}
			}
		}

		sBounds, sAlignFunc := v.layoutStack(sData, start, summedBoundsSize, sLayoutFn)
		sFrameSize := sBounds.Size().Shrink(sData.Padding).Shrink(sData.Border)
		for _, aligner := range aligners {
			aligner(sFrameSize)
		}

		return sBounds, func(offset CGPoint) {
			sAlignFunc(offset)
			for _, af := range alignFuncs {
				af(offset)
			}
		}
	}
}

func (v *formulaStack) layoutStack(data preloadData, start CGPoint, summedBoundsSize CGSize, layoutFn layoutFunc) (bounds CGRect, alignFunc alignFunc) {
	finalFrameSize := summedBoundsSize
	finalFrame := data.FrameSize
	{
		if data.IsInfWidth {
			finalFrame.Width = finalFrameSize.Width
		}

		if data.IsInfHeight {
			finalFrame.Height = finalFrameSize.Height
		}
	}

	finalBounds := finalFrame.Expand(data.Padding).Expand(data.Border)

	return layoutFn(start, finalBounds)
}

func (v *formulaStack) newAligner(childBounds CGRect, alignFunc alignFunc) func(stackFrame CGSize) {
	a := v.stackCtx.transitionAlign.Value()
	return func(stackFrame CGSize) {
		dw, dh := stackFrame.Width-childBounds.Dx(), stackFrame.Height-childBounds.Dy()
		offset := CGPoint{}
		offset.X = a.X * dw
		offset.Y = a.Y * dh

		switch v.types {
		case formulaVStack:
			offset.Y = 0
		case formulaHStack:
			offset.X = 0
		}

		alignFunc(offset)
	}
}
