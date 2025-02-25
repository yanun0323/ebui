package ebui

// viewCtxParam 提供所有 View 共用的參數
type viewCtxParam struct {
	_debug          string
	_systemSetFrame CGRect // 不包含 Padding 的內部邊界
	backgroundColor *Binding[AnyColor]
	frameSize       *Binding[CGSize]
	inset           *Binding[CGInset]
	roundCorner     *Binding[float64]

	scaleToFit      *Binding[bool]
	keepAspectRatio *Binding[bool]
}

func newParam() *viewCtxParam {
	return &viewCtxParam{
		frameSize: Bind(NewSize(Inf, Inf)),
	}
}

func (p *viewCtxParam) userSetFrameSize() flexibleSize {
	frame := p.frameSize.Get()
	return newFlexibleSize(frame.Width, frame.Height)
}

// systemSetFrame 回傳的是內部邊界
func (p *viewCtxParam) systemSetFrame() CGRect {
	return p._systemSetFrame
}

// systemSetFrame 回傳的是外部邊界
func (p *viewCtxParam) systemSetBounds() CGRect {
	inset := p.inset.Get()
	return NewRect(
		p._systemSetFrame.Start.X-inset.Left,
		p._systemSetFrame.Start.Y-inset.Top,
		p._systemSetFrame.End.X+inset.Right,
		p._systemSetFrame.End.Y+inset.Bottom,
	)
}

func (p *viewCtxParam) debugPrint(frame CGRect) {
	if len(p._debug) != 0 {
		logf("\x1b[35m[%s]\x1b[0m\tStart(%4.f, %4.f)\tEnd(%4.f, %4.f)\tSize(%4.f, %4.f)",
			p._debug,
			frame.Start.X, frame.Start.Y,
			frame.End.X, frame.End.Y,
			frame.Dx(), frame.Dy(),
		)
	}
}
