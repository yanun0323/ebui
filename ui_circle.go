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

func (c *circleImpl) userSetFrameSize() CGSize {
	frameSize := c.viewCtx.userSetFrameSize()
	frameSize.Width = min(frameSize.Width, frameSize.Height)
	frameSize.Height = frameSize.Width

	return frameSize
}

func (c *circleImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	drawFrame := c._owner.systemSetBounds()

	bOpt := c.drawOption(drawFrame, hook...)

	if c.backgroundColor == nil {
		return
	}

	bgColor := c.backgroundColor.Get()
	if !drawFrame.drawable() {
		return
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
	opt.GeoM.Concat(bOpt.GeoM)
	opt.ColorScale.ScaleWithColorScale(bOpt.ColorScale)
	screen.DrawImage(img, opt)
}
