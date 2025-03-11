package ebui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func logf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func someViews(views ...View) []SomeView {
	someViews := make([]SomeView, 0, len(views))
	for _, view := range views {
		someViews = append(someViews, view.Body())
	}
	return someViews
}

// 繪製圓角矩形色塊
func drawRoundedAndBorderRect(screen *ebiten.Image, bounds CGRect, radius float64, bgColor CGColor, border CGInset, borderColor CGColor, bOpt *ebiten.DrawImageOptions) {
	w := int(bounds.Dx() * _roundedScale)
	h := int(bounds.Dy() * _roundedScale)
	r := (radius + border.Top) * _roundedScale // FIXME: using inset to calculate border
	b := border.Top * _roundedScale

	img := ebiten.NewImage(w, h)
	if bgColor != transparent {
		img.Fill(bgColor)
	}

	cornerHandler := newCornerHandler(w, h, r, b)
	cornerHandler.Execute(func(isOutside, isBorder bool, x, y int) {
		if isOutside {
			img.Set(x, y, color.Transparent)
			return
		}

		if isBorder {
			img.Set(x, y, borderColor)
			return
		}
	})

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(_roundedScaleInverse, _roundedScaleInverse)
	opt.Filter = ebiten.FilterLinear
	opt.GeoM.Concat(bOpt.GeoM)
	opt.ColorScale.ScaleWithColorScale(bOpt.ColorScale)

	screen.DrawImage(img, opt)
}

func drawBorderRect(screen *ebiten.Image, bounds CGRect, bgColor CGColor, border CGInset, borderColor CGColor, bOpt *ebiten.DrawImageOptions) {
	w := int(bounds.Dx())
	h := int(bounds.Dy())
	img := ebiten.NewImage(w, h)
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

	screen.DrawImage(img, bOpt)
}

// // 數值限制函數
// func clamp(v, min, max int) int {
// 	if v < min {
// 		return min
// 	}
// 	if v > max {
// 		return max
// 	}
// 	return v
// }

// // 浮點數值限制函數
// func clampF(v, min, max float64) float64 {
// 	if v < min {
// 		return min
// 	}
// 	if v > max {
// 		return max
// 	}
// 	return v
// }

// func dp(value float64) float64 {
// 	// 這裡可以根據螢幕 DPI 進行調整
// 	return value
// }
