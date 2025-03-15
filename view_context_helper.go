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
