package ebui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
)

var _ SomeView = &viewCtx{}

var Inf = math.Inf(1)
var isInf = func(f float64) bool {
	return math.IsInf(f, 1)
}

const (
	_defaultRoundCorner float64 = 8.0
)

// viewCtx 提供所有 View 共用的修飾器方法和狀態
type viewCtx struct {
	*viewCtxEnv
	*viewCtxParam

	_owner SomeView
}

// 初始化 ViewContext
func newViewContext(owner SomeView) *viewCtx {
	return &viewCtx{
		viewCtxEnv:   newEnv(),
		viewCtxParam: newParam(),
		_owner:       owner,
	}
}

func (c *viewCtx) wrap(modify func(*viewCtx)) SomeView {
	// 創建一個新的 zstackImpl 實例
	zs := &zstackImpl{
		children: []SomeView{c._owner},
	}

	zs.viewCtx = newViewContext(zs)
	zs.viewCtx.viewCtxEnv = c.viewCtxEnv
	zs.viewCtx.viewCtxParam.frameSize = c.frameSize

	// 應用修改
	modify(zs.viewCtx)

	return zs
}

func (c *viewCtx) setEnv(env *viewCtxEnv) {
	c.viewCtxEnv = env
}

func (c *viewCtx) count() int {
	return 1
}

func (c *viewCtx) preload(parent *viewCtxEnv) (flexibleSize, CGInset, layoutFunc) {
	c.viewCtxEnv.inheritFrom(parent)
	padding := c.inset.Get()
	userSetFrameSize := c._owner.userSetFrameSize()
	return userSetFrameSize, padding, func(start CGPoint, flexFrameSize CGSize) CGRect {
		finalFrame := CGRect{start, start.Add(flexFrameSize.ToCGPoint())}
		finalFrameSize := userSetFrameSize
		if !finalFrameSize.IsInfX {
			finalFrame.End.X = start.X + finalFrameSize.Frame.Width
		}

		if !finalFrameSize.IsInfY {
			finalFrame.End.Y = start.Y + finalFrameSize.Frame.Height
		}

		c._systemSetFrame = NewRect(
			finalFrame.Start.X+padding.Left,
			finalFrame.Start.Y+padding.Top,
			finalFrame.End.X+padding.Left,
			finalFrame.End.Y+padding.Top,
		)

		result := finalFrame.Expand(padding)
		c.debugPrint(result)
		return result
	}
}

func (c *viewCtx) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	drawFrame := c._owner.systemSetBounds()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(drawFrame.Start.X, drawFrame.Start.Y)
	for _, h := range hook {
		h(opt)
	}

	bgColor := c.backgroundColor.Get()
	if bgColor == nil {
		return opt
	}

	if !drawFrame.drawable() {
		return opt
	}

	if radius := c.roundCorner.Get(); radius > 0 {
		drawRoundedRect(screen, drawFrame, radius, bgColor, opt)
	} else {
		img := ebiten.NewImage(int(drawFrame.Dx()), int(drawFrame.Dy()))
		img.Fill(bgColor)
		screen.DrawImage(img, opt)
	}

	return opt
}

func (c *viewCtx) Body() SomeView {
	return c._owner
}

// Frame 修飾器
func (c *viewCtx) Frame(size *Binding[CGSize]) SomeView {
	if size == nil {
		size = Bind(NewSize(Inf, Inf))
	}

	c.frameSize = size

	return c._owner
}

// Padding 修飾器
func (c *viewCtx) Padding(padding *Binding[CGInset]) SomeView {
	return c.wrap(func(c *viewCtx) {
		c.inset = padding
	})
}

// ForegroundColor 修飾器
func (c *viewCtx) ForegroundColor(color *Binding[AnyColor]) SomeView {
	c.foregroundColor = color
	return c._owner
}

// BackgroundColor 修飾器
func (c *viewCtx) BackgroundColor(color *Binding[AnyColor]) SomeView {
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

func (c *viewCtx) FontLetterSpacing(spacing *Binding[float64]) SomeView {
	c.fontLetterSpacing = spacing
	return c._owner
}

func (c *viewCtx) FontAlignment(alignment *Binding[font.Alignment]) SomeView {
	c.fontAlignment = alignment
	return c._owner
}

func (c *viewCtx) FontItalic(italic ...*Binding[bool]) SomeView {
	if len(italic) != 0 {
		c.fontItalic = italic[0]
	} else {
		c.fontItalic.Set(true)
	}

	return c._owner
}

func (c *viewCtx) RoundCorner(radius ...*Binding[float64]) SomeView {
	if len(radius) != 0 {
		c.roundCorner = radius[0]
	} else {
		c.roundCorner = Bind(_defaultRoundCorner)
	}
	return c._owner
}

func (c *viewCtx) Debug(tag string) SomeView {
	c._debug = tag
	return c._owner
}

func (c *viewCtx) ScaleToFit(enable ...*Binding[bool]) SomeView {
	if len(enable) != 0 {
		c.scaleToFit = enable[0]
	} else {
		c.scaleToFit = Bind(true)
	}

	return c._owner
}

func (c *viewCtx) KeepAspectRatio(enable ...*Binding[bool]) SomeView {
	if len(enable) != 0 {
		c.keepAspectRatio = enable[0]
	} else {
		c.keepAspectRatio = Bind(true)
	}

	return c._owner
}
