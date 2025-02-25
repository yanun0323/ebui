package main

import (
	"context"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/font"
)

var (
	red   = NewColor(255, 0, 0, 0)
	blue  = NewColor(0, 0, 255, 0)
	green = NewColor(0, 255, 0, 0)
	black = NewColor(16, 16, 16, 16)
)

func NewContentView() View {
	pauseString := "Random Color Pause"
	startString := "Random Color Start "

	return &ContentView{
		pauseString: pauseString,
		startString: startString,
		isRandom:    &atomic.Bool{},
		color:       Bind(red),
		content:     Bind(pauseString),
	}
}

type ContentView struct {
	pauseString string
	startString string

	isRandom  *atomic.Bool
	pauseFunc context.CancelFunc
	color     *Binding[AnyColor]
	content   *Binding[string]
}

func (v *ContentView) Body() SomeView {
	println("ContentView.Body() called")

	// 原始的 VStack
	return VStack(
		HStack(
			Circle().
				Frame(Bind(NewSize(150, 300))).
				BackgroundColor(Bind(red)),
			EmptyView(),
			Rectangle().
				Frame(Bind(NewSize(200, 150))).
				BackgroundColor(Bind(blue)),
		),
		Rectangle().BackgroundColor(Bind(NewColor(125, 125, 255, 0))).
			Frame(Bind(NewSize(100, 100))).
			BackgroundColor(Bind(red)).
			Padding(Bind(CGInset{10, 10, 10, 10})).
			BackgroundColor(Bind(blue)),
		ZStack(
			Rectangle().
				BackgroundColor(v.color),
			VStack(
				Spacer(),
				HStack(
					Spacer(),
					Text(v.content).
						FontLetterSpacing(Bind(2.0)),
					Spacer(),
				),
				Spacer(),
			),
		).
			Padding(Bind(CGInset{10, 10, 10, 10})),
		HStack(
			Button(func() {
				if v.isRandom.Swap(true) {
					return
				}

				v.content.Set(v.startString)
				println(v.startString)
				go func() {
					ctx, cancel := context.WithCancel(context.Background())
					v.pauseFunc = cancel
					ticker := time.NewTicker(200 * time.Millisecond)
					for {
						select {
						case <-ticker.C:
							v.color.Set(NewColor(rand.Intn(255), rand.Intn(255), rand.Intn(255), 255))
						case <-ctx.Done():
							return
						}
					}
				}()
			}, func() SomeView {
				return Text("Start Random Color").
					Padding(Bind(CGInset{15, 15, 15, 15})).
					BackgroundColor(Bind(black)).
					RoundCorner(Bind(15.0))
			}),
			Spacer(),
			Button(func() {
				if v.isRandom.Swap(false) {
					if fn := v.pauseFunc; fn != nil {
						fn()
					}

					v.content.Set(v.pauseString)
					println(v.pauseString)
				}
			}, func() SomeView {
				return Text("Pause Random Color").
					Padding(Bind(CGInset{15, 15, 15, 15})).
					BackgroundColor(Bind(black)).
					RoundCorner(Bind(15.0))
			}),
		).
			FontSize(Bind(font.Caption)),
	).Padding(Bind(CGInset{15, 15, 15, 15}))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(32, 32, 32, 32))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(false)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
