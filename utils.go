package ebui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// 繪製圓角矩形
func drawRoundedRect(screen *ebiten.Image, rect image.Rectangle, radius float64, col color.Color) {
	if radius <= 0 {
		screen.SubImage(rect).(*ebiten.Image).Fill(col)
		return
	}

	// 創建臨時圖像
	w, h := rect.Dx(), rect.Dy()
	tmp := ebiten.NewImage(w, h)

	// 繪製四個角
	drawCorner := func(x, y int, startAngle float64) {
		for i := 0; i < int(radius); i++ {
			for j := 0; j < int(radius); j++ {
				dx, dy := float64(i), float64(j)
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist <= radius {
					// 根據不同角落調整像素位置
					switch startAngle {
					case math.Pi: // 左上
						tmp.Set(x+(int(radius)-i), y+(int(radius)-j), col)
					case math.Pi * 1.5: // 右上
						tmp.Set(x+i, y+(int(radius)-j), col)
					case math.Pi * 0.5: // 左下
						tmp.Set(x+(int(radius)-i), y+j, col)
					case 0: // 右下
						tmp.Set(x+i, y+j, col)
					}
				}
			}
		}
	}

	// 繪製四個角
	drawCorner(0, 0, math.Pi)                   // 左上
	drawCorner(w-int(radius), 0, math.Pi*1.5)   // 右上
	drawCorner(0, h-int(radius), math.Pi*0.5)   // 左下
	drawCorner(w-int(radius), h-int(radius), 0) // 右下

	// 填充中間部分
	tmp.SubImage(image.Rect(int(radius), 0, w-int(radius), h)).(*ebiten.Image).Fill(col)
	tmp.SubImage(image.Rect(0, int(radius), w, h-int(radius))).(*ebiten.Image).Fill(col)

	// 繪製到目標
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
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

func pt(value float64) float64 {
	return value
}

func drawRoundedRectBorder(screen *ebiten.Image, rect image.Rectangle, radius, width float64, col color.Color) {
	// 外框
	drawRoundedRect(screen, rect, radius, col)

	// 內框
	inner := image.Rect(
		rect.Min.X+int(width),
		rect.Min.Y+int(width),
		rect.Max.X-int(width),
		rect.Max.Y-int(width),
	)
	drawRoundedRect(screen, inner, radius-width, screen.At(0, 0))
}
