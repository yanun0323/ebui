package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/sys"
)

type viewModifier func(view) view

type viewParameter struct {
	frameSize       Size
	offset          Point
	foregroundColor color.Color
	backgroundColor color.Color
	opacity         float64
	disable         bool
	padding         []int
	fontSize        font.Size
	fontWeight      font.Weight
	fontLineSpacing float64
	fontKerning     int
	fontItalic      bool
}

func newViewParameter() viewParameter {
	return viewParameter{
		frameSize:       Size{W: 0, H: 0},
		offset:          Point{X: 0, Y: 0},
		foregroundColor: color.White,
		backgroundColor: color.Black,
		opacity:         1,
		disable:         false,
		padding:         []int{15, 15, 15, 15},
		fontSize:        font.Body,
		fontWeight:      font.Normal,
		fontLineSpacing: 0,
		fontKerning:     0,
		fontItalic:      false,
	}
}

type view struct {
	owner     SomeView
	param     viewParameter
	modifiers []viewModifier
}

func newView(owner SomeView) view {
	return view{
		owner:     owner,
		param:     newViewParameter(),
		modifiers: make([]viewModifier, 0, 10),
	}
}

var _ SomeView = (*view)(nil)

func (v *view) Body() SomeView {
	return v.owner
}

func (v *view) Frame(size Binding[Size]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.frameSize = size.Get()
		return v
	})

	return v.owner
}

func (v *view) Offset(offset Binding[Point]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.offset = offset.Get()
		return v
	})

	return v.owner
}

func (v *view) ForegroundColor(clr Binding[color.Color]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.foregroundColor = clr.Get()
		return v
	})

	return v.owner
}

func (v *view) BackgroundColor(clr Binding[color.Color]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.backgroundColor = clr.Get()
		return v
	})

	return v.owner
}

func (v *view) Opacity(opacity Binding[float64]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.opacity = opacity.Get()
		return v
	})

	return v.owner
}

func (v *view) Disable(disable ...Binding[bool]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		switch len(disable) {
		case 0:
			v.param.disable = true
		case 1:
			v.param.disable = disable[0].Get()
		}
		return v
	})

	return v.owner
}

func (v *view) Padding(padding ...Binding[int]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		if len(padding) != 0 {
			ps := make([]int, 0, len(padding))
			for _, p := range padding {
				ps = append(ps, p.Get())
			}
			v.param.padding = ps
		}
		return v
	})

	return v.owner
}

func (v *view) FontSize(size Binding[font.Size]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.fontSize = size.Get()
		return v
	})

	return v.owner
}

func (v *view) FontWeight(weight Binding[font.Weight]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.fontWeight = weight.Get()
		return v
	})

	return v.owner
}

func (v *view) LineSpacing(spacing Binding[float64]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.fontLineSpacing = spacing.Get()
		return v
	})

	return v.owner
}

func (v *view) Kerning(kerning Binding[int]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		v.param.fontKerning = kerning.Get()
		return v
	})

	return v.owner
}

func (v *view) Italic(italic ...Binding[bool]) SomeView {
	v.modifiers = append(v.modifiers, func(v view) view {
		switch len(italic) {
		case 0:
			v.param.fontItalic = true
		case 1:
			v.param.fontItalic = italic[0].Get()
		}
		return v
	})

	return v.owner
}

var _ internalView = (*view)(nil)

func (v *view) modify() {
	for _, fn := range v.modifiers {
		*v = fn(*v)
	}
}

func (v view) drawOption() *ebiten.DrawImageOptions {
	opacity := v.param.opacity
	opacity = sys.If(opacity > 1, 1, opacity)
	opacity = sys.If(opacity < 0, 0, opacity)

	opt := &ebiten.DrawImageOptions{}
	opt.ColorScale.ScaleAlpha(float32(opacity))

	return opt
}

func (v view) bounds() (min, current, max Size) {
	v.modify()

	return v.param.frameSize, v.param.frameSize, v.param.frameSize
}

func (v view) draw(screen *ebiten.Image) {
	v.modify()

	frameSize := v.param.frameSize
	img := ebiten.NewImage(frameSize.W, frameSize.H)
	img.Fill(v.param.foregroundColor)

	screen.DrawImage(img, v.drawOption())
}
