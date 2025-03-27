package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// draw a rounded and bordered rectangle
func (c *viewCtx) drawRoundedAndBorderRect(screen *ebiten.Image, bounds CGRect, radius float64, bgColor CGColor, border CGInset, borderColor CGColor, bOpt *ebiten.DrawImageOptions) {
	var img *ebiten.Image
	if c._cache.IsNextHashCached() {
		img = c._cache.Load()
	} else {
		w := int(bounds.Dx() * _roundedScale)
		h := int(bounds.Dy() * _roundedScale)
		r := (radius + border.Top) * _roundedScale // FIXME: using inset to calculate border
		b := border.Top * _roundedScale
		img = ebiten.NewImage(w, h)
		if bgColor != transparent {
			img.Fill(bgColor)
		}

		cornerHandler := newCornerHandler(w, h, r, b)
		cornerHandler.Execute(func(isOutside, isBorder bool, x, y int) {
			if isOutside {
				img.Set(x, y, color.Transparent)
				return
			}

			if isBorder && b > 0 && borderColor != transparent {
				img.Set(x, y, borderColor)
				return
			}
		})

		c._cache.Update(img)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(_roundedScaleInverse, _roundedScaleInverse)
	opt.Filter = ebiten.FilterLinear
	opt.GeoM.Concat(bOpt.GeoM)
	opt.ColorScale.ScaleWithColorScale(bOpt.ColorScale)

	screen.DrawImage(img, opt)
}

// draw a bordered rectangle
func (c *viewCtx) drawBorderRect(screen *ebiten.Image, bounds CGRect, bgColor CGColor, border CGInset, borderColor CGColor, bOpt *ebiten.DrawImageOptions) {
	var img *ebiten.Image
	if c._cache.IsNextHashCached() {
		img = c._cache.Load()
	} else {
		w := int(bounds.Dx())
		h := int(bounds.Dy())
		img = ebiten.NewImage(w, h)
		if bgColor != transparent {
			img.Fill(bgColor)
		}

		var (
			left     = int(border.Left)
			top      = int(border.Top)
			right    = w - int(border.Right)
			bottom   = h - int(border.Bottom)
			isBorder = func(x, y int) bool {
				return x < left || y < top || x >= right || y >= bottom
			}
		)

		if !border.IsZero() {
			for x := range w {
				for y := range h {
					if isBorder(x, y) {
						img.Set(x, y, borderColor)
						continue
					}
				}
			}
		}

		c._cache.Update(img)
	}

	screen.DrawImage(img, bOpt)
}

// draw a shadow
func (c *viewCtx) drawShadow(screen, source *ebiten.Image, shadowLength float64, shadowColor CGColor, bOpt *ebiten.DrawImageOptions, fix ...float64) {
	if shadowLength <= 0 {
		return
	}

	fixOffset := 0.0
	if len(fix) != 0 && fix[0] != 0 {
		fixOffset = fix[0]
	}
	fixOffset *= shadowLength

	blend := ebiten.BlendSourceOver
	finalOpt := &ebiten.DrawImageOptions{}
	finalOpt.Filter = ebiten.FilterLinear
	finalOpt.Blend = blend
	// finalOpt.GeoM.Scale(1.111, 1.111)
	finalOpt.GeoM.Concat(bOpt.GeoM)
	finalOpt.ColorScale.ScaleWithColorScale(bOpt.ColorScale)
	finalOpt.GeoM.Translate(-shadowLength+fixOffset, -shadowLength+fixOffset)

	if c._shadowCache.IsNextHashCached() {
		screen.DrawImage(c._shadowCache.Load(), finalOpt)
		return
	}

	w := int(source.Bounds().Dx())
	h := int(source.Bounds().Dy())

	shadowBlock := ebiten.NewImage(w, h)
	sOpt := &ebiten.DrawImageOptions{}
	sOpt.ColorScale.Scale(0, 0, 0, 1)
	sOpt.ColorScale.ScaleWithColor(shadowColor)
	shadowBlock.DrawImage(source, sOpt)

	count := int(abs(shadowLength))
	w += 2 * count
	tempImg := ebiten.NewImage(w, h)

	h += 2 * count
	shadowImg := ebiten.NewImage(w, h)

	alpha := 1.0 / float32(count)
	// 步驟1：水平模糊
	for i := -count; i <= count; i++ {
		opt := &ebiten.DrawImageOptions{}
		opt.Filter = ebiten.FilterLinear
		opt.Blend = blend
		opt.GeoM.Translate(float64(i)+shadowLength, 0.0)
		opt.ColorScale.ScaleAlpha(alpha)
		tempImg.DrawImage(shadowBlock, opt)
	}

	// 步驟2：垂直模糊 (將水平模糊的結果再進行垂直模糊)
	for j := -count; j <= count; j++ {
		if j == 0 {
			continue // 中心點已經在水平模糊中處理過
		}
		opt := &ebiten.DrawImageOptions{}
		opt.Filter = ebiten.FilterLinear
		opt.Blend = blend
		// opt.GeoM.Scale(0.9, 0.9)
		opt.GeoM.Translate(0.0, float64(j)+shadowLength)
		opt.ColorScale.ScaleAlpha(alpha)
		shadowImg.DrawImage(tempImg, opt)
	}

	c._shadowCache.Update(shadowImg)
	screen.DrawImage(shadowImg, finalOpt)

}
