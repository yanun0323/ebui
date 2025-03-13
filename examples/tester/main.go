package main

import (
	"log"
	"sync/atomic"
	"time"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/layout"
)

var (
	white  = NewColor(255)
	red    = NewColor(255, 0, 0)
	blue   = NewColor(0, 0, 255)
	green  = NewColor(0, 255, 0)
	yellow = NewColor(255, 255, 0)

	borderInset = NewInset(1)
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
		offset:      Bind(CGPoint{}).Animated(animation.EaseInOut(time.Second)),
		offset2:     Bind(CGPoint{}).Animated(animation.EaseInOut(time.Second)),
	}
}

type ContentView struct {
	pauseString string
	startString string

	isRandom *atomic.Bool
	color    *Binding[CGColor]
	content  *Binding[string]
	offset   *Binding[CGPoint]
	offset2  *Binding[CGPoint]
}

func (v *ContentView) Body() SomeView {
	println("ContentView.Body() called")

	return VStack(
		HStack(
			Divider(),

			VStack(
				Text("align leading"),
				Rectangle().
					Frame(Bind(NewSize(80, 40))).
					RoundCorner(Bind(15.0)).
					BackgroundColor(Bind(red)).
					Border(Bind(borderInset), Bind(white)).
					Padding(Bind(NewInset(5))),
				Rectangle().
					Frame(Bind(NewSize(100, 30))).
					Border(Bind(borderInset), Bind(white)).
					RoundCorner(Bind(15.0)).
					BackgroundColor(Bind(green)).
					Padding(Bind(NewInset(5))),
			).Modify(debugFunc).
				Align(Bind(layout.AlignLeading)),

			Divider(),

			Spacer().Padding(Bind(NewInset(5))).Debug(),

			VStack(
				Text("align center"),
				Rectangle().
					Frame(Bind(NewSize(Inf, 20))).
					BackgroundColor(Bind(blue)).
					RoundCorner(Bind(15.0)).
					Debug().
					Border(Bind(borderInset), Bind(white)).
					Padding(Bind(NewInset(5))).
					Debug(),
				Rectangle().
					Frame(Bind(NewSize(80, 40))).
					Border(Bind(borderInset), Bind(white)).
					Padding(Bind(NewInset(5))),
			).Modify(debugFunc).
				Align(Bind(layout.AlignCenter)),

			Divider(),

			VStack(
				Text("Counter Demo"),
				Text("align trailing"),
				Rectangle().
					Frame(Bind(NewSize(120, 30))).
					Border(Bind(borderInset), Bind(white)).
					BackgroundColor(Bind(yellow)).
					RoundCorner(Bind(15.0)).
					Padding(Bind(NewInset(5))),
				Rectangle().
					Frame(Bind(NewSize(100, 20))).
					Border(Bind(CGInset{}), Bind(white)).
					BackgroundColor(Bind(green)).
					RoundCorner(Bind(15.0)).
					Padding(Bind(NewInset(5))),
			).Modify(debugFunc).
				Align(Bind(layout.AlignTrailing)),
		).Debug(),

		Divider().Padding(Bind(NewInset(5))).Debug(),

		HStack(
			VStack(
				HStack(
					Text("align top"),
					rect(50, 20, red),
					rect(50, 40, green),
					rect(50, 60, blue),
				).Modify(debugFunc).
					Align(Bind(layout.AlignTop)),

				Divider(),

				HStack(
					Text("align center").Debug(),
					rect(50, 20, red),
					rect(50, 40, green),
					rect(50, 60, blue),
				).Modify(debugFunc).
					Align(Bind(layout.AlignCenter)),

				Divider(),

				HStack(
					Text("align bottom").FontKerning(Bind(1.0)).Debug(),
					rect(50, 20, red),
					rect(50, 40, green),
					rect(50, 60, blue),
				).Modify(debugFunc).
					Align(Bind(layout.AlignBottom)),

				Divider(),
			).Debug(),

			VStack(
				Text("測試開關"),

				Toggle(Bind(false)),
			).Debug(),
		),

		Divider().Padding(Bind(NewInset(5))).Debug(),

		VStack(
			HStack(
				Text("EaseInOut").
					Frame(Bind(NewSize(150, Inf))),
				Button("", func() {
					if v.offset.Get().X == 0 {
						v.offset.Set(CGPoint{X: 280, Y: 0})
						v.offset2.Set(CGPoint{X: 280, Y: 0})
					} else {
						v.offset.Set(CGPoint{X: 0, Y: 0})
						v.offset2.Set(CGPoint{X: 0, Y: 0})
					}
				}, func() SomeView {
					return HStack(
						Circle().
							Frame(Bind(NewSize(20))).
							Offset(v.offset).
							BackgroundColor(Bind(red)),
						Spacer(),
					).Frame(Bind(NewSize(300, 20))).
						Padding(Bind(NewInset(5))).
						BackgroundColor(Bind(blue)).
						RoundCorner(Bind(10.0))
				}).Padding(Bind(NewInset(5))),
			).Padding(Bind(NewInset(0, 0, 0, 10))).
				Align(Bind(layout.AlignCenter)),
			HStack(
				Text("EaseInOut Strength").
					Frame(Bind(NewSize(150, Inf))),
				HStack(
					Circle().
						Frame(Bind(NewSize(20))).
						Offset(v.offset2.Animated(animation.EaseInOut(time.Second).Strengthen(5))).
						BackgroundColor(Bind(red)),
					Spacer(),
				).Frame(Bind(NewSize(300, 20))).
					Padding(Bind(NewInset(5))).
					BackgroundColor(Bind(blue)).
					RoundCorner(Bind(10.0)).
					Padding(Bind(NewInset(5))),
			).Padding(Bind(NewInset(0, 0, 0, 10))).
				Align(Bind(layout.AlignCenter)),
		).Debug(),
	).Debug().FontSize(Const(font.Body)).
		Padding(Bind(NewInset(30))).
		Debug().
		FontWeight(Bind(font.Bold))
}

func rect(w, h int, color CGColor) SomeView {
	sli := []int{20, 40, 60, 80, 100}
	sz := Bind(NewSize(w, h))
	return Button("rect", func() {
		WithAnimation(func() {
			h := int(sz.Get().Height)
			idx := 0
			for i, s := range sli {
				if h == s {
					idx = i
					break
				}
			}
			idx = (idx + 1) % len(sli)
			sz.Set(NewSize(w, sli[idx]))
		})
	}, func() SomeView {
		return Rectangle().
			Frame(sz).
			BackgroundColor(Bind(color)).
			RoundCorner(Bind(15.0)).
			Padding(Bind(NewInset(5)))
	})
}

func debugFunc(view SomeView) SomeView {
	return view.
		Padding(Bind(NewInset(5))).
		Border(Bind(NewInset(1)), Bind(red)).
		Padding(Bind(NewInset(5)))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(64, 16, 64, 64))
	app.SetWindowSize(800, 600)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.VSyncEnabled(false)
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}

}
