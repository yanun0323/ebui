package ebui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func logf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func logfTo(v SomeView, id int64, format string, args ...any) {
	if v.id() == id {
		logf(format, args...)
	}
}

func someViews(views ...View) []SomeView {
	someViews := make([]SomeView, 0, len(views))
	for _, view := range views {
		someViews = append(someViews, view.Body())
	}
	return someViews
}

// 繪製圓角矩形
func drawRoundedRect(screen *ebiten.Image, bounds CGRect, radius float64, col color.Color, op *ebiten.DrawImageOptions) {
	if radius <= 0 {
		return
	}

	// 創建臨時圖像
	w, h := bounds.Dx(), bounds.Dy()
	tmp := ebiten.NewImage(int(w), int(h))

	// 繪製四個角
	drawCorner := func(x, y float64, startAngle float64) {
		const segments = 16
		angleStep := math.Pi / 2 / segments

		for i := 0; i <= segments; i++ {
			angle := startAngle + angleStep*float64(i)
			px := x + radius*math.Cos(angle)
			py := y + radius*math.Sin(angle)
			tmp.Set(int(px), int(py), col)
		}
	}

	// 繪製四個角
	drawCorner(0, 0, math.Pi)            // 左上
	drawCorner(w-radius, 0, math.Pi*1.5) // 右上
	drawCorner(0, h-radius, math.Pi*0.5) // 左下
	drawCorner(w-radius, h-radius, 0)    // 右下

	// 填充中間部分
	tmp.SubImage(bounds.Rect()).(*ebiten.Image).Fill(col)

	screen.DrawImage(tmp, op)
}

// 數值限制函數
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// 浮點數值限制函數
func clampF(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func dp(value float64) float64 {
	// 這裡可以根據螢幕 DPI 進行調整
	return value
}
