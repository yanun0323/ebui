package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Spacer 元件
type spacerImpl struct {
	size  float64
	frame image.Rectangle
}

func Spacer() ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &spacerImpl{size: 1.0}
		},
	}
}

// 實現 View 介面
func (s *spacerImpl) Build() View {
	return s
}

func (s *spacerImpl) Layout(bounds image.Rectangle) image.Rectangle {
	s.frame = bounds
	return bounds
}

func (s *spacerImpl) Draw(screen *ebiten.Image) {
	// Spacer 是空白元件，不需要繪製任何內容
}

func (v ViewBuilder) WithSize(size float64) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			spacer := v.Build().(*spacerImpl)
			spacer.size = size
			return spacer
		},
	}
}
