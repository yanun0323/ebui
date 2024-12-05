package ebui

import (
	"image/color"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/gollection/v2"
	"github.com/yanun0323/pkg/sys"
)

type view struct {
	owner    SomeView
	img      gollection.Value[*ebiten.Image]
	noChange atomic.Bool

	render     viewRenderParameter
	param      viewParameter
	paramOrder []paramIdentity
}

func newView(id identity, owner SomeView) *view {
	return &view{
		owner:      owner,
		img:        gollection.NewSyncValue[*ebiten.Image](nil),
		param:      newViewParameter(id),
		paramOrder: make([]paramIdentity, 0, 10),
	}
}

var _ SomeView = (*view)(nil)

func (v *view) Body() SomeView {
	return v.owner
}

func (v *view) Frame(size Binding[Size]) SomeView {
	v.param.frameSize = size
	v.paramOrder = append(v.paramOrder, paramIDFrame)

	return v.owner
}

func (v *view) Offset(offset Binding[Point]) SomeView {
	v.param.offsets.Enqueue(offset)
	v.paramOrder = append(v.paramOrder, paramIDOffset)

	return v.owner
}

func (v *view) ForegroundColor(clr Binding[color.Color]) SomeView {
	v.param.foregroundColor = clr
	v.paramOrder = append(v.paramOrder, paramIDForegroundColor)

	return v.owner
}

func (v *view) BackgroundColor(clr Binding[color.Color]) SomeView {
	v.param.backgroundColor = clr
	v.paramOrder = append(v.paramOrder, paramIDBackgroundColor)

	return v.owner
}

func (v *view) Opacity(opacity Binding[float64]) SomeView {
	v.param.opacity.Enqueue(opacity)
	v.paramOrder = append(v.paramOrder, paramIDOpacity)

	return v.owner
}

// func (v *view) Disable(disable ...Binding[bool]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		switch len(disable) {
// 		case 0:
// 			v.param.disable = true
// 		case 1:
// 			v.param.disable = disable[0].Get()
// 		}
// 	})

// 	return v.owner
// }

// func (v *view) Padding(padding ...Binding[int]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		if len(padding) != 0 {
// 			ps := make([]int, 0, len(padding))
// 			for _, p := range padding {
// 				ps = append(ps, p.Get())
// 			}
// 			v.param.padding = ps
// 		}
// 	})

// 	return v.owner
// }

// func (v *view) FontSize(size Binding[font.Size]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		v.param.fontSize = size.Get()
// 	})

// 	return v.owner
// }

// func (v *view) FontWeight(weight Binding[font.Weight]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		v.param.fontWeight = weight.Get()
// 	})

// 	return v.owner
// }

// func (v *view) LineSpacing(spacing Binding[float64]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		v.param.fontLineSpacing = spacing.Get()
// 	})

// 	return v.owner
// }

// func (v *view) Kerning(kerning Binding[int]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		v.param.fontKerning = kerning.Get()
// 	})

// 	return v.owner
// }

// func (v *view) Italic(italic ...Binding[bool]) SomeView {
// 	v.modifiers = append(v.modifiers, func(v *view) {
// 		switch len(italic) {
// 		case 0:
// 			v.param.fontItalic = true
// 		case 1:
// 			v.param.fontItalic = italic[0].Get()
// 		}
// 	})

// 	return v.owner
// }

var _ internalView = (*view)(nil)

func (v *view) id() identity {
	return v.param.id
}

func (v *view) bounds() (min, current, max Size) {
	size := v.param.frameSize.Get()
	return size, size, size
}

func (v *view) update(container Size) {
	oldRender := v.render
	v.render = v.param.ViewRenderParameter(v.paramOrder)
	v.noChange.Store(v.render.Eq(oldRender))
}

func (v *view) getRenderImage() *ebiten.Image {
	if v.noChange.Swap(true) {
		return v.img.Load()
	}

	size := v.render.size
	img := ebiten.NewImage(size.W, size.H)
	img.Fill(v.render.foregroundColor)
	v.img.Store(img)

	return img
}

func (v *view) drawable() bool {
	return v.render.Drawable()
}

func (v *view) draw(screen *ebiten.Image) {
	if !v.owner.drawable() {
		return
	}

	img := v.owner.getRenderImage()
	if img != nil {
		v.drawResult(screen, img, nil)
	}
}

func (v *view) drawResult(screen, img *ebiten.Image, options *ebiten.DrawImageOptions) {
	if options == nil {
		options = &ebiten.DrawImageOptions{}
	}

	opacity := v.render.opacity
	opacity = sys.If(opacity > 1, 1, opacity)
	opacity = sys.If(opacity < 0, 0, opacity)
	options.ColorScale.ScaleAlpha(float32(opacity))

	offset := v.render.offset
	options.GeoM.Translate(float64(offset.X), float64(offset.Y))

	screen.DrawImage(img, options)
}
