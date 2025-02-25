package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type circleImpl struct {
	*viewCtx
}

func Circle() SomeView {
	circle := &circleImpl{}
	circle.viewCtx = newViewContext(circle)
	return circle
}

func (c *circleImpl) userSetFrameSize() flexibleSize {
	frameSize := c.viewCtx.userSetFrameSize()
	frameSize.Frame = NewSize(
		min(frameSize.Frame.Width, frameSize.Frame.Height),
		min(frameSize.Frame.Width, frameSize.Frame.Height),
	)

	return frameSize
}

func (c *circleImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	drawFrame := c._owner.systemSetBounds()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(drawFrame.Start.X, drawFrame.Start.Y)
	for _, h := range hook {
		h(op)
	}

	bgColor := c.backgroundColor.Get()
	if bgColor == nil {
		return op
	}

	if !drawFrame.drawable() {
		return op
	}

	w := int(drawFrame.Dx() * _roundedScale)
	h := int(drawFrame.Dy() * _roundedScale)
	diameter := min(w, h)
	radius := diameter / 2

	img := ebiten.NewImage(diameter, diameter)
	img.Fill(bgColor)

	for x := 0; x < diameter; x++ {
		for y := 0; y < diameter; y++ {
			if (x-radius)*(x-radius)+(y-radius)*(y-radius) > radius*radius {
				img.Set(x, y, color.Transparent)
			}
		}
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(_roundedScaleInverse, _roundedScaleInverse)
	opt.Filter = ebiten.FilterLinear
	opt.GeoM.Concat(op.GeoM)
	opt.ColorScale.ScaleWithColorScale(op.ColorScale)
	screen.DrawImage(img, opt)

	return op
}
