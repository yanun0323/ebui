package examples

import (
	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/font"
)

type ExampleContentView struct {
	Title *Binding[string]
}

func (v *ExampleContentView) Body() SomeView {
	return HStack(
		Spacer(),
		VStack(
			Spacer(),
			Text(v.Title).
				FontSize(Const(font.Title3)),
			Button("Click Me", func() {
				if v.Title.Get() != "Hello, World!" {
					v.Title.Set("Hello, World!")
				} else {
					v.Title.Set("Hello, Ebui!")
				}
			}),
			Spacer(),
		).Spacing(),
		Spacer(),
	).
		BackgroundColor(Const(NewColor(200, 100, 100, 255))).
		Padding(Const(NewInset(10, 10, 10, 10)))
}

func Preview_MyButton() View {
	return &ExampleContentView{
		Title: Bind("Hello, World!"),
	}
}
