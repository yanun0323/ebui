package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/logs"
)

type View interface {
	Body() SomeView
}

type SomeView interface {
	View

	Frame(w, h int) SomeView
	ForegroundColor(clr color.Color) SomeView
	BackgroundColor(clr color.Color) SomeView

	// Padding makes view have padding.
	//
	// # parameter count:
	// 	 - [1] all
	// 	 - [2] vertical, horizontal
	// 	 - [4] top, right, bottom, left
	Padding(padding ...int) SomeView

	Offset(x, y int) SomeView

	FontSize(size font.Size) SomeView
	FontWeight(weight font.Weight) SomeView

	CornerRadius(radius ...int) SomeView
	LineSpacing(spacing float64) SomeView
	Kerning(kerning int) SomeView
	Italic() SomeView
}

type uiViewDelegator interface {
	UIView() *uiView
}

type uiViewModifier struct {
	uiViewLayout
	uiViewParameter
}

func newViewModifier() uiViewModifier {
	return uiViewModifier{
		uiViewLayout: _zeroUIViewLayout,
	}
}

type uiView struct {
	uiViewLayout
	uiViewParameter
	uiViewAction
	uiViewEnvironment

	owner SomeView
	types types

	size         size
	flexibleSize size
	modifiers    []uiViewModifier
	contents     []SomeView
}

func newView(types types, owner SomeView, contents ...View) *uiView {
	cts := make([]SomeView, 0, len(contents))
	for _, v := range contents {
		cts = append(cts, v.Body())
	}

	return &uiView{
		uiViewLayout: _zeroUIViewLayout,
		types:        types,
		owner:        owner,
		contents:     cts,
	}
}

/*
	########  #######  ########  ##     ## ##     ## ##          ###    ########
	##       ##     ## ##     ## ###   ### ##     ## ##         ## ##   ##     ##
	##       ##     ## ##     ## #### #### ##     ## ##        ##   ##  ##     ##
	######   ##     ## ########  ## ### ## ##     ## ##       ##     ## ########
	##       ##     ## ##   ##   ##     ## ##     ## ##       ######### ##   ##
	##       ##     ## ##    ##  ##     ## ##     ## ##       ##     ## ##    ##
	##        #######  ##     ## ##     ##  #######  ######## ##     ## ##     ##
*/

type preloadSizes struct {
	size    size
	minSize size
}

func (p *uiView) preloadSize() preloadSizes {
	l := logs.Default().
		WithField("types", p.types)

	preload := preloadSizes{
		minSize: p.getFrame(),
	}

	preload.minSize.w = rpEq(preload.minSize.w, -1, 0)
	preload.minSize.h = rpEq(preload.minSize.h, -1, 0)

	wSpacer, hSpacer := false, false
	for _, content := range p.contents {
		invokeSomeView(content, func(c *uiView) {
			p := c.preloadSize()
			l.Infof("preloadSize: size: %v", p.size)
			if p.size.w != -1 {
				preload.minSize.w += rpEq(p.size.w, -1, 0)
			} else {
				wSpacer = true

				preload.minSize.w += rpEq(p.minSize.w, -1, 0)
			}

			if p.size.h != -1 {
				preload.minSize.h += rpEq(p.size.h, -1, 0)
			} else {
				hSpacer = true
				preload.minSize.h += rpEq(p.minSize.h, -1, 0)
			}
		})
	}

	if wSpacer {
		preload.size.w = -1
	} else {
		preload.size.w = rpZero(preload.minSize.w, -1)
	}

	if hSpacer {
		preload.size.h = -1
	} else {
		preload.size.h = rpZero(preload.minSize.h, -1)
	}

	p.size = preload.size

	l.Infof("preloadSize: size: %v, preload: %v", p.size, preload)

	return preload
}

func (p *uiView) getFrame() size {
	frame := _zeroSize
	for i := len(p.modifiers) - 1; i >= 0 && frame.IsZero(); i-- {
		frame = p.modifiers[i].frame
	}

	return frame
}

func (p *uiView) setSizePosition(flexSize size, pos *point) {
	p.start.x = pos.x
	p.start.y = pos.y
	p.size.w = rpEq(p.size.w, -1, flexSize.w)
	p.size.h = rpEq(p.size.h, -1, flexSize.h)

	flexCount := p.getFlexibleCount()
	nextFlexSize := size{
		w: (p.size.w - flexCount.summedSize.w) / rpZero(flexCount.count.x, 1),
		h: (p.size.h - flexCount.summedSize.h) / rpZero(flexCount.count.y, 1),
	}

	var (
		postCache   = pos.Add(0, 0)
		afterPos    = func() {}
		contentSize size
	)
	for _, content := range p.contents {
		invokeSomeView(content, func(c *uiView) {
			contentSize = size{
				w: rpEq(c.size.w, -1, nextFlexSize.w),
				h: rpEq(c.size.h, -1, nextFlexSize.h),
			}
			switch p.types {
			case typesVStack:
				pos.x += (p.size.w - contentSize.w) / 2
				afterPos = func() {
					pos.x = postCache.x
					pos.y += contentSize.h
				}
			case typesHStack:
				pos.y += (p.size.h - contentSize.h) / 2
				afterPos = func() {
					pos.x += contentSize.w
					pos.y = postCache.y
				}
			case typesZStack:
			default:
				pos.x += (p.size.w - contentSize.w) / 2
				pos.y += (p.size.h - contentSize.h) / 2
				afterPos = func() {
					pos.x = postCache.x + contentSize.w
					pos.y = postCache.y + contentSize.h
				}
			}

			c.setSizePosition(nextFlexSize, pos)
			afterPos()
		})
	}

	logs.Default().
		WithField("types", p.types).
		Infof("setSizePosition: size: %v, start: %v, pos: %v, nextFlexSize: %v, lastPos: %v",
			p.size, p.start, pos, nextFlexSize, *pos)
}

type flexibleCount struct {
	count      point
	summedSize size
}

func (p *uiView) getFlexibleCount() flexibleCount {
	fc := flexibleCount{}
	for _, content := range p.contents {
		invokeSomeView(content, func(c *uiView) {
			if c.size.w == -1 {
				fc.count.x++
			} else {
				fc.summedSize.w += c.size.w
			}

			if c.size.h == -1 {
				fc.count.y++
			} else {
				fc.summedSize.h += c.size.h
			}
		})
	}

	return fc
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

func (p *uiView) first() *uiViewModifier {
	if len(p.modifiers) == 0 {
		p.pushFirst()
	}
	return &p.modifiers[0]
}

func (p *uiView) last() *uiViewModifier {
	if len(p.modifiers) == 0 {
		p.pushLast()
	}
	return &p.modifiers[len(p.modifiers)-1]
}

func (p *uiView) pushFirst(v ...uiViewModifier) {
	vv := uiViewModifier{uiViewLayout: _zeroUIViewLayout}
	if len(v) != 0 {
		vv = v[0]
	}

	p.modifiers = append(append(make([]uiViewModifier, 0, len(p.modifiers)*2), vv), p.modifiers...)
}

func (p *uiView) pushLast(v ...uiViewModifier) {
	vv := uiViewModifier{uiViewLayout: _zeroUIViewLayout}
	if len(v) != 0 {
		vv = v[0]
	}

	if len(p.modifiers) != 0 {
		anchor := p.modifiers[len(p.modifiers)-1]
		vv.frame = anchor.frame
		vv.margin = anchor.padding

	}

	p.modifiers = append(p.modifiers, vv)
}

func (p *uiView) Start(x, y int, replace ...bool) {
	if len(replace) != 0 && replace[0] {
		p.start.x = x
		p.start.y = y
	} else {
		p.start.x += x
		p.start.y += y
	}
}

func (p uiView) Copy() *uiView {
	return &p
}

func (p *uiView) GetFrameWithMargin() size {
	frame := _zeroSize
	margin := bounds{}
	for i := len(p.modifiers) - 1; i >= 0; i-- {
		if frame.IsZero() {
			frame = p.modifiers[i].frame
			if !margin.IsZero() {
				margin = p.modifiers[i].margin
			}
		}

		if frame.w != -1 || frame.h != -1 {
			break
		}
	}

	return frame.Add(margin.left+margin.right, margin.top+margin.bottom)
}

func (p *uiView) ActionUpdate() {
	// TODO: Update for actions
}

func (p *uiView) Draw(screen *ebiten.Image) {
	p.drawModifiers(screen)
	p.drawContent(screen)
}

func (p *uiView) drawModifiers(screen *ebiten.Image) {
	for _, v := range p.modifiers {
		w := rpEq(v.frame.w, -1, p.flexibleSize.w) - v.padding.left - v.padding.right
		h := rpEq(v.frame.h, -1, p.flexibleSize.h) - v.padding.top - v.padding.bottom

		if w <= 0 || h <= 0 {
			continue
		}

		img := ebiten.NewImage(w, h)
		img.Fill(rpZero(v.background, color.Color(color.Transparent)))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(p.start.x), float64(p.start.y))
		opt.GeoM.Translate(float64(v.offset.x), float64(v.offset.y))
		opt.GeoM.Translate(float64(v.margin.left), float64(v.margin.top))
		screen.DrawImage(img, opt)
	}
}

func (p *uiView) drawContent(screen *ebiten.Image) {
	for _, p := range p.contents {
		invokeSomeView(p, func(ui *uiView) {
			ui.Draw(screen)
		})
	}
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

/* Check Interface Implement */
var _ uiViewDelegator = (*uiView)(nil)

func (p *uiView) UIView() *uiView {
	return p
}

/* Check Interface Implementation */
var _ SomeView = (*uiView)(nil)

func (p *uiView) Body() SomeView {
	return p.owner
}

func (p *uiView) Frame(w, h int) SomeView {
	if w < 0 {
		w = -1
	}

	if h < 0 {
		h = -1
	}

	last := p.last()
	last.frame = size{w, h}
	last.padding = bounds{}
	return p.owner
}

func (p *uiView) ForegroundColor(clr color.Color) SomeView {
	return p.owner
}

func (p *uiView) BackgroundColor(clr color.Color) SomeView {
	last := p.last()
	last.background = rpZero(last.background, clr)

	// add frame to front views
	w := last.frame.w
	if w != -1 {
		for i := len(p.modifiers) - 2; i >= 0; i-- {
			if p.modifiers[i].frame.w != -1 {
				break
			}

			p.modifiers[i].frame.w = w
		}
	}

	h := last.frame.h
	if h != -1 {
		for i := len(p.modifiers) - 2; i >= 0; i-- {
			if p.modifiers[i].frame.h != -1 {
				break
			}

			p.modifiers[i].frame.h = h
		}
	}

	p.pushLast()

	return p.owner
}

func (p *uiView) Padding(paddings ...int) SomeView {
	top, right, bottom, left := 0, 0, 0, 0

	switch len(paddings) {
	case 1: /* all */
		top = paddings[0]
		right = paddings[0]
		bottom = paddings[0]
		left = paddings[0]
	case 2: /* vertical, horizontal */
		top = paddings[0]
		right = paddings[1]
		bottom = paddings[0]
		left = paddings[1]
	case 4: /* top, right, bottom, left */
		top = paddings[0]
		right = paddings[1]
		bottom = paddings[2]
		left = paddings[3]
	}

	if top < 0 {
		top = 0
	}

	if right < 0 {
		right = 0
	}

	if bottom < 0 {
		bottom = 0
	}

	if left < 0 {
		left = 0
	}

	last := p.last()
	if last.frame.w != -1 || last.frame.h != -1 {
		last.margin = last.margin.Add(top, right, bottom, left)
	} else {
		last.padding = last.padding.Add(top, right, bottom, left)
	}

	return p.owner
}

func (p *uiView) Offset(x, y int) SomeView {
	p.offset.x += x
	p.offset.y += y
	return p.owner
}

func (p *uiView) FontSize(size font.Size) SomeView {
	p.fontSize = size
	return p.owner
}

func (p *uiView) FontWeight(weight font.Weight) SomeView {
	p.fontWeight = weight
	return p.owner
}

func (p *uiView) CornerRadius(radius ...int) SomeView {
	p.cornerRadius = sliceFirst(radius, 7)
	return p.owner
}

func (p *uiView) LineSpacing(spacing float64) SomeView {
	p.lineSpacing = spacing
	return p.owner
}

func (p *uiView) Kerning(kerning int) SomeView {
	p.kerning = kerning
	return p.owner
}

func (p *uiView) Italic() SomeView {
	p.italic = true
	return p.owner
}
