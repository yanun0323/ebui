package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

type viewOption struct {
	holder                                               SomeView
	x, y                                                 int
	xx, yy                                               int /* absolute x, y */
	w, h                                                 int
	paddingTop, paddingBottom, paddingLeft, paddingRight int
	flexibleW, flexibleH                                 int
	fColor                                               color.Color
	bColor                                               color.Color
	viewModifiers                                        []viewModifier
	subviews                                             []viewModifier
	fontSizes                                            font.Size
	opacities                                            float32
	isPressing                                           bool
	cornerRadius                                         int
}

func newViewOption(holder SomeView, svs ...View) viewOption {
	return viewOption{
		holder:    holder,
		w:         -1,
		h:         -1,
		opacities: 1.0,
		subviews:  newViewModifiers(svs...),
	}
}

func (opt *viewOption) Draw(screen *ebiten.Image, painter func(screen *ebiten.Image)) {
	if opt.w <= 0 || opt.h <= 0 {
		return
	}

	img := ebiten.NewImage(opt.w, opt.h)
	painter(img)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(opt.opacity())
	op.GeoM.Translate(float64(opt.xx), float64(opt.yy))

	makeImageRounded(img, opt.cornerRadius)
	screen.DrawImage(img, op)
}

func (opt viewOption) CreateChild(holder SomeView, x, y, xx, yy, flexibleW, flexibleH int, size ...image.Rectangle) viewOption {
	w, h := opt.w, opt.h
	if len(size) != 0 {
		w, h = size[0].Dx(), size[0].Dy()
	}

	return viewOption{
		holder:        holder,
		x:             x,
		y:             y,
		xx:            xx,
		yy:            yy,
		w:             w,
		h:             h,
		paddingTop:    opt.paddingTop,
		paddingBottom: opt.paddingBottom,
		paddingLeft:   opt.paddingLeft,
		paddingRight:  opt.paddingRight,
		flexibleW:     flexibleW,
		flexibleH:     flexibleH,
		fColor:        opt.fColor,
		bColor:        opt.bColor,
		fontSizes:     opt.fontSizes,
		opacities:     opt.opacities,
		isPressing:    opt.isPressing,
	}
}

func (opt viewOption) X() int {
	if opt.paddingLeft >= opt.w {
		return opt.w
	}

	return opt.x + opt.paddingLeft
}

func (opt viewOption) Y() int {
	if opt.paddingTop >= opt.h {
		return opt.h
	}

	return opt.y + opt.paddingTop
}

func (opt viewOption) XX() int {
	if opt.paddingLeft >= opt.w {
		return 0
	}

	return opt.paddingLeft
}

func (opt viewOption) YY() int {
	if opt.paddingTop >= opt.h {
		return 0
	}

	return opt.paddingTop
}

func (opt viewOption) Width() int {
	if opt.paddingLeft >= opt.w ||
		opt.paddingRight >= opt.w {
		return 0
	}

	result := opt.w - opt.paddingLeft - opt.paddingRight
	if result <= 0 {
		result = 0
	}

	return result
}

func (opt viewOption) Height() int {
	if opt.paddingTop >= opt.h ||
		opt.paddingBottom >= opt.h {
		return 0
	}

	result := opt.h - opt.paddingTop - opt.paddingBottom
	if result <= 0 {
		result = 0
	}

	return result
}

func (opt viewOption) foregroundColor() color.Color {
	if opt.fColor == nil {
		return _defaultForegroundColor
	}

	return opt.fColor
}

func (opt viewOption) fontSize() font.Size {
	if opt.fontSizes <= 0 {
		return font.Body
	}

	return opt.fontSizes
}

func (opt viewOption) opacity() float32 {
	if opt.opacities <= 0 {
		return 0
	}

	if opt.opacities > 1 {
		return 1
	}

	if opt.isPressing {
		return opt.opacities * 0.75
	}

	return opt.opacities
}

func (opt viewOption) DrawnArea() image.Rectangle {
	return image.Rect(opt.x, opt.y, opt.x+opt.w, opt.y+opt.h)
}

func (opt viewOption) Contain(x, y int) bool {
	px, py := opt.X(), opt.Y()
	pw, ph := opt.Width(), opt.Height()
	return px <= x && x <= px+pw && py <= y && y <= py+ph
}

func (opt *viewOption) IterateViewModifiers(modifier func(viewModifier), subviews ...func(viewModifier)) {
	subviewsHandler := modifier
	if len(subviews) != 0 && subviews[0] != nil {
		subviewsHandler = subviews[0]
	}

	for _, vm := range opt.viewModifiers {
		modifier(vm)
	}

	for _, sv := range opt.subviews {
		subviewsHandler(sv)
	}
}

func (opt *viewOption) calculateViewOption(current viewOption) viewOption {
	opt.updateSizeTo(opt, &current)
	opt.updateColorTo(opt, &current)
	opt.updateOthersTo(opt, &current)
	current.calculateFlexibleSizeTo()
	return current
}

func (viewOption) updateSizeTo(parent, current *viewOption) {
	current.x = parent.X()
	current.y = parent.Y()
	current.xx = parent.xx
	current.yy = parent.yy

	if current.w == -1 {
		current.w = parent.flexibleW
	}

	if current.w > parent.Width() {
		current.w = parent.Width()
	}

	if current.h == -1 {
		current.h = parent.flexibleH
	}

	if current.h > parent.Height() {
		current.h = parent.Height()
	}
}

func (viewOption) updateColorTo(parent, current *viewOption) {
	if current.fColor == nil {
		current.fColor = parent.fColor
	}
}

func (opt *viewOption) calculateFlexibleSizeTo() {
	wCount, hCount := 0, 0
	w, h := opt.Width(), opt.Height()
	for _, vm := range opt.subviews {
		v := vm(nil, nil)
		if v != nil {
			vw, vh := v.bounds()
			if vw != -1 {
				w -= vw
			} else {
				wCount++
			}

			if vh != -1 {
				h -= vh
			} else {
				hCount++
			}
		}
	}

	if wCount == 0 {
		wCount = 1
	}

	if hCount == 0 {
		hCount = 1
	}

	opt.flexibleW = w / wCount
	opt.flexibleH = h / hCount
}

func (viewOption) updateOthersTo(parent, current *viewOption) {
	current.fontSizes = parent.fontSizes
	current.opacities = parent.opacities
	current.isPressing = parent.isPressing
}

/*
	Interface Implement
*/

/* Check Interface Implement */
var _ someViewOption = (*viewOption)(nil)

func (opt *viewOption) option() *viewOption {
	return opt
}

func (opt *viewOption) Frame(w, h int) SomeView {
	opt.w = w
	opt.h = h
	return opt.holder
}

func (opt *viewOption) Body() SomeView {
	return opt.holder
}

func (opt *viewOption) ForegroundColor(clr color.Color) SomeView {
	opt.fColor = clr
	return opt.holder
}

func (opt *viewOption) BackgroundColor(clr color.Color) SomeView {
	opt.viewModifiers = append(opt.viewModifiers, backgroundColorViewModifier(clr))
	return opt.holder
}

func (opt *viewOption) Padding(top, right, bottom, left int) SomeView {
	opt.viewModifiers = append(opt.viewModifiers, paddingViewModifier(top, right, bottom, left))
	return opt.holder
}

func (opt *viewOption) FontSize(size font.Size) SomeView {
	opt.fontSizes = size
	return opt.holder
}

func (opt *viewOption) CornerRadius(radius ...int) SomeView {
	if len(radius) != 0 {
		opt.cornerRadius = radius[0]
	} else {
		opt.cornerRadius = 7
	}

	return opt.holder
}
