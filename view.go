package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
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

	flexible  size
	modifiers []uiViewModifier
	contents  []SomeView
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
	rootStart := loadRootStart()
	for _, v := range p.modifiers {
		w := rpEq(v.frame.w, -1, p.flexible.w) - v.padding.left - v.padding.right
		h := rpEq(v.frame.h, -1, p.flexible.h) - v.padding.top - v.padding.bottom

		if w <= 0 || h <= 0 {
			continue
		}

		img := ebiten.NewImage(w, h)
		img.Fill(rpZero(v.background, color.Color(color.Transparent)))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(rootStart.x), float64(rootStart.y))
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
