package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implementation */
var (
	_ SomeView = (*textView)(nil)
)

func Text(t string) *textView {
	v := &textView{
		t: t,
	}

	v.view = newView(typeText, v)
	return v
}

type textView struct {
	*view

	t string
}

func (v *textView) initBounds() (int, int) {
	w, h := v.w, v.h

	if w == -1 {

	}

	if h == -1 {

	}

	return w, h
}

func (*textView) textFace(opt *view, size ...font.Size) *text.GoTextFace {
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

func (v *textView) calculateSizeFromText() {
	if len(v.t) == 0 {
		return
	}

}

func (v *textView) getTextBounds(t string, face *text.GoTextFace) (int, int) {
	w, h := text.Measure(t, face, v.lineSpacing())
	return int(w) + v.kern(), int(h)
}

// flexibleTextBounds returns the width, height, text and text face of the text view.
//
// if the frame did'nt set, it will truncate the text to fit the flexible bounds.
func (v *textView) flexibleTextBounds(opt *view) (int, int, string, *text.GoTextFace) {
	if len(v.t) == 0 {
		return 0, 0, "", nil
	}

	face := v.textFace(opt)
	width, height := opt.initW, opt.initH
	width -= opt.kern() * len(v.t)

	if v.w == width && v.h == height {
		return width, height, v.t, face
	}

	tt := v.flexibleTruncateText(v.t, opt, face)
	if len(tt) == 0 {
		t := "..."
		f := v.textFace(opt, font.Body)
		w, h := text.Measure(t, f, opt.lineSpacing())
		return int(w), int(h), t, f
	}

	w, h := text.Measure(tt, face, opt.lineSpacing())

	return int(w), int(h), tt, face
}

func (v *textView) flexibleTruncateText(tt string, opt *view, face *text.GoTextFace, ignoreDots ...bool) string {
	if len(tt) == 0 {
		return tt
	}

	width, height := opt.Width(), opt.Height()
	width -= opt.kern() * len(tt)

	wf, hf := text.Measure(tt, face, opt.lineSpacing())
	if int(hf) > height {
		return ""
	}

	if int(wf) <= width {
		return tt
	}

	l, m, r := 0, (len(tt)-2)/2, len(tt)-1
	for l < r && m != l && m != r {
		wf, _ := text.Measure(tt[:m], face, opt.lineSpacing())
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
		return v.flexibleTruncateText(tt, opt, face, true)
	}

	return tt[:m-truncateNum] + "..."
}

func (v *textView) draw(screen *ebiten.Image) {
	cache := v.Copy()
	cache.Draw(screen, func(screen *ebiten.Image) {
		// FIXME: Fix me
		cache.IterateViewModifiers(func(vm viewModifier) {
			v := vm(screen, cache)
			if v != nil {
				v.draw(screen)
			}
		})

		w, h, tt, face := v.flexibleTextBounds(cache)
		if len(tt) != 0 && face != nil {
			dx := float64(cache.Width()-int(w)) / 2
			dy := float64(cache.Height()-int(h)) / 2

			logs.Warnf("cache.Width(): %d, cache.Height(): %d\ncache.iniW: %d, cache.iniH: %d\nw: %d, h: %d, dx: %.1f, dy: %.1f", cache.Width(), cache.Height(), cache.initW, cache.initH, w, h, dx, dy)

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
