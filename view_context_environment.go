package ebui

import (
	"bytes"

	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/internal/helper"
	"github.com/yanun0323/ebui/layout"
)

// viewCtxEnv provides the environment variables for all views
//
// It can be inherited by subviews.
type viewCtxEnv struct {
	foregroundColor *Binding[CGColor] /* only use for drawOption */
	fontSize        *Binding[font.Size]
	fontWeight      *Binding[font.Weight]
	fontLineHeight  *Binding[float64]
	fontKerning     *Binding[float64]
	fontAlignment   *Binding[font.TextAlign]
	fontItalic      *Binding[bool]
	lineLimit       *Binding[int]
	opacity         *Binding[float64] /* only use for drawOption */
	disabled        *Binding[bool]    /* only use for drawOption */
	alignment       *Binding[layout.Align]
	transition      *Binding[float64] /* for all animation progress */
	transitionAlign *Binding[CGPoint] /* for all animation offset */

	// ScrollView

	scrollViewDirection *Binding[layout.Direction]
}

func newEnv() *viewCtxEnv {
	return &viewCtxEnv{}
}

func (e *viewCtxEnv) initRootEnv() {
	e.foregroundColor = Bind(white)
}

func (e *viewCtxEnv) inheritEnvFrom(parent *viewCtx) {
	if parent == nil {
		return
	}

	e.foregroundColor = inheritValue(e.foregroundColor, parent.foregroundColor)
	e.fontSize = inheritValue(e.fontSize, parent.fontSize)
	e.fontWeight = inheritValue(e.fontWeight, parent.fontWeight)
	e.fontLineHeight = inheritValue(e.fontLineHeight, parent.fontLineHeight)
	e.fontKerning = inheritValue(e.fontKerning, parent.fontKerning)
	e.fontAlignment = inheritValue(e.fontAlignment, parent.fontAlignment)
	e.fontItalic = inheritValue(e.fontItalic, parent.fontItalic)
	e.lineLimit = inheritValue(e.lineLimit, parent.lineLimit)
	e.opacity = inheritValue(e.opacity, parent.opacity)
	e.disabled = inheritValue(e.disabled, parent.disabled)
	e.alignment = inheritValue(e.alignment, parent.alignment)
	e.transition = inheritValue(e.transition, parent.transition)
	e.transitionAlign = inheritValue(e.transitionAlign, parent.transitionAlign)
	e.scrollViewDirection = inheritValue(e.scrollViewDirection, parent.scrollViewDirection)
}

// inheritValue returns the value if it is not nil, otherwise it returns the parent value
func inheritValue[T bindable](nullableVal, val *Binding[T]) *Binding[T] {
	if nullableVal != nil {
		return nullableVal
	}

	return val
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

func (e *viewCtxEnv) Bytes(withFont bool) []byte {
	b := bytes.Buffer{}

	if withFont {
		b.Write(e.fontSize.Value().Bytes())
		b.Write(e.fontWeight.Value().Bytes())
		b.Write(helper.BytesFloat64(e.fontLineHeight.Value()))
		b.Write(helper.BytesFloat64(e.fontKerning.Value()))
		b.Write(e.fontAlignment.Value().Bytes())
		b.Write(helper.BytesBool(e.fontItalic.Value()))
		b.Write(helper.BytesInt(e.lineLimit.Value()))
	}

	b.Write(e.foregroundColor.Value().Bytes())
	b.Write(e.alignment.Value().Hash())
	b.Write(helper.BytesFloat64(e.transition.Value()))
	b.Write(e.transitionAlign.Value().Bytes())
	b.Write(helper.BytesInt8(e.scrollViewDirection.Value()))
	return b.Bytes()
}
