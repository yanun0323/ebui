package ebui

type circleImpl struct {
	*viewCtx
}

func Circle() SomeView {
	circle := &circleImpl{}
	circle.viewCtx = newViewContext(circle)
	circle.viewCtx.RoundCorner(Const(Inf))
	return circle
}

func (c *circleImpl) userSetFrameSize() CGSize {
	frameSize := c.viewCtx.userSetFrameSize()
	frameSize.Width = min(frameSize.Width, frameSize.Height)
	frameSize.Height = frameSize.Width

	return frameSize
}

func (c *circleImpl) Frame(size *Binding[CGSize]) SomeView {
	if size == nil {
		size = Const(NewSize(Inf, Inf))
	}

	c.frameSize = BindOneWay(size, func(s CGSize) CGSize {
		sz := min(s.Width, s.Height)
		return NewSize(sz, sz)
	})

	return c
}

func (c *circleImpl) RoundCorner(radius ...*Binding[float64]) SomeView {
	return c
}
