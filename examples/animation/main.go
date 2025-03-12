package main

import (
	"fmt"
	"log"
	"time"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/layout"
)

type ContentView struct {
}

func NewContentView() View {
	return &ContentView{}
}

func (v *ContentView) Body() SomeView {
	// 創建一些綁定以進行動畫操作
	count := Bind(0.0)
	opacity := Bind(1.0)
	scale := Bind(1.0)
	enabled := Bind(false)

	// Stack內的所有視圖
	return HStack(
		// 範例1: 基本動畫 - 使用直接設置帶動畫參數的方式
		VStack(
			Text("基本動畫").FontSize(Bind(font.Body)),
			VStack(
				Button("+100(Linear)", func() {
					// 直接在Set時傳入動畫參數
					count.Update(func(val float64) float64 {
						return val + 100
					}, animation.Linear(time.Second*3))
				}).Padding(Bind(NewInset(5))),

				Button("-100(EaseInOut)", func() {
					count.Update(func(val float64) float64 {
						return val - 100
					}, animation.EaseInOut(time.Millisecond*500))
				}).Padding(Bind(NewInset(5))),
			),
			Text(BindFunc(func() string {
				return fmt.Sprintf("計數: %.1f", count.Get())
			}, nil)).FontSize(Bind(font.Body)).Padding(Bind(NewInset(5))),
		).Padding(Bind(NewInset(5))).Debug(),

		// 範例2: 使用WithAnimation上下文 - SwiftUI風格
		VStack(
			Text("上下文動畫").FontSize(Bind(font.Body)),
			Rectangle().
				Frame(Bind(NewSize(100))).
				BackgroundColor(Bind(NewColor(255, 0, 0))).
				Padding(Bind(NewInset(5))).
				Opacity(opacity),
			HStack(
				Button("淡出", func() {
					// 使用動畫上下文包裹多個更改
					WithAnimation(func() {
						opacity.Set(0.2)
						scale.Set(0.8)
					})
				}).Padding(Bind(NewInset(5))),

				Button("重置", func() {
					WithAnimation(func() {
						opacity.Set(1.0)
						scale.Set(1.0)
					})
				}).Padding(Bind(NewInset(5))),
			),
		).Padding(Bind(NewInset(5))).Debug(),

		// 範例3: 使用 BackgroundColor 綁定動畫
		VStack(
			Text("背景動畫").FontSize(Bind(font.Body)),
			Rectangle().
				Frame(Bind(NewSize(100, 100))).
				BackgroundColor(BindFunc(func() CGColor {
					if enabled.Get() {
						return NewColor(255, 0, 0)
					}
					return NewColor(0, 0, 255)
				}, func(CGColor) {})).
				Padding(Bind(NewInset(5))),
			HStack(
				Text("Toggle"),
				Toggle(enabled).Padding(Bind(NewInset(5))),
			),
		).Padding(Bind(NewInset(5))).Debug(),
	).Padding(Bind(NewInset(5))).
		Center().
		Align(Bind(layout.AlignCenter))
}

func debugFunc(view SomeView) SomeView {
	return view.Modify(func(view SomeView) SomeView {
		return view.
			Padding(Bind(NewInset(5))).
			Border(Bind(NewInset(1)), Bind(NewColor(255, 0, 0)))
	})
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(32))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}
}
