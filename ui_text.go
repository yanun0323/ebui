package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/sys"
)

// import (
// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/text/v2"
// 	"github.com/yanun0323/ebui/font"
// )

// /* Check Interface Implementation */
var (
	_ SomeView = (*textView)(nil)
)

func Text(t string) *textView {
	v := &textView{
		t: t,
	}

	v.uiView = newView(typesText, v)
	return v
}

type textView struct {
	*uiView

	t string
}

func (*textView) textFace(v *uiView, size ...font.Size) *text.GoTextFace {
	if len(size) != 0 {
		return &text.GoTextFace{
			Source: _defaultFontResource,
			Size:   size[0].F64(),
		}
	}

	face := &text.GoTextFace{
		Source: _defaultFontResource,
		Size:   sys.If(v.fontSize == 0, font.Body, v.fontSize).F64(),
	}

	face.SetVariation(_fontTagWeight, v.fontWeight.F32())
	face.SetVariation(_fontTagItalic, sys.If[float32](v.italic, 1, 0))

	return face
}

// calculateSizeFromText returns the width, height, text and text face of the text view.
//
// It will truncate the text to fit the bounds.
func (v *textView) calculateSizeFromText(view *uiView, skipTruncate ...bool) (int, int, string, *text.GoTextFace) {
	if len(v.t) == 0 {
		return 0, 0, "", nil
	}

	face := v.textFace(view)
	w, h := v.getTextBounds(v.t, face)
	drawSize := view.getDrawSize(view.cachedSize)
	if w <= drawSize.w && h <= drawSize.h {
		return w, h, v.t, face
	}

	tt := v.t
	if len(skipTruncate) == 0 || !skipTruncate[0] {
		tt = v.flexibleTruncateText(v.t, view, face)
	}

	if len(tt) == 0 {
		t := "..."
		f := v.textFace(view, font.Body)
		w, h := text.Measure(t, f, view.lineSpacing)
		return int(w), int(h), t, f
	}

	w, h = v.getTextBounds(tt, face)
	return int(w), int(h), tt, face
}

func (v *textView) getTextBounds(t string, face *text.GoTextFace) (int, int) {
	w, h := text.Measure(t, face, v.lineSpacing)
	return int(w) + v.kerning, int(h)
}

func (v *textView) flexibleTruncateText(tt string, view *uiView, face *text.GoTextFace, ignoreDots ...bool) string {
	if len(tt) == 0 {
		return tt
	}

	drawSize := view.getDrawSize(view.cachedSize)
	drawSize.w -= view.kerning * len(tt)

	wf, _ := text.Measure(tt, face, view.lineSpacing)
	if int(wf) <= drawSize.w {
		return tt
	}

	l, m, r := 0, (len(tt)-2)/2, len(tt)-1
	for l < r && m != l && m != r {
		wf, _ := text.Measure(tt[:m], face, view.lineSpacing)
		if int(wf) >= drawSize.w {
			r = m
			m = (l + m) / 2
		} else {
			l = m
			m = (r + m) / 2
		}
	}

	if len(ignoreDots) != 0 && ignoreDots[0] {
		return tt[:m]
	}

	truncateNum := 2
	if m <= truncateNum {
		return v.flexibleTruncateText(tt, view, face, true)
	}

	return tt[:m-truncateNum] + "..."
}

func (v *textView) getSize() size {
	if v.isCached {
		return v.cachedSize
	}

	s := v.uiView.getSize()
	if s.IsZero() {
		w, h, _, _ := v.calculateSizeFromText(v.uiView, true)
		s = size{w, h}
	}

	v.isCached = true
	v.cachedSize = s
	return v.cachedSize
}

func (v *textView) draw(screen *ebiten.Image) {
	w, h, tt, face := v.calculateSizeFromText(v.uiView)
	if len(tt) != 0 && face != nil {
		size := v.getDrawSize(v.uiView.cachedSize)
		dx := float64(v.uiView.start.x) + float64(size.w-int(w))/2
		dy := float64(v.uiView.start.y) + float64(size.h-int(h))/2

		for i, gl := range text.AppendGlyphs(nil, tt, face, &text.LayoutOptions{LineSpacing: v.uiView.lineSpacing}) {
			if gl.Image == nil {
				continue
			}
			opt := &ebiten.DrawImageOptions{}
			opt.ColorScale.ScaleWithColor(sys.If(v.background == nil, color.Color(color.White), v.background))

			opt.GeoM.Translate(float64(i*v.uiView.kerning), 0)
			opt.GeoM.Translate(dx, dy)
			opt.GeoM.Translate(gl.X, gl.Y)
			v.handlePreference(opt)
			screen.DrawImage(gl.Image, opt)
		}
	}
	v.drawModifiers(screen)
}
