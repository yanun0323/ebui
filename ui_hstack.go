package ebui

import "github.com/yanun0323/pkg/sys"

/* Check Interface Implementation */
var _ SomeView = (*hstackView)(nil)

func HStack(views ...View) SomeView {
	v := &hstackView{}
	v.uiView = newView(typesHStack, v, views...)
	return v
}

type hstackView struct {
	*uiView
}

func (p *hstackView) getSize() size {
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
	for _, child := range p.contents {
		childSize := child.getSize()
		result.w += sys.If(childSize.w >= 0, childSize.w, 0)
		result.h = max(result.h, childSize.h)
	}

	result.w += sys.If(result.w != -1, 1, 0)

	result.w = max(result.w, size.w)
	result.h = max(result.h, size.h)

	p.isCached = true
	p.cachedSize = result

	return result
}

func (p *hstackView) getStackSubViewStart(offset point) point {
	return point{offset.x, 0}
}

func (p *hstackView) getStackSubViewCenterOffset(offset point) point {
	return point{0, offset.y}
}

func (p *hstackView) stepSubView(pos point, childSize size) point {
	pos.x += childSize.w
	return pos
}
