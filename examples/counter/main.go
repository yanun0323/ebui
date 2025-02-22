package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/font"
)

type ContentView struct {
	count     *Binding[int]
	countText *Binding[string]
}

func NewContentView() View {
	return &ContentView{
		count:     Bind(0),
		countText: Bind("Current Value: "),
	}
}

func (v *ContentView) Body() SomeView {
	return VStack(
		Text("Counter Example").
			FontSize(Bind(font.Body)).
			FontAlignment(Bind(font.AlignCenter)).
			FontLineHeight(Bind(2.0)).
			BackgroundColor(Bind[color.Color](color.RGBA{125, 125, 255, 0})),
		Text(v.countText).
			FontSize(Bind(font.Body)).
			FontAlignment(Bind(font.AlignCenter)).
			FontLineHeight(Bind(2.0)),
		ZStack(
			VStack(
				Spacer(),
				Text("TEST TEXT").
					FontSize(Bind(font.Title3)),
			).
				// Padding(NewBinding(20.0)).
				BackgroundColor(Bind[color.Color](color.Gray{200})),
		).
			BackgroundColor(Bind[color.Color](color.Gray{100})),
		// VStack(
		// 	Button(func() {
		// 		v.count.Set(v.count.Get() + 1)
		// 		v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
		// 		println(fmt.Sprintf("set text to: %s", v.countText.Get()))
		// 	}, func() SomeView {
		// 		return Text("增加")
		// 		// BackgroundColor(NewBinding[color.Color](color.Gray{200}))
		// 	}),
		// 	Button(func() {
		// 		v.count.Set(v.count.Get() - 1)
		// 		v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
		// 		println(fmt.Sprintf("set text to: %s", v.countText.Get()))
		// 	}, func() SomeView {
		// 		return Text("減少").
		// 			BackgroundColor(NewBinding[color.Color](color.Gray{200}))
		// 	}),
		// 	Spacer(),
		// ),
	).
		// BackgroundColor(NewBinding[color.Color](color.RGBA{255, 0, 0, 0})).
		Padding(Bind(15.0))
}

func main() {
	ebiten.SetWindowSize(400, 300)
	ebiten.SetWindowTitle("EBUI Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := NewApplication(NewContentView())
	app.SetBackgroundColor(color.RGBA{100, 0, 0, 0})

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
