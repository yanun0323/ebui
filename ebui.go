package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

func EbitenDraw(screen *ebiten.Image, view SomeView) {
	logs.Debug("")
	w, h := ebiten.WindowSize()
	opt := newViewOption(view).CreateChild(view, 0, 0, 0, 0, w, h, image.Rect(0, 0, w, h))
	view.draw(screen, opt)
}
