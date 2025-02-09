package main

import (
	"fmt"
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
		count:     NewBinding(0),
		countText: NewBinding("當前數值: "),
	}
}

func (v *ContentView) Body() SomeView {
	return VStack(
		Text("計數器示例").
			FontSize(NewBinding(font.Body)).
			FontAlignment(NewBinding(font.AlignCenter)).
			FontLineHeight(NewBinding(2.0)),
		Text(v.countText).
			FontSize(NewBinding(font.Body)).
			FontAlignment(NewBinding(font.AlignCenter)).
			FontLineHeight(NewBinding(2.0)),
		VStack(
			Button(func() {
				v.count.Set(v.count.Get() + 1)
				v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, func() SomeView {
				return Text("增加")
				// BackgroundColor(NewBinding[color.Color](color.Gray{200}))
			}),
			// Button(func() {
			// 	v.count.Set(v.count.Get() - 1)
			// 	v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
			// 	println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			// }, func() SomeView {
			// 	return Text("減少").
			// 		BackgroundColor(NewBinding[color.Color](color.Gray{200}))
			// }),
			// Spacer(),
		),
	).Padding(NewBinding(15.0)).BackgroundColor(NewBinding[color.Color](color.Gray{30}))
}

func main() {
	ebiten.SetWindowSize(400, 300)
	ebiten.SetWindowTitle("EBUI Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := NewApplication(NewContentView())

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
