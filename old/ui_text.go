package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/pkg/sys"
)

type texts struct {
	*view

	localeKey Binding[string]
}

func Text(localeKey Binding[string]) SomeView {
	v := &texts{
		localeKey: localeKey,
	}
	v.view = newView(idText, v)
	return v
}

func (v *texts) getRenderImage() *ebiten.Image {
	if v.noChange.Swap(true) {
		return v.img.Load()
	}

	size := v.render.size

	t := v.localeKey.Get()
	face := &text.GoTextFace{
		Source: _defaultFontResource,
		Size:   v.param.fontSize.Get().F64(),
	}

	face.SetVariation(_fontTagWeight, v.param.fontWeight.Get().F32())
	face.SetVariation(_fontTagItalic, sys.If[float32](v.param.fontItalic.Get(), 1, 0))
	wf, hf := text.Measure(t, face, v.param.fontLineSpacing.Get())

	chars := text.AppendGlyphs(nil, t, face, &text.LayoutOptions{LineSpacing: v.param.fontLineSpacing.Get()})
	w := int(wf) + (v.param.fontKerning.Get() * len(chars))
	h := int(hf)

	w = max(w, size.W)
	h = max(h, size.H)

	img := ebiten.NewImage(w, h)
	// dx := /*  float64(v.uiView.start.x) */ +float64(size.W-int(w)) / 2
	// dy := /*  float64(v.uiView.start.y) */ +float64(size.H-int(h)) / 2
	for i, gl := range chars {
		if gl.Image == nil {
			continue
		}
		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(v.param.foregroundColor.Get())

		opt.GeoM.Translate(float64(i*v.param.fontKerning.Get()), 0)
		// opt.GeoM.Translate(dx, dy)
		opt.GeoM.Translate(gl.X, gl.Y)
		img.DrawImage(gl.Image, opt)
	}

	v.img.Store(img)
	return img
}
