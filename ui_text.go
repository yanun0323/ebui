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

func (v *texts) draw(screen *ebiten.Image) {
	v.updateRenderCache()

	t := v.localeKey.Get()
	face := &text.GoTextFace{
		Source: _defaultFontResource,
		Size:   v.param.fontSize.F64(),
	}

	face.SetVariation(_fontTagWeight, v.param.fontWeight.F32())
	face.SetVariation(_fontTagItalic, sys.If[float32](v.param.fontItalic, 1, 0))
	wf, hf := text.Measure(t, face, v.param.fontLineSpacing)

	w := int(wf) + v.param.fontKerning
	h := int(hf)

	size := v.param.frameSize

	dx := /*  float64(v.uiView.start.x) */ +float64(size.W-int(w)) / 2
	dy := /*  float64(v.uiView.start.y) */ +float64(size.H-int(h)) / 2

	for i, gl := range text.AppendGlyphs(nil, t, face, &text.LayoutOptions{LineSpacing: v.param.fontLineSpacing}) {
		if gl.Image == nil {
			continue
		}
		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(v.param.foregroundColor)

		opt.GeoM.Translate(float64(i*v.param.fontKerning), 0)
		opt.GeoM.Translate(dx, dy)
		opt.GeoM.Translate(gl.X, gl.Y)
		v.drawResult(screen, gl.Image, opt)
	}
}
