package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Application struct {
	StateManager    *StateManager
	EventManager    *EventManager
	AnimManager     *AnimationManager
	RootView        View
	BackgroundColor color.Color
	Bounds          image.Rectangle
}

func (app *Application) Update() error {
	// 1. 更新動畫
	app.AnimManager.Update()

	// 2. 處理狀態更新
	if app.StateManager.isDirty() {
		// 重新構建視圖樹
		app.RootView = app.RootView.Build()
		// 重新計算佈局
		app.RootView.Layout(app.Bounds)
	}

	// 3. 處理輸入事件
	app.handleInput()

	return nil
}

func (app *Application) handleInput() {
	// 處理輸入
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		app.EventManager.DispatchTouchEvent(TouchEvent{
			Phase:    TouchPhaseBegan,
			Position: image.Pt(x, y),
		})
	}
}

func (app *Application) Draw(screen *ebiten.Image) {
	// 清空畫面
	screen.Fill(app.BackgroundColor)

	// 繪製視圖樹
	app.RootView.Draw(screen)
}
