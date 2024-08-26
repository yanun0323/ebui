package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

type View interface {
	Body() SomeView
}

//go:generate domaingen -destination=view.option.go -package=ebui -name=viewOption
type SomeView interface {
	View

	draw(screen *ebiten.Image)
	view() *uiView
	initBounds() (int, int)

	Frame(w, h int) SomeView
	ForegroundColor(clr color.Color) SomeView
	BackgroundColor(clr color.Color) SomeView
	Padding(top, right, bottom, left int) SomeView
	FontSize(size font.Size) SomeView
	FontWeight(weight font.Weight) SomeView
	CornerRadius(radius ...int) SomeView
	LineSpacing(spacing float64) SomeView
	Kern(kern int) SomeView
	Italic() SomeView
}

type uiView struct {
	holder SomeView
	types  types

	size frame /* auto calculated */

	/* params */
	initSize         frame /* initial bounds */
	minSize, maxSize frame
	pos              point /* absolute coords, relative to (0,0) */
	xx, yy           int   /* relative coords, relative to parent (x, y) */
	padding          bounds
	viewModifiers    []viewModifier
	subviews         []SomeView
	viewOpacity      float32

	/* no inherit */

	cornerRadius    int         /* no inherit */
	fontBold        bool        /* no inherit */
	fontWeight      font.Weight /* no inherit */
	fontItalic      bool        /* no inherit */
	fontKern        int         /* no inherit */
	fontLineSpacing float64     /* no inherit */

	/* inherit if nil/zero */

	fColor     color.Color /* inherit if nil/zero */
	fontSizes  font.Size   /* inherit if nil/zero */
	isPressing bool        /* inherit if nil/zero */
}

func newUIView(id types, holder SomeView, svs ...View) *uiView {
	subviews := make([]SomeView, 0, len(svs))
	for _, sv := range svs {
		subviews = append(subviews, sv.Body())
	}

	return &uiView{
		holder:      holder,
		types:       id,
		size:        frame{-1, -1},
		initSize:    frame{-1, -1},
		minSize:     frame{-1, -1},
		maxSize:     frame{-1, -1},
		viewOpacity: 1.0,
		subviews:    subviews,
	}
}

func (v *uiView) Draw(screen *ebiten.Image, painter func(screen *ebiten.Image)) {
	if v.size.w <= 0 || v.size.h <= 0 {
		return
	}

	img := ebiten.NewImage(v.size.w, v.size.h)
	painter(img)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(v.opacity())
	switch v.types {

	default:
		op.GeoM.Translate(float64(v.xx), float64(v.yy))
	}

	makeImageRounded(img, v.cornerRadius)
	screen.DrawImage(img, op)
}

func (v uiView) Copy() *uiView {
	return &v
}

func (v uiView) DrawnArea() image.Rectangle {
	return image.Rect(v.pos.x, v.pos.y, v.pos.x+v.size.w, v.pos.y+v.size.h)
}

func (v uiView) Contain(x, y int) bool {
	px, py := v.X(), v.Y()
	pw, ph := v.Width(), v.Height()
	return px <= x && x <= px+pw && py <= y && y <= py+ph
}

func (v *uiView) IterateViewModifiers(modifier func(viewModifier), subviews ...func(viewModifier)) {
	subviewsHandler := modifier
	if len(subviews) != 0 && subviews[0] != nil {
		subviewsHandler = subviews[0]
	}

	for _, vm := range v.viewModifiers {
		modifier(vm)
	}

	for _, sv := range v.subviews {
		subviewsHandler(func(*ebiten.Image, *uiView) SomeView {
			return sv
		})
	}
}

/*
	########     ###    ########     ###    ##     ## ######## ######## ######## ########   ######
	##     ##   ## ##   ##     ##   ## ##   ###   ### ##          ##    ##       ##     ## ##    ##
	##     ##  ##   ##  ##     ##  ##   ##  #### #### ##          ##    ##       ##     ## ##
	########  ##     ## ########  ##     ## ## ### ## ######      ##    ######   ########   ######
	##        ######### ##   ##   ######### ##     ## ##          ##    ##       ##   ##         ##
	##        ##     ## ##    ##  ##     ## ##     ## ##          ##    ##       ##    ##  ##    ##
	##        ##     ## ##     ## ##     ## ##     ## ########    ##    ######## ##     ##  ######
*/

func (v uiView) X() int {
	if v.padding.left >= v.size.w {
		return v.size.w
	}

	return v.pos.x + v.padding.left
}

func (v uiView) Y() int {
	if v.padding.top >= v.size.h {
		return v.size.h
	}

	return v.pos.y + v.padding.top
}

func (v uiView) Width() int {
	result := v.size.w - v.padding.left - v.padding.right
	if result <= 0 {
		result = 0
	}

	return result
}

func (v uiView) Height() int {
	result := v.size.h - v.padding.top - v.padding.bottom
	if result <= 0 {
		result = 0
	}

	return result
}

func (v uiView) foregroundColor() color.Color {
	if v.fColor == nil {
		return _defaultForegroundColor
	}

	return v.fColor
}

func (v uiView) fontSize() font.Size {
	if v.fontSizes <= 0 {
		return font.Body
	}

	return v.fontSizes
}

func (v uiView) weight() float32 {
	if v.fontBold {
		return font.Bold.F32()
	}

	if v.fontWeight <= 0 {
		return font.Normal.F32()
	}

	return float32(v.fontWeight)
}

func (v uiView) lineSpacing() float64 {
	return v.fontLineSpacing
}

func (v uiView) kern() int {
	return v.fontKern
}

func (v uiView) italic() float32 {
	if v.fontItalic {
		return 1
	}

	return 0
}

func (v uiView) opacity() float32 {
	if v.viewOpacity <= 0 {
		return 0
	}

	if v.viewOpacity > 1 {
		return 1
	}

	if v.isPressing {
		return v.viewOpacity * 0.75
	}

	return v.viewOpacity
}

/*
	#### ##     ## ########  ##       ######## ##     ## ######## ##    ## ########
	 ##  ###   ### ##     ## ##       ##       ###   ### ##       ###   ##    ##
	 ##  #### #### ##     ## ##       ##       #### #### ##       ####  ##    ##
	 ##  ## ### ## ########  ##       ######   ## ### ## ######   ## ## ##    ##
	 ##  ##     ## ##        ##       ##       ##     ## ##       ##  ####    ##
	 ##  ##     ## ##        ##       ##       ##     ## ##       ##   ###    ##
	#### ##     ## ##        ######## ######## ##     ## ######## ##    ##    ##
*/

/* Check Interface Implementation */
var _ SomeView = (*uiView)(nil)

func (v *uiView) Body() SomeView {
	return v.holder
}

func (uiView) draw(*ebiten.Image) {}

func (v *uiView) initBounds() (int, int) {
	cache := v.Copy()
	cache.IterateViewModifiers(func(vm viewModifier) {
		_ = vm(nil, cache)
	})

	return rpNeq(cache.initSize.w, -1, cache.size.w), rpNeq(cache.initSize.h, -1, cache.size.h)
}

func (v *uiView) view() *uiView {
	return v
}

func (v *uiView) Frame(w, h int) SomeView {
	v.initSize = frame{w, h}
	v.viewModifiers = append(v.viewModifiers, frameViewModifier(w, h))
	return v.holder
}

func (v *uiView) ForegroundColor(clr color.Color) SomeView {
	v.fColor = clr
	return v.holder
}

func (v *uiView) BackgroundColor(clr color.Color) SomeView {
	v.viewModifiers = append(v.viewModifiers, backgroundColorViewModifier(clr))
	return v.holder
}

func (v *uiView) Padding(top, right, bottom, left int) SomeView {
	v.viewModifiers = append(v.viewModifiers, paddingViewModifier(top, right, bottom, left))
	return v.holder
}

func (v *uiView) FontSize(size font.Size) SomeView {
	v.fontSizes = size
	return v.holder
}

func (v *uiView) FontWeight(weight font.Weight) SomeView {
	v.fontWeight = weight
	return v.holder
}

func (v *uiView) LineSpacing(spacing float64) SomeView {
	v.fontLineSpacing = spacing
	return v.holder
}

func (v *uiView) Kern(kern int) SomeView {
	v.fontKern = kern
	return v.holder
}

func (v *uiView) Italic() SomeView {
	v.fontItalic = true
	return v.holder
}

func (v *uiView) Bold() SomeView {
	v.fontBold = true
	return v.holder
}

func (v *uiView) CornerRadius(radius ...int) SomeView {
	if len(radius) != 0 {
		v.cornerRadius = radius[0]
	} else {
		v.cornerRadius = 7
	}

	return v.holder
}
