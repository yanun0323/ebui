package ebui

type formulaStack struct {
	stackCtx *ctx
	children []SomeView
}

func (v *formulaStack) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	summedSize := sz(0, 0)
	childrenLayoutFns := make([]func(CGPoint, CGSize) CGRect, 0, len(v.children))
	flexCount := pt(0, 0)
	for _, child := range v.children {
		childSize, childInset, layoutFn := child.preload()
		childBounds := childSize.Expand(childInset)
		if isInf(childBounds.Width) {
			flexCount.X++
			childBounds.Width = 0
		}
		if isInf(childBounds.Height) {
			flexCount.Y++
			childBounds.Height = 0
		}

		{ // 計算子視圖大小總和
			switch v.stackCtx._tag {
			case tagVStack:
				summedSize = sz(max(summedSize.Width, childBounds.Width), summedSize.Height+childBounds.Height)
			case tagHStack:
				summedSize = sz(summedSize.Width+childBounds.Width, max(summedSize.Height, childBounds.Height))
			case tagZStack:
				summedSize = sz(max(summedSize.Width, childBounds.Width), max(summedSize.Height, childBounds.Height))
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
		if isInf(sSize.Width) {
			if flexCount.X > 0 {
				sSize.Width = Inf
			} else {
				sSize.Width = summedSize.Width
			}
		}

		if isInf(sSize.Height) {
			if flexCount.Y > 0 {
				sSize.Height = Inf
			} else {
				sSize.Height = summedSize.Height
			}
		}
	}

	return sSize, sInset, func(start CGPoint, flexSize CGSize) CGRect {
		perFlexSize := flexSize.Shrink(sInset)
		{ // 計算彈性大小公式
			switch v.stackCtx._tag {
			case tagVStack:
				hFlexSize := max(perFlexSize.Height-summedSize.Height, 0)
				perFlexSize = sz(perFlexSize.Width, hFlexSize/max(flexCount.Y, 1))
			case tagHStack:
				wFlexSize := max(perFlexSize.Width-summedSize.Width, 0)
				perFlexSize = sz(wFlexSize/max(flexCount.X, 1), perFlexSize.Height)
			}
		}

		anchor := start.Add(pt(sInset.Left, sInset.Top))
		summedSize := sz(0, 0)

		for _, childLayoutFn := range childrenLayoutFns {
			childFrame := childLayoutFn(anchor, perFlexSize)
			childSize := childFrame.Size()
			{ // 計算子視圖的 layout 最後位置
				switch v.stackCtx._tag {
				case tagVStack:
					anchor = pt(anchor.X, childFrame.End.Y)
					summedSize = sz(
						max(summedSize.Width, childSize.Width),
						summedSize.Height+childSize.Height,
					)
				case tagHStack:
					anchor = pt(childFrame.End.X, anchor.Y)
					summedSize = sz(
						summedSize.Width+childSize.Width,
						max(summedSize.Height, childSize.Height),
					)
				case tagZStack:
					summedSize = summedSize.Add(childFrame.Size())
				}
			}
		}

		logf("[STACK] %s start: %+v, endSize: %+v, flexCount: %+v, flexSize: %+v\n", v.stackCtx.debug(), start, summedSize, flexCount, perFlexSize)
		return sLayoutFn(start, summedSize)
	}
}
