package ebui

import (
	"github.com/yanun0323/ebui/font"
)

// viewCtxEnv 提供所有 View 共用的環境變量
type viewCtxEnv struct {
	foregroundColor   *Binding[CGColor]
	fontSize          *Binding[font.Size]
	fontWeight        *Binding[font.Weight]
	fontLineHeight    *Binding[float64]
	fontLetterSpacing *Binding[float64]
	fontAlignment     *Binding[font.Alignment]
	fontItalic        *Binding[bool]
	opacity           *Binding[float64]
	disabled          *Binding[bool]
}

func newEnv() *viewCtxEnv {
	return &viewCtxEnv{
		foregroundColor: Bind[CGColor](white),
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
	e.fontLetterSpacing = getNewIfOldNil(parent.fontLetterSpacing, e.fontLetterSpacing)
	e.fontAlignment = getNewIfOldNil(parent.fontAlignment, e.fontAlignment)
	e.fontItalic = getNewIfOldNil(parent.fontItalic, e.fontItalic)
	e.opacity = getNewIfOldNil(parent.opacity, e.opacity)
	e.disabled = getNewIfOldNil(parent.disabled, e.disabled)

	return e
}

func getNewIfOldNil[T bindable](newValue, oldValue *Binding[T]) *Binding[T] {
	if oldValue != nil {
		return oldValue
	}

	return newValue
}
