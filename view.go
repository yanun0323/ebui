package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

type view struct {
	holder SomeView
	types  types

	w, h int /* auto calculated */

	/* params */
	initW, initH                                         int /* initial bounds */
	x, y                                                 int /* absolute coords, relative to (0,0) */
	xx, yy                                               int /* relative coords, relative to parent (x, y) */
	minW, minH                                           int
	maxW, maxH                                           int
	paddingTop, paddingBottom, paddingLeft, paddingRight int
	viewModifiers                                        []viewModifier
	subviews                                             []SomeView
	viewOpacity                                          float32

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

func newView(id types, holder SomeView, svs ...View) *view {
	subviews := make([]SomeView, 0, len(svs))
	for _, sv := range svs {
		subviews = append(subviews, sv.Body())
	}

	return &view{
		holder:      holder,
		types:       id,
		initW:       -1,
		initH:       -1,
		w:           -1,
		h:           -1,
		viewOpacity: 1.0,
		subviews:    subviews,
	}
}

func (v *view) setInitSize(w, h int) {
	v.initW, v.initH = w, h
	v.w, v.h = w, h
}

func (v *view) Draw(screen *ebiten.Image, painter func(screen *ebiten.Image)) {
	if v.w <= 0 || v.h <= 0 {
		return
	}

	img := ebiten.NewImage(v.w, v.h)
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

func (v view) Copy() *view {
	return &v
}

func (v view) DrawnArea() image.Rectangle {
	return image.Rect(v.x, v.y, v.x+v.w, v.y+v.h)
}

func (v view) Contain(x, y int) bool {
	px, py := v.X(), v.Y()
	pw, ph := v.Width(), v.Height()
	return px <= x && x <= px+pw && py <= y && y <= py+ph
}

func (v *view) IterateViewModifiers(modifier func(viewModifier), subviews ...func(viewModifier)) {
	subviewsHandler := modifier
	if len(subviews) != 0 && subviews[0] != nil {
		subviewsHandler = subviews[0]
	}

	for _, vm := range v.viewModifiers {
		modifier(vm)
	}

	for _, sv := range v.subviews {
		subviewsHandler(func(*ebiten.Image, *view) SomeView {
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

func (v view) X() int {
	if v.paddingLeft >= v.w {
		return v.w
	}

	return v.x + v.paddingLeft
}

func (v view) Y() int {
	if v.paddingTop >= v.h {
		return v.h
	}

	return v.y + v.paddingTop
}

func (v view) XX() int {
	return v.xx + v.paddingLeft
}

func (v view) YY() int {
	return v.xx + v.paddingLeft
}

func (v view) Width() int {
	result := v.w - v.paddingLeft - v.paddingRight
	if result <= 0 {
		result = 0
	}

	return result
}

func (v view) Height() int {
	result := v.h - v.paddingTop - v.paddingBottom
	if result <= 0 {
		result = 0
	}

	return result
}

func (v view) foregroundColor() color.Color {
	if v.fColor == nil {
		return _defaultForegroundColor
	}

	return v.fColor
}

func (v view) fontSize() font.Size {
	if v.fontSizes <= 0 {
		return font.Body
	}

	return v.fontSizes
}

func (v view) weight() float32 {
	if v.fontBold {
		return font.Bold.F32()
	}

	if v.fontWeight <= 0 {
		return font.Normal.F32()
	}

	return float32(v.fontWeight)
}

func (v view) lineSpacing() float64 {
	return v.fontLineSpacing
}

func (v view) kern() int {
	return v.fontKern
}

func (v view) italic() float32 {
	if v.fontItalic {
		return 1
	}

	return 0
}

func (v view) opacity() float32 {
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
var _ SomeView = (*view)(nil)

func (v *view) Body() SomeView {
	return v.holder
}

func (view) draw(*ebiten.Image) {}

func (v *view) initBounds() (int, int) {
	return v.initW, v.initH
}

func (v *view) params() *view {
	return v
}

func (v *view) Frame(w, h int) SomeView {
	v.initW = w
	v.initH = h
	return v.holder
}

func (v *view) ForegroundColor(clr color.Color) SomeView {
	v.fColor = clr
	return v.holder
}

func (v *view) BackgroundColor(clr color.Color) SomeView {
	v.viewModifiers = append(v.viewModifiers, backgroundColorViewModifier(clr))
	return v.holder
}

func (v *view) Padding(top, right, bottom, left int) SomeView {
	v.viewModifiers = append(v.viewModifiers, paddingViewModifier(top, right, bottom, left))
	return v.holder
}

func (v *view) FontSize(size font.Size) SomeView {
	v.fontSizes = size
	return v.holder
}

func (v *view) FontWeight(weight font.Weight) SomeView {
	v.fontWeight = weight
	return v.holder
}

func (v *view) LineSpacing(spacing float64) SomeView {
	v.fontLineSpacing = spacing
	return v.holder
}

func (v *view) Kern(kern int) SomeView {
	v.fontKern = kern
	return v.holder
}

func (v *view) Italic() SomeView {
	v.fontItalic = true
	return v.holder
}

func (v *view) Bold() SomeView {
	v.fontBold = true
	return v.holder
}

func (v *view) CornerRadius(radius ...int) SomeView {
	if len(radius) != 0 {
		v.cornerRadius = radius[0]
	} else {
		v.cornerRadius = 7
	}

	return v.holder
}
