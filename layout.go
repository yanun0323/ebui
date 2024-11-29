package ebui

import (
	"github.com/yanun0323/pkg/sys"
)

// func reset(v SomeView) {
// 	v.reset()
// 	for _, child := range v.subView() {
// 		reset(child)
// 	}
// }

// func update(v SomeView) {
// 	v.update()
// 	for _, sv := range v.subView() {
// 		update(sv)
// 	}
// }

func layout(v SomeView, pos point, container size) {
	childNoWidthCount := 0
	childNoHeightCount := 0
	childHStackNoWidthCount := 0
	childVStackNoHeightCount := 0
	children := v.subView()
	summedChildrenSize := size{}
	callbacks := make([]func(point, size) size, 0, len(children))
	for i := range children {
		child := children[i]
		childSize := child.getSize()
		childTypes := child.getTypes()
		summedChildrenSize.w += sys.If(childSize.w >= 0, childSize.w, 0)
		summedChildrenSize.h += sys.If(childSize.h >= 0, childSize.h, 0)
		childNoWidthCount += sys.If(childSize.w < 0, 1, 0)
		childNoHeightCount += sys.If(childSize.h < 0, 1, 0)
		// childHStackNoWidthCount += sys.If(childSize.w < 0 && childTypes == typesHStack, 1, 0)
		// childVStackNoHeightCount += sys.If(childSize.h < 0 && childTypes == typesVStack, 1, 0)

		callbacks = append(callbacks, func(childPos point, childContainerDefault size) size {
			childFinalSize := layoutChildFinalSize(childSize, childContainerDefault, childNoWidthCount, childNoHeightCount, v.getTypes(), childTypes)
			childFinalPos := layoutChildFinalPos(v, childPos, container, childFinalSize)
			child.setSize(childFinalSize)
			child.setPosition(childFinalPos)
			layout(child, childFinalPos, childFinalSize)
			return childFinalSize
		})
	}

	pos = layoutStartPos(v, pos, container, summedChildrenSize, childNoWidthCount, childNoHeightCount, childHStackNoWidthCount, childVStackNoHeightCount)
	childContainerDefault := layoutChildContainerDefault(container, summedChildrenSize, childNoWidthCount, childNoHeightCount)

	for _, callback := range callbacks {
		childSize := callback(pos, childContainerDefault)
		pos = v.stepSubView(pos, childSize)
	}
}

func layoutStartPos(v SomeView, pos point, container, summedChildrenSize size, childNoWidthCount, childNoHeightCount, childHStackNoWidthCount, childVStackNoHeightCount int) point {
	centerStartOffset := v.getStackSubViewStart(point{
		x: (container.w - summedChildrenSize.w) / 2,
		y: (container.h - summedChildrenSize.h) / 2,
	})

	// TODO: 判斷 types VStack/HStack 若有 NoWidthCount/NoHeightCount  則起始點為 0
	// switch v.getTypes() {
	// case typesVStack:
	// 	if childNoWidthCount != 0 {
	// 		centerStartOffset.x = 0
	// 	}
	// case typesHStack:
	// 	if childNoHeightCount != 0 {
	// 		centerStartOffset.y = 0
	// 	}
	// }

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

func layoutChildFinalSize(childSize, childContainerDefault size, childNoWidthCount, childNoHeightCount int, parentTypes, childTypes types) size {
	w := sys.If(childSize.w >= 0, childSize.w, childContainerDefault.w)
	h := sys.If(childSize.h >= 0, childSize.h, childContainerDefault.h)

	// switch parentTypes {
	// case typesVStack:
	// 	if childTypes == typesHStack && childNoWidthCount != 0 {
	// 		w = max(w, childContainerDefault.w)
	// 	}
	// case typesHStack:
	// 	if childTypes == typesVStack && childNoHeightCount != 0 {
	// 		logs.Info(h)
	// 		h = max(h, childContainerDefault.h)
	// 	}
	// }

	return size{w, h}
}

func layoutChildFinalPos(parent SomeView, childPos point, parentContainer, childFinalSize size) point {
	centerOffset := point{
		x: (parentContainer.w - childFinalSize.w) / 2,
		y: (parentContainer.h - childFinalSize.h) / 2,
	}

	return childPos.Adds(parent.getStackSubViewOffsetToCenter(centerOffset))
}
