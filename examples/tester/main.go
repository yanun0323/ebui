package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct{}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255})

	vector.DrawFilledCircle(screen, 100, 100, 100, color.RGBA{255, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, 300, 300, 100, 100, color.RGBA{0, 255, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten 文字輸入框示例")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
