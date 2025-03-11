package ebui

import (
	"encoding/json"
	"fmt"
)

// viewCtxParam 提供所有 View 共用的參數
type viewCtxParam struct {
	_debug          string
	_systemSetFrame CGRect // 不包含 Padding 的內部邊界
	backgroundColor *Binding[AnyColor]
	frameSize       *Binding[CGSize]
	inset           *Binding[CGInset]
	roundCorner     *Binding[float64]
	borderInset     *Binding[CGInset]
	borderColor     *Binding[AnyColor]
	scaleToFit      *Binding[bool]
	keepAspectRatio *Binding[bool]
}

func newParam() *viewCtxParam {
	return &viewCtxParam{
		frameSize: Bind(NewSize(Inf, Inf)),
	}
}

func (p *viewCtxParam) userSetFrameSize() CGSize {
	return p.frameSize.Get()
}

// systemSetFrame 回傳的是內部邊界
func (p *viewCtxParam) systemSetFrame() CGRect {
	return p._systemSetFrame
}

// systemSetFrame 回傳的是外部邊界
func (p *viewCtxParam) systemSetBounds() CGRect {
	padding := p.padding()
	border := p.border()
	return NewRect(
		p._systemSetFrame.Start.X-padding.Left-border.Left,
		p._systemSetFrame.Start.Y-padding.Top-border.Top,
		p._systemSetFrame.End.X+padding.Right+border.Right,
		p._systemSetFrame.End.Y+padding.Bottom+border.Bottom,
	)
}

func (p *viewCtxParam) padding() CGInset {
	return p.inset.Get()
}

func (p *viewCtxParam) border() CGInset {
	return p.borderInset.Get()
}

func (p *viewCtxParam) debugPrint(prefix string, format string, a ...any) {
	if len(p._debug) != 0 {
		tag := fmt.Sprintf("%s \x1b[35m[%s]\x1b[0m\t", prefix, p._debug)
		logf(tag+format, a...)
	}
}

func (p *viewCtxParam) debugPrintPreload(frame CGRect, flexFrameSize CGSize, sData preloadData) {
	if len(p._debug) != 0 {
		tag := fmt.Sprintf("\x1b[35m[%s]\x1b[0m\t", p._debug)
		logf("preload %s\tStart(%4.f, %4.f)  End(%4.f, %4.f)  Size(%4.f, %4.f)  FlexSize(%4.f, %4.f), sData:\n%s",
			tag,
			frame.Start.X, frame.Start.Y,
			frame.End.X, frame.End.Y,
			frame.Dx(), frame.Dy(),
			flexFrameSize.Width, flexFrameSize.Height,
			serialize(sData),
		)
	}
}

func serialize(a any) string {
	s, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return fmt.Sprintf("RAW(%v)", a)
	}

	return string(s)
}
