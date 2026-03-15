package examples

import (
	. "github.com/yanun0323/ebui"
)

func NewTesterContentView() View {
	return &TesterContentView{
		lightMode: Bind(true).Animated(),
		length:    Bind(5.0).Animated(),
	}
}

type TesterContentView struct {
	lightMode *Binding[bool]
	length    *Binding[float64]
}

func (v *TesterContentView) Body() SomeView {
	return VStack(
		Rectangle().Fill(Const(green)),
		Text("Hello, World!"),
		Rectangle().Fill(Const(red)),
	).Spacing().Padding().Debug()
}

func (v *TesterContentView) CurrentBackgroundColorText() SomeView {
	return Text(BindOneWay(v.lightMode, func(lightMode bool) string {
		if lightMode {
			return "White"
		}
		return "Black"
	}))
}

func (v *TesterContentView) ChangeBackgroundColor() {
	v.lightMode.Set(!v.lightMode.Get())
}

func Preview_Tester() View {
	return NewTesterContentView()
}
