package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
)

/* Check Interface Implementation */
var (
	_ SomeView = (*textView)(nil)
)

func Text(t string) *textView {
	v := &textView{
		t: t,
	}

	v.uiView = newUIView(typesText, v)
	return v
}

type textView struct {
	*uiView

	t string
}

func (*textView) textFace(opt *uiView, size ...font.Size) *text.GoTextFace {
	if len(size) != 0 {
		return &text.GoTextFace{
			Source: _defaultFontResource,
			Size:   size[0].F64(),
		}
	}

	face := &text.GoTextFace{
		Source: _defaultFontResource,
		Size:   opt.fontSize().F64(),
	}

	face.SetVariation(_fontTagWeight, opt.weight())
	face.SetVariation(_fontTagItalic, opt.italic())

	return face
}

// calculateSizeFromText returns the width, height, text and text face of the text view.
//
// It will truncate the text to fit the bounds.
func (v *textView) calculateSizeFromText(view *uiView) (int, int, string, *text.GoTextFace) {
	if len(v.t) == 0 {
		return 0, 0, "", nil
	}

	face := v.textFace(view)
	w, h := v.getTextBounds(v.t, face)
	if w <= view.size.h && h <= view.size.h {
		return w, h, v.t, face
	}

	tt := v.flexibleTruncateText(v.t, view, face)
	if len(tt) == 0 {
		t := "..."
		f := v.textFace(view, font.Body)
		w, h := text.Measure(t, f, view.lineSpacing())
		return int(w), int(h), t, f
	}

	w, h = v.getTextBounds(tt, face)
	return int(w), int(h), tt, face
}

func (v *textView) getTextBounds(t string, face *text.GoTextFace) (int, int) {
	w, h := text.Measure(t, face, v.lineSpacing())
	return int(w) + v.kern(), int(h)
}

func (v *textView) flexibleTruncateText(tt string, view *uiView, face *text.GoTextFace, ignoreDots ...bool) string {
	if len(tt) == 0 {
		return tt
	}

	width, height := view.Width(), view.Height()
	width -= view.kern() * len(tt)

	wf, hf := text.Measure(tt, face, view.lineSpacing())
	if int(hf) > height {
		return ""
	}

	if int(wf) <= width {
		return tt
	}

	l, m, r := 0, (len(tt)-2)/2, len(tt)-1
	for l < r && m != l && m != r {
		wf, _ := text.Measure(tt[:m], face, view.lineSpacing())
		if int(wf) >= width {
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

func (v *textView) draw(screen *ebiten.Image) {
	cache := v.Copy()
	cache.Draw(screen, func(screen *ebiten.Image) {
		cache.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, cache)
			if v != nil {
				v.draw(screen)
			}
		})

		w, h, tt, face := v.calculateSizeFromText(cache)
		if len(tt) != 0 && face != nil {
			dx := float64(cache.padding.left) + float64(cache.Width()-int(w))/2
			dy := float64(cache.padding.top) + float64(cache.Height()-int(h))/2

			for i, gl := range text.AppendGlyphs(nil, tt, face, &text.LayoutOptions{
				LineSpacing: cache.lineSpacing(),
			}) {
				if gl.Image == nil {
					continue
				}

				op := &ebiten.DrawImageOptions{}
				op.ColorScale.ScaleWithColor(cache.foregroundColor())
				// op.ColorScale.ScaleWithColor(color.Black)
				op.GeoM.Translate(float64(i*cache.kern()), 0)
				op.GeoM.Translate(dx, dy)
				op.GeoM.Translate(gl.X, gl.Y)
				screen.DrawImage(gl.Image, op)
			}
		}
	})
}
