package ebui

import (
	"image/color"

	"github.com/yanun0323/ebui/font"
	gol "github.com/yanun0323/gollection/v2"
)

type viewParameter struct {
	id              identity
	offsets         gol.Queue[Binding[Point]]
	opacity         gol.Queue[Binding[float64]]
	frameSize       Binding[Size]
	foregroundColor Binding[color.Color]
	backgroundColor Binding[color.Color]
	disable         Binding[bool]
	hidden          Binding[bool]
	padding         []Binding[int]
	fontSize        Binding[font.Size]
	fontWeight      Binding[font.Weight]
	fontLineSpacing Binding[float64]
	fontKerning     Binding[int]
	fontItalic      Binding[bool]
}

func newViewParameter(id identity) viewParameter {
	return viewParameter{
		id:              id,
		offsets:         gol.NewQueue[Binding[Point]](),
		opacity:         gol.NewQueue[Binding[float64]](),
		frameSize:       New[Size](Size{W: 0, H: 0}),
		foregroundColor: New[color.Color](color.White),
		backgroundColor: New[color.Color](color.Black),
		disable:         New(false),
		padding:         []Binding[int]{New(15), New(15), New(15), New(15)},
		fontSize:        New(font.Body),
		fontWeight:      New(font.Normal),
		fontLineSpacing: New(0.0),
		fontKerning:     New(0),
		fontItalic:      New(false),
	}
}

func (param viewParameter) ViewRenderParameter(paramOrder []paramIdentity) viewRenderParameter {
	render := viewRenderParameter{
		opacity:         1.0,
		foregroundColor: param.foregroundColor.Get(),
		backgroundColor: param.backgroundColor.Get(),
		disable:         param.disable.Get(),
		hidden:          param.hidden.Get(),
		// padding:         param.padding,
		fontSize:        param.fontSize.Get(),
		fontWeight:      param.fontWeight.Get(),
		fontLineSpacing: param.fontLineSpacing.Get(),
		fontKerning:     param.fontKerning.Get(),
		fontItalic:      param.fontItalic.Get(),
	}

	for _, id := range paramOrder {
		switch id {
		case paramIDFrame:
			render.size = param.frameSize.Get()
		case paramIDOffset:
			b := param.offsets.Dequeue()
			render.offset = render.offset.Add(b.Get())
			param.offsets.Enqueue(b)
		case paramIDOpacity:
			b := param.opacity.Dequeue()
			render.opacity *= b.Get()
			param.opacity.Enqueue(b)
		case paramIDForegroundColor, paramIDBackgroundColor:
		}
	}

	return render
}

type viewRenderParameter struct {
	size            Size
	offset          Point
	stackOffset     Point
	centerOffset    Point
	opacity         float64
	foregroundColor color.Color
	backgroundColor color.Color
	disable         bool
	hidden          bool
	padding         []int
	fontSize        font.Size
	fontWeight      font.Weight
	fontLineSpacing float64
	fontKerning     int
	fontItalic      bool
}

func (v viewRenderParameter) Eq(v2 viewRenderParameter) bool {
	return v.size.Eq(v2.size) &&
		// v.offset.Eq(v2.offset) &&
		// v.stackOffset.Eq(v2.stackOffset) &&
		// v.centerOffset.Eq(v2.centerOffset) &&
		// v.opacity == v2.opacity &&
		v.foregroundColor == v2.foregroundColor &&
		v.backgroundColor == v2.backgroundColor &&
		// v.disable == v2.disable &&
		// v.hidden == v2.hidden &&
		// slices.Equal(v.padding, v2.padding) &&
		v.fontSize == v2.fontSize &&
		v.fontWeight == v2.fontWeight &&
		v.fontLineSpacing == v2.fontLineSpacing &&
		v.fontKerning == v2.fontKerning &&
		v.fontItalic == v2.fontItalic
}

func (v viewRenderParameter) Drawable() bool {
	return !v.hidden && v.size.W > 0 && v.size.H > 0 && v.opacity > 0.0
}
