package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/sys"
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

	Resizable() SomeView
	AspectRatio(ratio ...float64) SomeView

	// PRIVATE
	draw(screen *ebiten.Image)
	getTypes() types
	// reset 清除 getSize 的 cache
	//
	// 在每次 ebiten update 的時候調用
	reset()
	// getSize 取得這個視圖的大小並快取，（如果沒設定大小，就會計算子視圖的總大小）
	getSize() size
	setSize(size)
	getPosition() point
	setPosition(point)
	// stepSubView 設定 Stack 子視圖的位置，與子視圖大小的移動關係
	//
	// e.g. VStack：x 不會被子視圖大小影響；y 會加上「子視圖的高度」單位
	stepSubView(pos point, childSize size) point

	// getStackSubViewStart 取得 Stack 視圖，子視圖的起始位置
	getStackSubViewStart(offset point) point
	// setStackSubViewCenterOffset 設定子視圖的中心偏移
	getStackSubViewCenterOffset(offset point) point
	subView() []SomeView
	update()
	isPress(x, y int) bool
	setEnvironment(uiViewEnvironment)
	handlePreference(*ebiten.DrawImageOptions)
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

	isCached   bool
	cachedSize size
	modifiers  []uiViewModifier
	contents   []SomeView
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
	########  ########  #### ##     ##    ###    ######## ########
	##     ## ##     ##  ##  ##     ##   ## ##      ##    ##
	##     ## ##     ##  ##  ##     ##  ##   ##     ##    ##
	########  ########   ##  ##     ## ##     ##    ##    ######
	##        ##   ##    ##   ##   ##  #########    ##    ##
	##        ##    ##   ##    ## ##   ##     ##    ##    ##
	##        ##     ## ####    ###    ##     ##    ##    ########
*/

func (p *uiView) getFrame() size {
	frame := _zeroSize
	for i := len(p.modifiers) - 1; i >= 0 && frame.IsZero(); i-- {
		frame = p.modifiers[i].frame
	}

	return frame
}

func (p *uiView) last() *uiViewModifier {
	if len(p.modifiers) == 0 {
		p.pushLast()
	}
	return &p.modifiers[len(p.modifiers)-1]
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

func (p *uiView) actionUpdate() {
	// TODO: Update for actions
}

/*
	########  ########     ###    ##      ##
	##     ## ##     ##   ## ##   ##  ##  ##
	##     ## ##     ##  ##   ##  ##  ##  ##
	##     ## ########  ##     ## ##  ##  ##
	##     ## ##   ##   ######### ##  ##  ##
	##     ## ##    ##  ##     ## ##  ##  ##
	########  ##     ## ##     ##  ###  ###
*/

func (p *uiView) draw(screen *ebiten.Image) {
	p.drawModifiers(screen)
	p.drawContent(screen)
}

func (p *uiView) drawModifiers(screen *ebiten.Image) {
	for _, v := range p.modifiers {
		pos := v.getDrawSize(p.cachedSize)
		if pos.w <= 0 || pos.h <= 0 {
			continue
		}

		img := ebiten.NewImage(pos.w, pos.h)
		img.Fill(sys.If(v.background == nil, color.Color(color.Transparent), v.background))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(p.start.x), float64(p.start.y))
		opt.GeoM.Translate(float64(v.offset.x), float64(v.offset.y))
		opt.GeoM.Translate(float64(v.margin.left), float64(v.margin.top))
		p.handlePreference(opt)
		screen.DrawImage(img, opt)
	}
}

func (p *uiView) drawContent(screen *ebiten.Image) {
	for _, p := range p.contents {
		p.draw(screen)
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
	last.background = sys.If(last.background == nil, clr, last.background)

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
	p.cornerRadius = first(radius, 7)
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

func (p *uiView) Resizable() SomeView {
	p.resizable = true
	return p.owner
}

func (p *uiView) AspectRatio(ratio ...float64) SomeView {
	if len(ratio) != 0 {
		p.aspectRatio = ratio[0]
	} else {
		p.aspectRatio = 1
	}

	return p.owner
}

func (p *uiView) getTypes() types {
	return p.types
}

func (p *uiView) reset() {
	p.isCached = false
	p.isPressing = false
}

func (p *uiView) getSize() size {
	if p.isCached {
		return p.cachedSize
	}

	size := p.getFrame()
	if size.w != -1 && size.h != -1 {
		p.isCached = true
		p.cachedSize = size
		return size
	}

	result := _zeroSize
	childNoWidthCount := 0
	childNoHeightCount := 0
	for _, child := range p.contents {
		childSize := child.getSize()
		result.w = max(result.w, childSize.w)
		result.h = max(result.h, childSize.h)
		childNoWidthCount += sys.If(childSize.w >= 0, 0, 1)
		childNoHeightCount += sys.If(childSize.h >= 0, 0, 1)
	}

	result.w = sys.If(size.w == -1, result.w, size.w)
	result.h = sys.If(size.h == -1, result.h, size.h)
	result.w = sys.If(childNoWidthCount != 0, -1, result.w)
	result.h = sys.If(childNoHeightCount != 0, -1, result.h)

	p.isCached = true
	p.cachedSize = result
	return result
}

func (p *uiView) setSize(size size) {
	p.isCached = true
	p.cachedSize = size
}

func (p *uiView) getPosition() point {
	return p.start
}

func (p *uiView) setPosition(pos point) {
	p.start = pos
}

func (p *uiView) getStackSubViewStart(offset point) point {
	return point{}
}

func (p *uiView) getStackSubViewCenterOffset(offset point) point {
	return offset
}

func (p *uiView) stepSubView(pos point, childSize size) point {
	return pos
}

func (p *uiView) subView() []SomeView {
	return p.contents
}

func (p *uiView) update() {
	p.actionUpdate()
	for _, v := range p.contents {
		v.setEnvironment(p.uiViewEnvironment)
	}
}

func (p *uiView) isPress(x, y int) bool {
	drawSize := p.getDrawSize(p.cachedSize)
	start := p.start

	return x >= start.x && x <= start.x+drawSize.w && y >= start.y && y <= start.y+drawSize.h
}

func (p *uiView) setEnvironment(env uiViewEnvironment) {
	p.uiViewEnvironment.set(env)
}

func (p *uiView) handlePreference(opt *ebiten.DrawImageOptions) {
	if p.isPressing {
		opt.ColorScale.ScaleAlpha(0.3)
	}
}
