package ebui

import "github.com/yanun0323/pkg/sys"

/* Check Interface Implementation */
var _ SomeView = (*vstackView)(nil)

func VStack(views ...View) SomeView {
	v := &vstackView{}
	v.uiView = newView(typesVStack, v, views...)
	return v
}

type vstackView struct {
	*uiView
}

func (p *vstackView) getSize() size {
	if p.isCached {
		return p.cachedSize
	}

	size := p.getFrame()
	if size.w != -1 && size.h != -1 {
		p.isCached = true
		p.cachedSize = size
		return size
	}

	result := _zeroSize
	childNoWidthCount := 0
	childNoHeightCount := 0
	for _, child := range p.contents {
		childSize := child.getSize()
		result.w = max(result.w, childSize.w)
		result.h += sys.If(childSize.h >= 0, childSize.h, 0)
		childNoWidthCount += sys.If(childSize.w >= 0 || child.getTypes() == typesSpacer, 0, 1)
		childNoHeightCount += sys.If(childSize.h >= 0, 0, 1)
	}

	result.h += sys.If(result.h != -1, 1, 0)

	result.w = max(result.w, size.w)
	result.h = max(result.h, size.h)
	result.w = sys.If(childNoWidthCount != 0, -1, result.w)
	result.h = sys.If(childNoHeightCount != 0, -1, result.h)

	p.isCached = true
	p.cachedSize = result

	return result
}

func (p *vstackView) getStackSubViewStart(offset point) point {
	return point{0, offset.y}
}

func (p *vstackView) getStackSubViewCenterOffset(offset point) point {
	return point{offset.x, 0}
}

func (p *vstackView) stepSubView(pos point, childSize size) point {
	pos.y += childSize.h
	return pos
}
