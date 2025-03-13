package ebui

import (
	"bytes"
	"math"
	"strconv"

	"github.com/cespare/xxhash/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/internal/helper"
	"github.com/yanun0323/ebui/layout"
)

var _ SomeView = &viewCtx{}

var Inf = math.Inf(1)
var isInf = func(f float64) bool {
	return math.IsInf(f, 1)
}

// viewCtx provides the common modifiers and states for all views
type viewCtx struct {
	*viewCtxEnv
	*viewCtxParam

	_owner SomeView
	_cache *helper.HashCache[*ebiten.Image]
}

// initialize ViewContext
func newViewContext(owner SomeView) *viewCtx {
	return &viewCtx{
		viewCtxEnv:   newEnv(),
		viewCtxParam: newParam(),
		_owner:       owner,
		_cache:       helper.NewHashCache[*ebiten.Image](),
	}
}

func (c *viewCtx) Bytes(withFont bool) []byte {
	b := bytes.Buffer{}
	b.Write(c.viewCtxEnv.Bytes(withFont))
	b.Write(c.viewCtxParam.Bytes())

	return b.Bytes()
}

func (c *viewCtx) Hash(withFont bool) string {
	h := xxhash.New()
	h.Write(c.viewCtxEnv.Bytes(withFont))
	h.Write(c.viewCtxParam.Bytes())
	return strconv.FormatUint(h.Sum64(), 16)
}

func (c *viewCtx) wrap(modify ...func(*viewCtx)) SomeView {
	// create a new zstackImpl instance
	zs := &zstackImpl{
		children: []SomeView{c._owner},
	}

	zs.viewCtx = newViewContext(zs)
	zs.viewCtx.viewCtxEnv = c.viewCtxEnv
	// zs.viewCtx.viewCtxParam.frameSize = c.frameSize

	// apply modifiers
	for _, m := range modify {
		m(zs.viewCtx)
	}

	return zs
}

func (c *viewCtx) setEnv(env *viewCtxEnv) {
	c.viewCtxEnv = env
}

func (c *viewCtx) count() int {
	return 1
}

func (c *viewCtx) align(offset CGPoint) {
	c._systemSetFrame = c._systemSetFrame.Move(offset)
}

func (c *viewCtx) preload(parent *viewCtxEnv, _ ...formulaType) (preloadData, layoutFunc) {
	c.viewCtxEnv.inheritFrom(parent)
	padding := c.padding()
	border := c.border()
	userSetFrameSize := c._owner.userSetFrameSize()
	data := newPreloadData(userSetFrameSize, padding, border)
	return data, func(start CGPoint, flexBoundsSize CGSize) (CGRect, alignFunc) {
		flexFrameSize := flexBoundsSize.Shrink(padding).Shrink(border)
		flexibleFrame := CGRect{start, start.Add(flexFrameSize.ToCGPoint())}
		finalFrame := flexibleFrame
		if !userSetFrameSize.IsInfWidth() {
			finalFrame.End.X = start.X + userSetFrameSize.Width
		}

		if !userSetFrameSize.IsInfHeight() {
			finalFrame.End.Y = start.Y + userSetFrameSize.Height
		}

		c._systemSetFrame = NewRect(
			finalFrame.Start.X+padding.Left+border.Left,
			finalFrame.Start.Y+padding.Top+border.Top,
			finalFrame.End.X+padding.Left+border.Left,
			finalFrame.End.Y+padding.Top+border.Top,
		)

		c._cache.SetNextHash(c.Hash(false))
		c.debugPrintPreload(finalFrame, flexFrameSize, data)

		return finalFrame.Expand(padding).Expand(border), c.align
	}
}

func (c *viewCtx) drawOption(rect CGRect, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(rect.Start.X, rect.Start.Y)
	for _, h := range hook {
		h(opt)
	}

	if c.opacity != nil {
		opt.ColorScale.ScaleAlpha(float32(c.opacity.Get()))
	}
	return opt
}

func (c *viewCtx) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	drawBounds := c._owner.systemSetBounds()

	if !drawBounds.drawable() {
		return
	}

	bgColor := c.backgroundColor.Get()
	borderLength := c.borderInset.Get()
	borderColor := black
	if c.borderColor != nil {
		borderColor = c.borderColor.Get()
	}

	opt := c.drawOption(drawBounds, hook...)
	if radius := c.roundCorner.Get(); radius > 0 {
		c.drawRoundedAndBorderRect(screen, drawBounds, radius, bgColor, borderLength, borderColor, opt)
	} else {
		c.drawBorderRect(screen, drawBounds, bgColor, borderLength, borderColor, opt)
	}
}

/*
	####   ##     ##   ########    ##         ########   ##     ##   ########   ##    ##   ########
	 ##    ###   ###   ##     ##   ##         ##         ###   ###   ##         ###   ##      ##
	 ##    #### ####   ##     ##   ##         ##         #### ####   ##         ####  ##      ##
	 ##    ## ### ##   ########    ##         ######     ## ### ##   ######     ## ## ##      ##
	 ##    ##     ##   ##          ##         ##         ##     ##   ##         ##  ####      ##
	 ##    ##     ##   ##          ##         ##         ##     ##   ##         ##   ###      ##
	####   ##     ##   ##          ########   ########   ##     ##   ########   ##    ##      ##
*/

func (c *viewCtx) Body() SomeView {
	return c._owner
}

func (c *viewCtx) Frame(size *Binding[CGSize]) SomeView {
	if size == nil {
		size = Const(NewSize(Inf, Inf))
	}

	c.frameSize = size

	return c._owner
}

func (c *viewCtx) Padding(padding ...*Binding[CGInset]) SomeView {
	if len(padding) != 0 {
		return c.wrap(func(c *viewCtx) {
			c.inset = padding[0]
		})
	} else {
		return c.wrap(func(c *viewCtx) {
			c.inset = DefaultPadding
		})
	}
}

func (c *viewCtx) ForegroundColor(color *Binding[CGColor]) SomeView {
	c.foregroundColor = color
	return c._owner
}

func (c *viewCtx) BackgroundColor(color *Binding[CGColor]) SomeView {
	c.backgroundColor = color
	return c._owner
}

func (c *viewCtx) FontSize(size *Binding[font.Size]) SomeView {
	c.fontSize = size
	return c._owner
}

func (c *viewCtx) FontWeight(weight *Binding[font.Weight]) SomeView {
	c.fontWeight = weight
	return c._owner
}

func (c *viewCtx) FontLineHeight(height *Binding[float64]) SomeView {
	c.fontLineHeight = height
	return c._owner
}

func (c *viewCtx) FontKerning(spacing *Binding[float64]) SomeView {
	c.fontKerning = spacing
	return c._owner
}

func (c *viewCtx) FontAlignment(alignment *Binding[font.TextAlign]) SomeView {
	c.fontAlignment = alignment
	return c._owner
}

func (c *viewCtx) FontItalic(italic ...*Binding[bool]) SomeView {
	if len(italic) != 0 {
		c.fontItalic = italic[0]
	} else {
		c.fontItalic = Const(true)
	}

	return c._owner
}

func (c *viewCtx) RoundCorner(radius ...*Binding[float64]) SomeView {
	if len(radius) != 0 {
		c.roundCorner = radius[0]
	} else {
		c.roundCorner = DefaultRoundCorner
	}
	return c.wrap()
}

func (c *viewCtx) DebugPrint(tag string) SomeView {
	c._debug = tag
	return c._owner
}

func (c *viewCtx) ScaleToFit(enable ...*Binding[bool]) SomeView {
	if len(enable) != 0 {
		c.scaleToFit = enable[0]
	} else {
		c.scaleToFit = Const(true)
	}

	return c._owner
}

func (c *viewCtx) KeepAspectRatio(enable ...*Binding[bool]) SomeView {
	if len(enable) != 0 {
		c.keepAspectRatio = enable[0]
	} else {
		c.keepAspectRatio = Const(true)
	}

	return c._owner
}

func (c *viewCtx) Border(border *Binding[CGInset], color ...*Binding[CGColor]) SomeView {
	c.borderInset = border
	if len(color) != 0 {
		c.borderColor = color[0]
	} else {
		c.borderColor = DefaultBorderColor
	}

	return c._owner
}

func (c *viewCtx) Opacity(opacity *Binding[float64]) SomeView {
	c.opacity = opacity
	return c._owner
}

func (c *viewCtx) Modifier(modifier ViewModifier) SomeView {
	return modifier.Body(c)
}

func (c *viewCtx) Modify(with func(SomeView) SomeView) SomeView {
	return with(c)
}

func (c *viewCtx) Disabled(disabled ...*Binding[bool]) SomeView {
	if len(disabled) != 0 {
		c.disabled = disabled[0]
	} else {
		c.disabled = Const(true)
	}

	return c._owner
}

func (c *viewCtx) Align(alignment *Binding[layout.Align]) SomeView {
	c.alignment = alignment

	if alignment == nil {
		return c._owner
	}

	if c.transitionAlign == nil {
		c.transitionAlign = Bind(CGPoint{})
	}

	c.transitionAlign.Set(alignToCGPoint(alignment.Get()), nil)
	alignment.addListener(func(oldVal, newVal layout.Align, animStyle animation.Style, isAnimating bool) {
		println("align changed", "old", serialize(oldVal), "new", serialize(newVal), "isAnimating", serialize(isAnimating))
		if isAnimating {
			return
		}
		c.transitionAlign.Set(alignToCGPoint(newVal), animStyle)
	})

	return c._owner
}

func (c *viewCtx) Debug() SomeView {
	c.Border(Const(NewInset(1)), Const(NewColor(255, 0, 0)))
	return c.wrap()
}

func (c *viewCtx) Center() SomeView {
	return VStack(
		Spacer(),
		HStack(
			Spacer(),
			c._owner,
			Spacer(),
		),
		Spacer(),
	)
}

func (c *viewCtx) Offset(offset *Binding[CGPoint]) SomeView {
	c.offset = offset
	return c._owner
}

func (c *viewCtx) Spacing(spacing ...*Binding[float64]) SomeView {
	if len(spacing) != 0 {
		c.spacing = spacing[0]
	} else {
		c.spacing = DefaultSpacing
	}

	return c._owner
}
