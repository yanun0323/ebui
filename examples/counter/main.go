package main

import (
	"fmt"
	"log"

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
		Rectangle().BackgroundColor(Bind(NewColor(125, 125, 255, 255))),
		Rectangle().BackgroundColor(Bind(NewColor(255, 125, 125, 255))),
		VStack(
			Rectangle().BackgroundColor(Bind(NewColor(0, 0, 125, 255))),
			HStack(
				Rectangle().BackgroundColor(Bind(NewColor(125, 0, 0, 255))),
				Rectangle().BackgroundColor(Bind(NewColor(0, 125, 0, 255))),
			),
		),
		Image(Bind("image.jpg")).
			ScaleToFit().
			KeepAspectRatio(),
		Text(Bind("Counter Example")).
			FontSize(Bind(font.Body)).
			FontAlignment(Bind(font.AlignCenter)).
			FontLineHeight(Bind(2.0)).
			BackgroundColor(Bind(NewColor(125, 125, 255, 255))),
		Text(v.countText).
			FontSize(Bind(font.Body)).
			FontAlignment(Bind(font.AlignCenter)).
			FontLineHeight(Bind(2.0)),
		ZStack(
			VStack(
				Spacer(),
				HStack(
					Spacer(),
					Text(Bind("TEXT")).
						FontSize(Bind(font.Title3)).
						BackgroundColor(Bind(NewColor(125, 64, 64, 255))),
					Spacer(),
				),
				Spacer(),
			).
				Padding(Bind(NewInset(20, 20, 20, 20))).
				BackgroundColor(Bind(NewColor(160, 160, 160, 255))).
				Padding(Bind(NewInset(20, 20, 20, 20))).
				BackgroundColor(Bind(NewColor(200, 200, 200, 255))),
		).
			BackgroundColor(Bind(NewColor(100, 100, 100, 255))),
		HStack(
			Button(func() {
				v.count.Set(v.count.Get() + 1)
				v.countText.Set(fmt.Sprintf("Current Value: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, func() SomeView {
				return Text(Bind("Increase")).
					Padding(Bind(NewInset(5, 5, 5, 5))).
					BackgroundColor(Bind(NewColor(200, 200, 200, 255))).
					RoundCorner(Bind(10.0))
			}).
				Padding(Bind(NewInset(15, 15, 15, 15))),
			Spacer(),
			Button(func() {
				v.count.Set(v.count.Get() - 1)
				v.countText.Set(fmt.Sprintf("Current Value: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, func() SomeView {
				return Text(Bind("Decrease")).
					Padding(Bind(NewInset(5, 5, 5, 5))).
					BackgroundColor(Bind(NewColor(200, 200, 200, 255))).
					RoundCorner(Bind(20.0))
			}).
				Padding(Bind(NewInset(15, 15, 15, 15))),
		),
	).
		Padding(Bind(NewInset(15, 15, 15, 15)))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(32, 32, 32, 32))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
