package main

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/example/component"
)

func main() {
	ebiten.SetWindowSize(1366, 768)
	ebiten.SetWindowTitle("ebui")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	contentView := component.ContentView("title", "content")
	g := NewGame(contentView)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("run game", "error", err)
	}
}

func NewGame(contentView ebui.View) *Game {
	return &Game{
		contentView: contentView.Body(),
	}
}

type Game struct {
	contentView ebui.SomeView
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebui.EbitenDraw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
