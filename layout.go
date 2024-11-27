package ebui

import (
	"github.com/yanun0323/pkg/sys"
)

func clearCache(v SomeView) {
	v.clearCache()
	for _, child := range v.subView() {
		clearCache(child)
	}
}

func layout(v SomeView, pos point, container size) {
	childNoWidthCount := 0
	childNoHeightCount := 0
	children := v.subView()
	summedChildrenSize := size{}
	callbacks := make([]func(pos point, container size) size, 0, len(children))
	for i := range children {
		child := children[i]
		childSize := child.getSize()
		summedChildrenSize.w += sys.If(childSize.w >= 0, childSize.w, 0)
		summedChildrenSize.h += sys.If(childSize.h >= 0, childSize.h, 0)
		childNoWidthCount += sys.If(childSize.w >= 0, 0, 1)
		childNoHeightCount += sys.If(childSize.h >= 0, 0, 1)

		callbacks = append(callbacks, func(childPos point, childContainerDefault size) size {
			childFinalSize := layoutChildFinalSize(childSize, childContainerDefault)
			childFinalPos := layoutChildFinalPos(v, childPos, container, childFinalSize)
			child.setSize(childFinalSize)
			child.setPosition(childFinalPos)
			layout(child, childFinalPos, childFinalSize)
			return childFinalSize
		})
	}

	pos = layoutStartPos(v, pos, container, summedChildrenSize, childNoWidthCount, childNoHeightCount)
	childContainerDefault := layoutChildContainerDefault(container, summedChildrenSize, childNoWidthCount, childNoHeightCount)

	for _, callback := range callbacks {
		childSize := callback(pos, childContainerDefault)
		pos = v.stepSubView(pos, childSize)
	}
}

func layoutStartPos(v SomeView, pos point, container, summedChildrenSize size, childNoWidthCount, childNoHeightCount int) point {
	centerStartOffset := v.getStackSubViewStart(point{
		x: (container.w - summedChildrenSize.w) / 2,
		y: (container.h - summedChildrenSize.h) / 2,
	})
	centerStartOffset.x = sys.If(centerStartOffset.x < 0 || childNoWidthCount != 0, 0, centerStartOffset.x)
	centerStartOffset.y = sys.If(centerStartOffset.y < 0 || childNoHeightCount != 0, 0, centerStartOffset.y)
	return pos.Adds(centerStartOffset)
}

func layoutChildContainerDefault(container, summedChildrenSize size, childNoWidthCount, childNoHeightCount int) size {
	childNoWidthCount = sys.If(childNoWidthCount <= 0, 1, childNoWidthCount)
	childNoHeightCount = sys.If(childNoHeightCount <= 0, 1, childNoHeightCount)

	childContainer := container
	childContainer.w -= summedChildrenSize.w
	childContainer.h -= summedChildrenSize.h
	childContainer.w /= childNoWidthCount
	childContainer.h /= childNoHeightCount

	childContainer.w = sys.If(childContainer.w < 1, 1, childContainer.w)
	childContainer.h = sys.If(childContainer.h < 1, 1, childContainer.h)

	return childContainer
}

func layoutChildFinalSize(childSize, childContainerDefault size) size {
	w := sys.If(childSize.w >= 0, childSize.w, childContainerDefault.w)
	h := sys.If(childSize.h >= 0, childSize.h, childContainerDefault.h)
	return size{w, h}
}

func layoutChildFinalPos(parent SomeView, childPos point, parentContainer, childFinalSize size) point {
	centerOffset := point{
		x: (parentContainer.w - childFinalSize.w) / 2,
		y: (parentContainer.h - childFinalSize.h) / 2,
	}
	return childPos.Adds(parent.getStackSubViewCenterOffset(centerOffset))
}
