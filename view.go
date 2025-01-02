package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type View interface {
	SomeView

	Layout(bounds image.Rectangle) image.Rectangle
	Draw(screen *ebiten.Image)
}

// SomeView 是所有 View 的基礎介面
type SomeView interface {
	Build() View
}

// ViewBuilder 用來構建具體的 View
type ViewBuilder struct {
	build func() View
}

func (v ViewBuilder) Build() View {
	return v.build()
}
