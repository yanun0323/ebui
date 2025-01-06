package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui"
)

type ContentView struct {
	count     *ebui.Binding[int]
	countText *ebui.Binding[string]
}

func NewContentView() *ContentView {
	return &ContentView{
		count:     ebui.NewBinding(0),
		countText: ebui.NewBinding(""),
	}
}

func (v *ContentView) Build() ebui.View {
	return v.Body().Build()
}

func (v *ContentView) Body() ebui.SomeView {
	return ebui.VStack(
		ebui.TextStatic("計數器示例").WithStyle(ebui.TextStyle{
			Size:       12,
			Color:      color.Black,
			Alignment:  ebui.TextAlignCenter,
			LineHeight: 2.0,
		}),
		ebui.Text(v.countText).WithStyle(ebui.TextStyle{
			Size:       12,
			Color:      color.Black,
			Alignment:  ebui.TextAlignCenter,
			LineHeight: 2.0,
		}),
		ebui.VStack(
			ebui.Button(func() {
				v.count.Set(v.count.Get() + 1)
				v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, ebui.TextStatic("增加")),
			ebui.Button(func() {
				v.count.Set(v.count.Get() - 1)
				v.countText.Set(fmt.Sprintf("當前數值: %d", v.count.Get()))
				println(fmt.Sprintf("set text to: %s", v.countText.Get()))
			}, ebui.TextStatic("減少")),
			ebui.Spacer().WithSize(8),
		).WithPadding(8),
	).WithPadding(15).BackgroundColor(color.White)
}

func main() {
	ebiten.SetWindowSize(400, 300)
	ebiten.SetWindowTitle("EBUI Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := ebui.NewApplication(NewContentView())

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
