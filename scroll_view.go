package ebui

import (
	"image"
	"image/color"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

type ScrollView struct {
	content      View
	frame        image.Rectangle
	contentSize  image.Point
	offset       image.Point
	maxOffset    image.Point
	isDragging   bool
	startDrag    image.Point
	startOffset  image.Point
	cache        *ViewCache
}

func NewScrollView(content SomeView) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &ScrollView{
				content: content.Build(),
				cache:   NewViewCache(),
			}
		},
	}
}

func (sv *ScrollView) Layout(bounds image.Rectangle) image.Rectangle {
	sv.frame = bounds
	
	// 計算內容大小
	contentBounds := sv.content.Layout(image.Rectangle{
		Min: bounds.Min.Sub(sv.offset),
		Max: image.Point{X: bounds.Max.X * 2, Y: bounds.Max.Y * 2},
	})
	
	sv.contentSize = contentBounds.Size()
	sv.maxOffset = image.Point{
		X: max(0, sv.contentSize.X - bounds.Dx()),
		Y: max(0, sv.contentSize.Y - bounds.Dy()),
	}
	
	return bounds
}

func (sv *ScrollView) Draw(screen *ebiten.Image) {
	// 創建裁剪遮罩
	mask := ebiten.NewImage(sv.frame.Dx(), sv.frame.Dy())
	mask.Fill(color.White)
	
	// 設置裁剪區域
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sv.frame.Min.X), float64(sv.frame.Min.Y))
	screen.DrawImage(mask, op)
	
	// 繪製內容
	contentOp := &ebiten.DrawImageOptions{}
	contentOp.GeoM.Translate(float64(-sv.offset.X), float64(-sv.offset.Y))
	sv.content.Draw(screen)
}

func (sv *ScrollView) HandleTouchEvent(event TouchEvent) bool {
	switch event.Phase {
	case TouchPhaseBegan:
		sv.isDragging = true
		sv.startDrag = event.Position
		sv.startOffset = sv.offset
		return true
		
	case TouchPhaseMoved:
		if sv.isDragging {
			delta := event.Position.Sub(sv.startDrag)
			newOffset := sv.startOffset.Sub(delta)
			
			// 限制滾動範圍
			sv.offset = image.Point{
				X: clamp(newOffset.X, 0, sv.maxOffset.X),
				Y: clamp(newOffset.Y, 0, sv.maxOffset.Y),
			}
			return true
		}
		
	case TouchPhaseEnded, TouchPhaseCancelled:
		sv.isDragging = false
		return true
	}
	
	return false
} 