package examples

import (
	"fmt"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/layout"
)

var (
	gray   = NewColor(128)
	white  = NewColor(239)
	black  = NewColor(16)
	red    = NewColor(128, 0, 0)
	green  = NewColor(0, 128, 0)
	blue   = NewColor(0, 0, 128)
	yellow = NewColor(128, 128, 0)
)

func NewContentView() View {
	return &ContentView{
		lightMode: Bind(true).Animated(),
		length:    Bind(5.0).Animated(),
	}
}

type ContentView struct {
	lightMode *Binding[bool]
	length    *Binding[float64]
}

func (v *ContentView) Body() SomeView {
	return ScrollView(
		VStack(
			HStack(
				Text("Light Mode"),
				Toggle(v.lightMode),
			),

			Text(BindOneWay(v.length, func(length float64) string {
				return fmt.Sprintf("length: %.2f", length)
			})),

			Slider(v.length, Const(0.0), Const(100.0)),

			HStack(
				Rectangle().Fill(Const(red)).Frame(Const(NewSize(100, 100))),
				Rectangle().Fill(Const(yellow)).Frame(Const(NewSize(100, 100))),
			).RoundCorner(),

			HStack(
				Circle().Fill(Const(red)).Frame(Const(NewSize(100))).Shadow(v.length),
				Rectangle().Fill(Const(yellow)).Frame(Const(NewSize(150))).Shadow(v.length).RoundCorner(),
				Rectangle().Fill(Const(green)).Frame(Const(NewSize(150))).Border(Const(NewInset(10)), Const(blue)).RoundCorner().Shadow(v.length),
				Rectangle().Fill(Const(blue)).Frame(Const(NewSize(150))).Border(Const(NewInset(10)), Const(green)).RoundCorner(Const(150.0)).Shadow(v.length),
			).Spacing(Const(30.0)).Align(Const(layout.AlignCenter)),

			Rectangle().Fill(Const(red)).Frame(Const(NewSize(100, 1000))),
		).Spacing(Const(30.0)).
			Center().
			Align(Const(layout.AlignCenter)).
			ForegroundColor(BindOneWay(v.lightMode, func(lightMode bool) CGColor {
				if lightMode {
					return black
				}
				return white
			})).
			BackgroundColor(BindOneWay(v.lightMode, func(lightMode bool) CGColor {
				if lightMode {
					return white
				}
				return black
			})),
	)
}

func (v *ContentView) CurrentBackgroundColorText() SomeView {
	return Text(BindOneWay(v.lightMode, func(lightMode bool) string {
		if lightMode {
			return "White"
		}
		return "Black"
	}))
}

func (v *ContentView) ChangeBackgroundColor() {
	v.lightMode.Set(!v.lightMode.Get())
}

func Preview_ScrollView() View {
	return NewContentView()
}
