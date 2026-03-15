package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (c *viewCtx) loadBackgroundImage(bounds CGRect, radius float64, bgColor CGColor, border CGInset, borderColor CGColor) *ebiten.Image {
	if c._cache.IsNextHashCached() {
		return c._cache.Load()
	}

	img := createRoundedRect(bounds, radius, bgColor, border, borderColor)
	c._cache.Update(img)
	return img
}

func createRoundedRect(bounds CGRect, radius float64, bgColor CGColor, border CGInset, borderColor CGColor) *ebiten.Image {
	var (
		w   = bounds.Dx()
		h   = bounds.Dy()
		img = ebiten.NewImage(int(w), int(h))
	)

	if radius <= 0 {
		// draw border
		if !border.IsZero() {
			vector.DrawFilledRect(img, 0, 0, float32(w), float32(border.Top), borderColor, true)                           // top
			vector.DrawFilledRect(img, float32(w-border.Right), 0, float32(border.Right), float32(h), borderColor, true)   // right
			vector.DrawFilledRect(img, 0, float32(h-border.Bottom), float32(w), float32(border.Bottom), borderColor, true) // bottom
			vector.DrawFilledRect(img, 0, 0, float32(border.Left), float32(h), borderColor, true)                          // left
		}
		// draw background
		vector.DrawFilledRect(img, float32(border.Left), float32(border.Top), float32(w-border.Left-border.Right), float32(h-border.Top-border.Bottom), bgColor, true)

		return img
	}

	drawRoundedRect := func(img *ebiten.Image, x, y, width, height float32, r float32, clr CGColor) {
		r = min(min(r, width/2), height/2)
		r = max(r, 0)

		isCircle := width == height && r >= width/2
		if isCircle {
			vector.DrawFilledCircle(img, x+width/2, y+height/2, r, clr, true)
			return
		}

		var (
			minX          = x + r
			minY          = y + r
			maxX          = x + width - r
			maxY          = y + height - r
			diameter      = r * 2
			inboundWidth  = width - diameter
			inboundHeight = height - diameter
		)

		// draw corner circle
		vector.DrawFilledCircle(img, minX, minY, r, clr, true) // top left
		vector.DrawFilledCircle(img, maxX, minY, r, clr, true) // top right
		vector.DrawFilledCircle(img, minX, maxY, r, clr, true) // bottom left
		vector.DrawFilledCircle(img, maxX, maxY, r, clr, true) // bottom right

		// draw sides
		vector.DrawFilledRect(img, minX, y, inboundWidth, height, clr, true) // top bottom
		vector.DrawFilledRect(img, x, minY, width, inboundHeight, clr, true) // left right

	}

	if !border.IsZero() { // draw border
		drawRoundedRect(img, 0, 0, float32(w), float32(h), float32(radius), borderColor)
	}

	// draw background
	drawRoundedRect(img, float32(border.Left), float32(border.Top), float32(w-border.Left-border.Right), float32(h-border.Top-border.Bottom), float32(radius), bgColor)

	return img
}

// draw a shadow
func (c *viewCtx) drawShadow(screen, source *ebiten.Image, shadowLength float64, shadowColor CGColor, bOpt *ebiten.DrawImageOptions) {
	if shadowLength <= 0 {
		return
	}

	var (
		scale    = 2.0
		blend    = ebiten.BlendSourceOver
		finalOpt = &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
			Blend:  blend,
		}
	)

	shadowLength /= scale

	finalOpt.GeoM.Translate(-shadowLength, -shadowLength)
	finalOpt.GeoM.Scale(scale, scale)
	finalOpt.GeoM.Translate(shadowLength, shadowLength)
	finalOpt.GeoM.Concat(bOpt.GeoM)
	finalOpt.ColorScale.ScaleWithColor(shadowColor)
	finalOpt.GeoM.Translate(-shadowLength, -shadowLength)

	if c._shadowCache.IsNextHashCached() {
		screen.DrawImage(c._shadowCache.Load(), finalOpt)
		return
	}

	var (
		dl        = 2 * int(ceil(shadowLength))
		w         = source.Bounds().Dx() / int(ceil(scale))
		h         = source.Bounds().Dy() / int(ceil(scale))
		tempImg   = ebiten.NewImage(w+dl, h)
		shadowImg = ebiten.NewImage(w+dl, h+dl)
		step      = 16
		offset    = func(i int) float64 {
			return float64(i) * shadowLength / float64(step)
		}
		opt = &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
			Blend:  blend,
		}
		alpha             = 1 / float32(step)
		baseAlpha float32 = 0.01

		scaleReverse = 1 / scale
	)

	{ // horizontal
		for i := range step {
			opt.ColorScale.Reset()
			opt.ColorScale.Scale(0, 0, 0, alpha*float32(step-i)/float32(step)+baseAlpha)
			{
				opt.GeoM.Reset()
				opt.GeoM.Scale(scaleReverse, scaleReverse)
				opt.GeoM.Translate(shadowLength-offset(i), 0)
				tempImg.DrawImage(source, opt)

				opt.GeoM.Reset()
				opt.GeoM.Scale(scaleReverse, scaleReverse)
				opt.GeoM.Translate(shadowLength+offset(i), 0)
				tempImg.DrawImage(source, opt)
			}
		}
	}

	{ // vertical
		for i := range step {
			if i == 0 {
				continue
			}

			opt.ColorScale.Reset()
			opt.ColorScale.Scale(0, 0, 0, alpha*float32(step-i)/float32(step)+baseAlpha)
			{
				opt.GeoM.Reset()
				opt.GeoM.Translate(0, shadowLength-offset(i))
				shadowImg.DrawImage(tempImg, opt)

				opt.GeoM.Reset()
				opt.GeoM.Translate(0, shadowLength+offset(i))
				shadowImg.DrawImage(tempImg, opt)
			}
		}
	}

	c._shadowCache.Update(shadowImg)
	screen.DrawImage(shadowImg, finalOpt)
}
