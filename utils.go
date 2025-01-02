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
                    tmp.Set(x+i, y+j, col)
                }
            }
        }
    }

    // 左上角
    drawCorner(0, 0, math.Pi)
    // 右上角
    drawCorner(w-int(radius), 0, math.Pi*1.5)
    // 左下角
    drawCorner(0, h-int(radius), math.Pi*0.5)
    // 右下角
    drawCorner(w-int(radius), h-int(radius), 0)

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