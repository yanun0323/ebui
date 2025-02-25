package ebui

import (
	"image/color"

	"github.com/yanun0323/ebui/font"
)

// viewCtxEnv 提供所有 View 共用的環境變量
type viewCtxEnv struct {
	foregroundColor   *Binding[AnyColor]
	fontSize          *Binding[font.Size]
	fontWeight        *Binding[font.Weight]
	fontLineHeight    *Binding[float64]
	fontLetterSpacing *Binding[float64]
	fontAlignment     *Binding[font.Alignment]
	fontItalic        *Binding[bool]
}

func newEnv() *viewCtxEnv {
	return &viewCtxEnv{
		foregroundColor: Bind[AnyColor](color.White),
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

	return e
}

func getNewIfOldNil[T comparable](newValue, oldValue *Binding[T]) *Binding[T] {
	if oldValue != nil {
		return oldValue
	}

	return newValue
}
