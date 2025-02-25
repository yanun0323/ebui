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
func drawRoundedRect(screen *ebiten.Image, bounds CGRect, radius float64, col color.Color, op *ebiten.DrawImageOptions) {
	scale := 3.0
	w := int(bounds.Dx() * scale)
	h := int(bounds.Dy() * scale)
	r := int(radius * scale)

	img := ebiten.NewImage(w, h)
	img.Fill(col)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if (x >= r && x <= w-r) || (y >= r && y <= h-r) {
				continue
			}

			// 左上角
			if x < r && y < r && (x-r)*(x-r)+(y-r)*(y-r) > r*r {
				img.Set(x, y, color.Transparent)
			}

			// 右上角
			if x >= w-r && y < r && (x-(w-r))*(x-(w-r))+(y-r)*(y-r) > r*r {
				img.Set(x, y, color.Transparent)
			}

			// 左下角
			if x < r && y >= h-r && (x-r)*(x-r)+(y-(h-r))*(y-(h-r)) > r*r {
				img.Set(x, y, color.Transparent)
			}

			// 右下角
			if x >= w-r && y >= h-r && (x-(w-r))*(x-(w-r))+(y-(h-r))*(y-(h-r)) > r*r {
				img.Set(x, y, color.Transparent)
			}
		}
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(1/scale, 1/scale)
	opt.Filter = ebiten.FilterLinear
	opt.GeoM.Concat(op.GeoM)
	opt.ColorScale.ScaleWithColorScale(op.ColorScale)

	screen.DrawImage(img, opt)
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
