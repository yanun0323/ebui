package ebui

type formulaType int

const (
	formulaVStack formulaType = iota
	formulaHStack
	formulaZStack
)

type formulaStack struct {
	types    formulaType
	stackCtx *ctx
	children []SomeView
}

func (v *formulaStack) preload() (flexibleSize, Inset, layoutFunc) {
	childrenSummedBounds := CGSize(0, 0)
	childrenLayoutFns := make([]layoutFunc, 0, len(v.children))
	flexCount := CGPoint(0, 0)
	for _, child := range v.children {
		childFrame, childInset, layoutFn := child.preload()
		childFrameSize := childFrame.Frame

		if childFrame.IsSpacer {
			switch v.types {
			case formulaVStack:
				flexCount.Y++
				childrenLayoutFns = append(childrenLayoutFns, func(start Point, flexFrameSize Size) (bounds Rect) {
					return Rect{start, CGPoint(start.X, start.Y+flexFrameSize.Height)}
				})
			case formulaHStack:
				flexCount.X++
				childrenLayoutFns = append(childrenLayoutFns, func(start Point, flexFrameSize Size) (bounds Rect) {
					return Rect{start, CGPoint(start.X+flexFrameSize.Width, start.Y)}
				})
			}
		} else {
			if childFrame.IsInfX {
				flexCount.X++
				childFrameSize.Width = 0
			}

			if childFrame.IsInfY {
				flexCount.Y++
				childFrameSize.Height = 0
			}
		}

		childBoundsSize := childFrameSize.Expand(childInset)
		{ // 計算子視圖大小總和
			switch v.types {
			case formulaVStack:
				childrenSummedBounds = CGSize(max(childrenSummedBounds.Width, childBoundsSize.Width), childrenSummedBounds.Height+childBoundsSize.Height)
			case formulaHStack:
				childrenSummedBounds = CGSize(childrenSummedBounds.Width+childBoundsSize.Width, max(childrenSummedBounds.Height, childBoundsSize.Height))
			case formulaZStack:
				childrenSummedBounds = CGSize(max(childrenSummedBounds.Width, childBoundsSize.Width), max(childrenSummedBounds.Height, childBoundsSize.Height))
			}
		}

		if layoutFn != nil {
			childrenLayoutFns = append(childrenLayoutFns, layoutFn)
		}
	}

	sSize, sInset, sLayoutFn := v.stackCtx.preload()
	{
		// 如果 Stack 本身沒有設定大小
		// 		-> 有彈性子視圖：使用無限大小
		// 		-> 沒有彈性子視圖：使用子視圖大小總和
		// 如果 Stack 本身有設定大小，則使用 Stack 的大小
		if sSize.IsInfX {
			sSize.Frame.Width = childrenSummedBounds.Width
			sSize.IsInfX = flexCount.X > 0
		}

		if sSize.IsInfY {
			sSize.Frame.Height = childrenSummedBounds.Height
			sSize.IsInfY = flexCount.Y > 0
		}
	}

	return sSize, sInset, func(start Point, flexFrameSize Size) (bounds Rect) {
		perFlexFrameSize := flexFrameSize.Shrink(sInset)
		{ // 計算彈性大小公式
			switch v.types {
			case formulaVStack:
				hFlexSize := max(perFlexFrameSize.Height-childrenSummedBounds.Height, 0)
				perFlexFrameSize = CGSize(perFlexFrameSize.Width, hFlexSize/max(flexCount.Y, 1))
			case formulaHStack:
				wFlexSize := max(perFlexFrameSize.Width-childrenSummedBounds.Width, 0)
				perFlexFrameSize = CGSize(wFlexSize/max(flexCount.X, 1), perFlexFrameSize.Height)
			}
		}

		anchor := start.Add(CGPoint(sInset.Left, sInset.Top))
		summedSize := CGSize(0, 0)

		for _, childLayoutFn := range childrenLayoutFns {
			childFrame := childLayoutFn(anchor, perFlexFrameSize)
			childSize := childFrame.Size()
			{ // 計算子視圖的 layout 最後位置
				switch v.types {
				case formulaVStack:
					anchor = CGPoint(anchor.X, childFrame.End.Y)
					summedSize = CGSize(
						max(summedSize.Width, childSize.Width),
						summedSize.Height+childSize.Height,
					)
				case formulaHStack:
					anchor = CGPoint(childFrame.End.X, anchor.Y)
					summedSize = CGSize(
						summedSize.Width+childSize.Width,
						max(summedSize.Height, childSize.Height),
					)
				case formulaZStack:
					summedSize = summedSize.Add(childFrame.Size())
				}
			}
		}

		finalFlexibleBounds := flexFrameSize.Shrink(sInset)
		finalBound := sSize.Frame
		{
			if sSize.IsInfX {
				finalBound.Width = finalFlexibleBounds.Width
			}

			if sSize.IsInfY {
				finalBound.Height = finalFlexibleBounds.Height
			}
		}

		return sLayoutFn(start, finalBound)
	}
}
