package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/pkg/logs"
)

func EbitenUpdate(sv SomeView) {
	root(sv).calculateStage()
	tickTock()
}

func EbitenDraw(screen *ebiten.Image, view SomeView) {
	println()

	if p := view.params(); p.w <= 0 || p.h <= 0 {
		logs.Warnf("view is not ready yet: size(%d, %d)", p.w, p.h)
		return
	}

	view.draw(screen)
}

type View interface {
	Body() SomeView
}

//go:generate domaingen -destination=view.option.go -package=ebui -name=viewOption
type SomeView interface {
	View

	draw(screen *ebiten.Image)
	params() *view
	initBounds() (int, int)

	Frame(w, h int) SomeView
	ForegroundColor(clr color.Color) SomeView
	BackgroundColor(clr color.Color) SomeView
	Padding(top, right, bottom, left int) SomeView
	FontSize(size font.Size) SomeView
	FontWeight(weight font.Weight) SomeView
	CornerRadius(radius ...int) SomeView
	LineSpacing(spacing float64) SomeView
	Kern(kern int) SomeView
	Italic() SomeView
}
