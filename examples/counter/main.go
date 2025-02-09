package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/yanun0323/ebui"
)

type ContentView struct {
	count     *Binding[int]
	countText *Binding[string]
}

func NewContentView() *ContentView {
	return &ContentView{
		count:     NewBinding(0),
		countText: NewBinding(""),
	}
}

func (v *ContentView) Build() View {
	return v.Body().Build()
}

func (v *ContentView) Body() SomeView {
	return VStack(
		TextStatic("計數器示例").WithStyle(TextStyle{
			Size:       12,
			Color:      color.Black,
			Alignment:  TextAlignCenter,
			LineHeight: 2.0,
		}),
		Text(v.countText).WithStyle(TextStyle{
			Size:       12,
			Color:      color.Black,
			Alignment:  TextAlignCenter,
			LineHeight: 2.0,
		}),
		VStack(
			Button(func() {
				v.count.Set(v.count.Get() + 1)
				v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, TextStatic("增加")),
			Button(func() {
				v.count.Set(v.count.Get() - 1)
				v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, TextStatic("減少")),
			Spacer().WithSize(8),
		).WithPadding(8),
	).WithPadding(15).BackgroundColor(color.White)
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
