package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Spacer 元件
type spacerImpl struct {
	*viewContext
	frame image.Rectangle
}

func Spacer() SomeView {
	return &spacerImpl{}
}

func (s *spacerImpl) Body() SomeView {
	return s
}

func (s *spacerImpl) layout(bounds image.Rectangle) image.Rectangle {
	s.frame = bounds
	return bounds
}

func (s *spacerImpl) draw(screen *ebiten.Image) {
	// Spacer 是空白元件，不需要繪製任何內容
}
