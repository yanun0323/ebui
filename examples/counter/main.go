package main

import (
	"log"

	. "github.com/yanun0323/ebui"
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
		Text("Hello").
			Padding(Bind(NewInset(10, 10, 10, 10))).
			BackgroundColor(Bind(NewColor(125, 125, 255, 0))).
			Padding(Bind(NewInset(15, 15, 15, 15))).
			BackgroundColor(Bind(NewColor(255, 125, 125, 0))),
		// Rectangle().BackgroundColor(Bind(NewColor(125, 125, 255, 0))),
		// Rectangle().BackgroundColor(Bind(NewColor(255, 125, 125, 0))),
		// VStack(
		// 	Rectangle().BackgroundColor(Bind(NewColor(0, 0, 125, 0))),
		// 	HStack(
		// 		Rectangle().BackgroundColor(Bind(NewColor(125, 0, 0, 0))),
		// 		Rectangle().BackgroundColor(Bind(NewColor(0, 125, 0, 0))),
		// 	),
		// ),
		// Image(Bind("image.jpg")).
		// 	ScaleToFit().
		// 	KeepAspectRatio(),
		// Text("Counter Example").
		// 	FontSize(Bind(font.Body)).
		// 	FontAlignment(Bind(font.AlignCenter)).
		// 	FontLineHeight(Bind(2.0)).
		// 	BackgroundColor(Bind(NewColor(125, 125, 255, 0))),
		// Text(v.countText).
		// 	FontSize(Bind(font.Body)).
		// 	FontAlignment(Bind(font.AlignCenter)).
		// 	FontLineHeight(Bind(2.0)),
		// ZStack(
		// 	VStack(
		// 		Spacer(),
		// 		HStack(
		// 			Spacer(),
		// 			Text("TEXT").
		// 				FontSize(Bind(font.Title3)).
		// 				BackgroundColor(Bind(NewColor(125, 64, 64, 125))),
		// 			Spacer(),
		// 		),
		// 		Spacer(),
		// 	).
		// 		Padding(Bind(NewInset(20, 20, 20, 20))).
		// 		BackgroundColor(Bind(NewColor(200, 200, 200, 0))),
		// ).
		// 	BackgroundColor(Bind(NewColor(100, 100, 100, 0))),
		// HStack(
		// 	Button(func() {
		// 		v.count.Set(v.count.Get() + 1)
		// 		v.countText.Set(fmt.Sprintf("Current Value: %d", v.count.Get()))
		// 		println(fmt.Sprintf("set text to: %s", v.countText.Get()))
		// 	}, func() SomeView {
		// 		return Text("Increase").
		// 			Padding(Bind(NewInset(5, 5, 5, 5)))
		// 	}).
		// 		Padding(Bind(NewInset(15, 15, 15, 15))).
		// 		RoundCorner(Bind(10.0)).
		// 		BackgroundColor(Bind(NewColor(100, 100, 100, 0))),
		// 	Spacer(),
		// 	Button(func() {
		// 		v.count.Set(v.count.Get() - 1)
		// 		v.countText.Set(fmt.Sprintf("Current Value: %d", v.count.Get()))
		// 		println(fmt.Sprintf("set text to: %s", v.countText.Get()))
		// 	}, func() SomeView {
		// 		return Text("Decrease").
		// 			BackgroundColor(Bind(NewColor(200, 200, 200, 0)))
		// 	}),
		// ),
	)
	// Padding(Bind(NewInset(15, 15, 15, 15))).
	// BackgroundColor(Bind(NewColor(255, 0, 0, 0)))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(100, 0, 0, 0))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
