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
	opacity         *Binding[float64] /* only use for drawOption */
	disabled        *Binding[bool]    /* only use for drawOption */
	alignment       *Binding[layout.Align]
	transition      *Binding[float64] /* for all animation progress */
	transitionAlign *Binding[CGPoint] /* for all animation offset */
}

func newEnv() *viewCtxEnv {
	return &viewCtxEnv{
		foregroundColor: Bind(white),
	}
}

func (e *viewCtxEnv) inheritFrom(parent *viewCtxEnv) *viewCtxEnv {
	if parent == nil {
		return e
	}

	e.foregroundColor = getNewIfOldNil(parent.foregroundColor, e.foregroundColor)
	e.fontSize = getNewIfOldNil(parent.fontSize, e.fontSize)
	e.fontWeight = getNewIfOldNil(parent.fontWeight, e.fontWeight)
	e.fontLineHeight = getNewIfOldNil(parent.fontLineHeight, e.fontLineHeight)
	e.fontKerning = getNewIfOldNil(parent.fontKerning, e.fontKerning)
	e.fontAlignment = getNewIfOldNil(parent.fontAlignment, e.fontAlignment)
	e.fontItalic = getNewIfOldNil(parent.fontItalic, e.fontItalic)
	e.opacity = getNewIfOldNil(parent.opacity, e.opacity)
	e.disabled = getNewIfOldNil(parent.disabled, e.disabled)
	e.alignment = getNewIfOldNil(parent.alignment, e.alignment)
	e.transition = getNewIfOldNil(parent.transition, e.transition)
	e.transitionAlign = getNewIfOldNil(parent.transitionAlign, e.transitionAlign)

	return e
}

func getNewIfOldNil[T bindable](newValue, oldValue *Binding[T]) *Binding[T] {
	if oldValue != nil {
		return oldValue
	}

	return newValue
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
		b.Write(e.fontSize.Get().Bytes())
		b.Write(e.fontWeight.Get().Bytes())
		b.Write(helper.BytesFloat64(e.fontLineHeight.Get()))
		b.Write(helper.BytesFloat64(e.fontKerning.Get()))
		b.Write(e.fontAlignment.Get().Bytes())
		b.Write(helper.BytesBool(e.fontItalic.Get()))
	}

	b.Write(e.alignment.Get().Hash())

	b.Write(helper.BytesFloat64(e.transition.Get()))
	b.Write(e.transitionAlign.Get().Bytes())

	return b.Bytes()
}
