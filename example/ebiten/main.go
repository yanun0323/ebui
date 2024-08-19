package main

import (
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui"
)

func main() {
	g := &Game{
		t: "Hello,",
	}

	ebiten.SetWindowSize(1366, 768)
	ebiten.SetWindowTitle("ebui")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.contentView = ebui.HStack(
		ebui.Spacer(),
		ebui.Button(
			func() {
				println("Clicked")
			},
			ebui.Text("Click Me").
				ForegroundColor(color.White),
		).
			Frame(150, 150).
			BackgroundColor(color.Gray{128}),
		ebui.Text(&g.t).
			ForegroundColor(color.Gray{128}).
			Padding(10, 10, 10, 10).
			BackgroundColor(color.White).
			Padding(30, 30, 30, 30).
			BackgroundColor(color.Gray{64}).
			Frame(400, -1),
		ebui.VStack(
			ebui.Text("VStack").
				BackgroundColor(color.Gray{32}),
			ebui.HStack(
				ebui.Text("HStack 1").
					BackgroundColor(color.Gray{160}),
				ebui.Text("HStack 2").
					Padding(30, 30, 30, 30).
					BackgroundColor(color.Gray{192}),
			),
		),
		ebui.Text("Trailing").
			BackgroundColor(color.Gray{16}),
		ebui.Spacer(),
		// ebui.Text(&g.t).
		// 	ForegroundColor(color.RGBA{0, 255, 0, 255}),
	).
		ForegroundColor(color.RGBA{200, 200, 200, 255}).
		BackgroundColor(color.RGBA{255, 0, 0, 255}).
		// Frame(1000, 768)
		Padding(5, 5, 5, 5)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("run game", "error", err)
	}
}

var _ ebiten.Game = (*Game)(nil)

type Game struct {
	counter     int
	t           string
	contentView ebui.SomeView
}

func (g *Game) Update() error {
	g.counter++

	if g.counter%ebiten.TPS() == 0 {
		g.t += "a"
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{128})
	ebui.EbitenDraw(screen, g.contentView)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
