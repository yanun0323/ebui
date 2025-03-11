package main

import (
	"fmt"
	"log"
	"time"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/font"
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
	return VStack(
		// 範例1: 基本動畫 - 使用直接設置帶動畫參數的方式
		VStack(
			Text("基本動畫範例").FontSize(Bind(font.Body)),
			HStack(
				Button("計數+1 (線性)", func() {
					// 直接在Set時傳入動畫參數
					count.Set(count.Get()+1, animation.Linear(time.Millisecond*500))
				}).Padding(Bind(NewInset(10, 10, 10, 10))),

				Button("計數-1 (緩入緩出)", func() {
					count.Set(count.Get()-1, animation.EaseInOut(time.Millisecond*500))
				}).Padding(Bind(NewInset(10, 10, 10, 10))),
			),
			Text(BindFunc(func() string {
				return fmt.Sprintf("計數: %.2f", count.Get())
			}, nil)).FontSize(Bind(font.Body)).Padding(Bind(NewInset(10, 10, 10, 10))),
		).Padding(Bind(NewInset(5, 5, 5, 5))),

		// 範例2: 使用WithAnimation上下文 - SwiftUI風格
		VStack(
			Text("上下文動畫範例").FontSize(Bind(font.Body)),
			Rectangle().
				Frame(Bind(NewSize(100, 100))).
				BackgroundColor(Bind(NewColor(255, 0, 0, 255))).
				Padding(Bind(NewInset(10, 10, 10, 10))).
				Opacity(opacity),
			HStack(
				Button("淡出", func() {
					// 使用動畫上下文包裹多個更改
					animation.WithAnimation(animation.EaseOut(time.Millisecond*800), func() {
						opacity.Set(0.2)
						scale.Set(0.8)
					})
				}).Padding(Bind(NewInset(10, 10, 10, 10))),

				Button("重置", func() {
					animation.WithAnimation(animation.EaseIn(time.Millisecond*500), func() {
						opacity.Set(1.0)
						scale.Set(1.0)
					})
				}).Padding(Bind(NewInset(10, 10, 10, 10))),
			),
		).Padding(Bind(NewInset(5, 5, 5, 5))),

		// 範例3: 使用 BackgroundColor 綁定動畫
		VStack(
			Text("背景動畫範例").FontSize(Bind(font.Body)),
			Rectangle().
				Frame(Bind(NewSize(100, 100))).
				BackgroundColor(BindFunc(func() CGColor {
					if enabled.Get() {
						return NewColor(255, 0, 0, 255)
					}
					return NewColor(0, 0, 255, 255)
				}, func(CGColor) {})).
				Padding(Bind(NewInset(10, 10, 10, 10))),
			HStack(
				Toggle("開關", enabled).Padding(Bind(NewInset(10, 10, 10, 10))),
			),
			Rectangle().
				Padding(Bind(NewInset(5, 5, 5, 5))),
		).Padding(Bind(NewInset(5, 5, 5, 5))),
	).Padding(Bind(NewInset(5, 5, 5, 5)))
}

func main() {
	app := NewApplication(NewContentView())
	app.SetWindowBackgroundColor(NewColor(32, 32, 32, 32))
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(WindowResizingModeEnabled)
	app.SetResourceFolder("resource")
	app.Debug()

	if err := app.Run("Counter Demo"); err != nil {
		log.Fatal(err)
	}
}
