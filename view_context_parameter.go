package ebui

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/yanun0323/ebui/internal/helper"
)

// viewCtxParam provides the parameters for all views
//
// It can NOT be inherited by subviews.
type viewCtxParam struct {
	_debug          string
	_systemSetFrame CGRect // the internal bounds without padding and border
	backgroundColor *Binding[CGColor]
	frameSize       *Binding[CGSize]
	inset           *Binding[CGInset]
	roundCorner     *Binding[float64]
	borderInset     *Binding[CGInset]
	borderColor     *Binding[CGColor]
	scaleToFit      *Binding[bool]
	keepAspectRatio *Binding[bool]
	offset          *Binding[CGPoint] /* only use for systemSetFrame */
	spacing         *Binding[float64]
}

func newParam() *viewCtxParam {
	return &viewCtxParam{
		frameSize: Bind(NewSize(Inf, Inf)),
	}
}

func (p *viewCtxParam) userSetFrameSize() CGSize {
	return p.frameSize.Get()
}

// systemSetFrame returns the internal bounds
func (p *viewCtxParam) systemSetFrame() CGRect {
	offset := p.offset.Get()
	return p._systemSetFrame.Move(offset)
}

// systemSetFrame returns the external bounds
func (p *viewCtxParam) systemSetBounds() CGRect {
	padding := p.padding()
	border := p.border()
	systemSetFrame := p.systemSetFrame()
	return NewRect(
		systemSetFrame.Start.X-padding.Left-border.Left,
		systemSetFrame.Start.Y-padding.Top-border.Top,
		systemSetFrame.End.X+padding.Right+border.Right,
		systemSetFrame.End.Y+padding.Bottom+border.Bottom,
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

/*
	##     ##    ###     ######  ##     ##    ###    ########  ##       ########
	##     ##   ## ##   ##    ## ##     ##   ## ##   ##     ## ##       ##
	##     ##  ##   ##  ##       ##     ##  ##   ##  ##     ## ##       ##
	######### ##     ##  ######  ######### ##     ## ########  ##       ######
	##     ## #########       ## ##     ## ######### ##     ## ##       ##
	##     ## ##     ## ##    ## ##     ## ##     ## ##     ## ##       ##
	##     ## ##     ##  ######  ##     ## ##     ## ########  ######## ########
*/

func (p *viewCtxParam) Bytes() []byte {
	b := bytes.Buffer{}
	b.Write(p._systemSetFrame.Bytes())
	b.Write(p.backgroundColor.Get().Bytes())
	b.Write(p.inset.Get().Bytes())
	b.Write(helper.BytesFloat64(p.roundCorner.Get()))
	b.Write(p.borderInset.Get().Bytes())
	b.Write(p.borderColor.Get().Bytes())
	b.Write(helper.BytesBool(p.scaleToFit.Get()))
	b.Write(helper.BytesBool(p.keepAspectRatio.Get()))
	b.Write(helper.BytesFloat64(p.spacing.Get()))

	return b.Bytes()
}
