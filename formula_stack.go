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

func (v *formulaStack) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	stackEnv := v.stackCtx.inheritFrom(parent)
	childrenSummedBounds := CGSize{}
	childrenLayoutFns := make([]layoutFunc, 0, len(v.children))
	flexCount := NewPoint(0, 0)
	for _, child := range v.children {
		childData, layoutFn := child.preload(stackEnv)

		if child.isSpacer() {
			switch v.types {
			case formulaVStack:
				flexCount.Y++
				childrenLayoutFns = append(childrenLayoutFns, func(start CGPoint, flexFrameSize CGSize) (bounds CGRect) {
					return CGRect{start, NewPoint(start.X, start.Y+flexFrameSize.Height)}
				})
			case formulaHStack:
				flexCount.X++
				childrenLayoutFns = append(childrenLayoutFns, func(start CGPoint, flexFrameSize CGSize) (bounds CGRect) {
					return CGRect{start, NewPoint(start.X+flexFrameSize.Width, start.Y)}
				})
			}
		} else {
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
		{ // 計算子視圖大小總和及最小允許彈性邊界
			switch v.types {
			case formulaVStack:
				childrenSummedBounds = NewSize(max(childrenSummedBounds.Width, childBoundsSize.Width), childrenSummedBounds.Height+childBoundsSize.Height)
			case formulaHStack:
				childrenSummedBounds = NewSize(childrenSummedBounds.Width+childBoundsSize.Width, max(childrenSummedBounds.Height, childBoundsSize.Height))
			case formulaZStack:
				childrenSummedBounds = NewSize(max(childrenSummedBounds.Width, childBoundsSize.Width), max(childrenSummedBounds.Height, childBoundsSize.Height))
			}
		}

		if layoutFn != nil {
			childrenLayoutFns = append(childrenLayoutFns, layoutFn)
		}
	}

	sData, sLayoutFn := v.stackCtx.preload(stackEnv)
	{
		// 如果 Stack 本身沒有設定大小
		// 		-> 有彈性子視圖：使用無限大小
		// 		-> 沒有彈性子視圖：使用子視圖大小總和
		// 如果 Stack 本身有設定大小，則使用 Stack 的大小
		if sData.IsInfWidth {
			sData.FrameSize.Width = childrenSummedBounds.Width
			sData.IsInfWidth = flexCount.X > 0
		}

		if sData.IsInfHeight {
			sData.FrameSize.Height = childrenSummedBounds.Height
			sData.IsInfHeight = flexCount.Y > 0
		}
	}

	return sData, func(start CGPoint, flexBoundsSize CGSize) (bounds CGRect) {
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

		{ // 計算彈性大小公式
			switch v.types {
			case formulaVStack:
				hFlexSize := max(flexFrameSize.Height-childrenSummedBounds.Height, 0)
				perFlexFrameSize = NewSize(flexFrameSize.Width, hFlexSize/max(flexCount.Y, 1))
				ensureWidthMinimum()
			case formulaHStack:
				wFlexSize := max(flexFrameSize.Width-childrenSummedBounds.Width, 0)
				perFlexFrameSize = NewSize(wFlexSize/max(flexCount.X, 1), flexFrameSize.Height)
				ensureHeightMinimum()
			case formulaZStack:
				perFlexFrameSize = flexFrameSize

				ensureWidthMinimum()
				ensureHeightMinimum()
			}
		}

		anchor := start.
			Add(NewPoint(sData.Padding.Left, sData.Padding.Top)).
			Add(NewPoint(sData.Border.Left, sData.Border.Top))
		summedBoundsSize := NewSize(0, 0)

		for _, childLayoutFn := range childrenLayoutFns {
			childBounds := childLayoutFn(anchor, perFlexFrameSize)
			childBoundsSize := childBounds.Size()
			{ // 計算子視圖的 layout 最後位置
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

		finalFrameSize := summedBoundsSize
		finalFrame := sData.FrameSize
		{
			if sData.IsInfWidth {
				finalFrame.Width = finalFrameSize.Width
			}

			if sData.IsInfHeight {
				finalFrame.Height = finalFrameSize.Height
			}
		}

		finalBounds := finalFrame.Expand(sData.Padding).Expand(sData.Border)

		return sLayoutFn(start, finalBounds)
	}
}
