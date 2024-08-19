package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/pkg/logs"
)

/* Check Interface Implement */
var (
	_ SomeView = (*textView[string])(nil)
	_ SomeView = (*textView[*string])(nil)
)

func Text[T string | *string](t T) *textView[T] {
	v := &textView[T]{
		t: t,
	}

	v.viewOption = newViewOption(v)
	return v
}

type textView[T string | *string] struct {
	viewOption

	t T
}

func (v *textView[T]) Body() SomeView {
	return v
}

func (v *textView[T]) text() string {
	switcher := func(a any) string {
		switch a := a.(type) {
		case string:
			return a
		case *string:
			if a == nil {
				return ""
			}

			return *a
		default:
			return ""
		}
	}

	return switcher(v.t)
}

func (v *textView[T]) textBounds(width, height int) (int, int, string, *text.GoTextFace) {
	face := &text.GoTextFace{
		Source: _mplusFaceSource,
		Size:   float64(v.viewOption.fontSize()),
	}

	tt := v.truncateText(v.text(), width, height, face)
	if len(tt) == 0 {
		return 0, 0, "", face
	}

	w, h := text.Measure(tt, face, 0)

	return int(w), int(h), tt, face
}

func (v *textView[T]) truncateText(tt string, width, height int, face *text.GoTextFace, ignoreDots ...bool) string {
	if len(tt) == 0 {
		return tt
	}

	wf, hf := text.Measure(tt, face, 0)
	if int(hf) > height {
		return ""
	}

	if int(wf) <= width {
		return tt
	}

	l, m, r := 0, (len(tt)-2)/2, len(tt)-1
	for l < r && m != l && m != r {
		wf, _ := text.Measure(tt[:m], face, 0)
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
		return v.truncateText(tt, width, height, face, true)
	}

	return tt[:m-truncateNum] + "..."
}

func (v *textView[T]) draw(screen *ebiten.Image, parent viewOption) image.Rectangle {
	if v != nil {
		current := parent.calculateViewOption(v.viewOption)
		logs.Debugf("TextView: %+v", current)
		result := image.Rect(current.x, current.y, current.x, current.y)
		current.Draw(screen, func(screen *ebiten.Image) {
			current.IterateViewModifiers(func(vm viewModifier) {
				v := vm(screen, &current)
				if v != nil {
					_ = v.draw(screen, current)
				}
			})

			ix, iy := current.XX(), current.YY()
			w, h, tt, face := v.textBounds(current.Width()-2, current.Height())
			if len(tt) != 0 {
				x := ix + (current.Width()-int(w))/2
				y := iy + (current.Height()-int(h))/2
				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(x), float64(y))

				text.Draw(screen, tt, face, op)
				result = current.DrawnArea()
			}
		})

		return result
	}

	return image.Rect(parent.x, parent.y, parent.x, parent.y)
}

func (v *textView[T]) bounds() (int, int) {
	w := v.w
	h := v.h
	return w, h
}
