package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

func EbitenUpdate(sv SomeView) {
	root(sv).calculateStage()
	tickTock()
}

func EbitenDraw(screen *ebiten.Image, view SomeView) {
	println()

	if p := view.params(); p.size.w <= 0 || p.size.h <= 0 {
		logs.Warnf("view is not ready yet: size(%d, %d)", p.size.w, p.size.h)
		return
	}

	view.draw(screen)
}
