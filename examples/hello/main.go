package main

import (
	"log"

	ui "github.com/yanun0323/ebui"
	ebitenbackend "github.com/yanun0323/ebui/backend/ebiten"
)

type helloScreen struct{}

func (helloScreen) Body(ctx ui.ViewContext) ui.Element {
	titleStyle := ui.DefaultTextStyle()
	titleStyle.Size = 28
	titleStyle.Weight = ui.WeightBold
	titleStyle.LineHeight = 6

	return ui.Apply(
		ui.ZStack(
			ui.Apply(
				ui.Rectangle(),
				ui.Fill(ui.RGBA(30, 34, 45, 255)),
				ui.Frame(ui.S(640, 360)),
			),
			ui.Apply(
				ui.Text("哈囉，EBUI v2"),
				ui.Foreground(ui.White),
				ui.Font(titleStyle),
				ui.Padding(ui.Symmetric(12, 18)),
				ui.Background(ui.RGBA(64, 124, 255, 255)),
			),
		),
		ui.Frame(ui.S(640, 360)),
		ui.Align(ui.AlignCenter),
	)
}

func main() {
	app, err := ebitenbackend.NewApp(
		helloScreen{},
		ebitenbackend.WithWindowTitle("EBUI v2 Hello"),
		ebitenbackend.WithWindowSize(640, 360),
		ebitenbackend.WithBackground(ui.RGBA(14, 18, 26, 255)),
		ebitenbackend.WithFontPath("resource/Taipei Sans TC Regular.ttf"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
