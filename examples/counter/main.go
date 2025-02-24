package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/yanun0323/ebui"
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
		countText: Bind("Current Value: 0"),
	}
}

func (v *ContentView) Body() SomeView {
	return VStack(
		Rectangle().BackgroundColor(Bind[color.Color](color.RGBA{125, 125, 255, 0})),
		Rectangle().BackgroundColor(Bind[color.Color](color.RGBA{255, 125, 125, 0})),
		VStack(
			Rectangle().BackgroundColor(Bind[color.Color](color.RGBA{0, 0, 125, 0})),
			HStack(
				Rectangle().BackgroundColor(Bind[color.Color](color.RGBA{125, 0, 0, 0})),
				Rectangle().BackgroundColor(Bind[color.Color](color.RGBA{0, 125, 0, 0})),
			),
		),
		Image(Bind("image.jpg")).
			ScaleToFit().
			KeepAspectRatio(),
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
				HStack(
					Spacer(),
					Text("TEXT").
						FontSize(Bind(font.Title3)).
						BackgroundColor(Bind[color.Color](color.RGBA{125, 64, 64, 125})),
					Spacer(),
				),
				Spacer(),
			).
				Padding(Bind(Inset{20, 20, 20, 20})).
				BackgroundColor(Bind[color.Color](color.Gray{200})),
		).
			BackgroundColor(Bind[color.Color](color.Gray{100})),
		HStack(
			Button(func() {
				v.count.Set(v.count.Get() + 1)
				v.countText.Set(fmt.Sprintf("Current Value: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, func() SomeView {
				return Text("Increase").
					Padding(Bind(Inset{5, 5, 5, 5}))
			}).
				Padding(Bind(Inset{15, 15, 15, 15})).
				RoundCorner(Bind(10.0)).
				BackgroundColor(Bind[color.Color](color.Gray{100})),
			Spacer(),
			Button(func() {
				v.count.Set(v.count.Get() - 1)
				v.countText.Set(fmt.Sprintf("Current Value: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, func() SomeView {
				return Text("Decrease").
					BackgroundColor(Bind[color.Color](color.Gray{200}))
			}),
		),
	).
		Padding(Bind(Inset{15, 15, 15, 15})).
		BackgroundColor(Bind[color.Color](color.RGBA{255, 0, 0, 0}))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(color.RGBA{100, 0, 0, 0})
	app.SetWindowTitle("EBUI Demo")
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	app.SetResourceFolder("resource")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
