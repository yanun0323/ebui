> [!IMPORTANT]
> 此套件目前處於 alpha 開發階段

# EBUI

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

EBUI 是一個受 [SwiftUI](https://developer.apple.com/documentation/swiftui) 啟發，建立在 [Ebitengine](https://github.com/hajimehoshi/ebiten) 框架之上的聲明式 UI 框架，專為 Go 語言打造。它讓開發者能夠使用簡潔、功能性的語法創建互動式應用程式。

## 特色功能

- **聲明式語法**：使用類似 SwiftUI 的簡潔聲明式語法構建 UI
- **資料綁定**：反應式程式設計模型，自動更新 UI
- **元件化架構**：創建可重複使用的 UI 元件
- **佈局系統**：靈活的堆疊式佈局系統（VStack、HStack、ZStack）
- **修飾符鏈**：串聯修飾符以自訂樣式和行為
- **即時預覽**：通過 VSCode 整合實現即時預覽變更
- **跨平台**：在任何 Ebitengine 支援的平台上運行
- **內建動畫**：支援流暢的 UI 動畫
- **Ebitengine 整合**：可作為獨立應用或嵌入現有 Ebitengine 專案中運行

## 安裝

```
# 即將推出至套件管理器
```

## 快速入門

### 定義內容視圖

```go
import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

type contentView struct {
	title   *ui.Binding[string]
	content *ui.Binding[string]
}

func ContentView(title, content *ui.Binding[string]) ui.View {
	return &contentView{
		title:   title,
		content: content,
	}
}

func (view *contentView) Body() ui.SomeView {
	return ui.HStack(
		ui.Spacer(),
		ui.VStack(
			ui.Spacer(),
			ui.Text(view.title).
				Padding(ui.Bind(ui.Inset{0, 15, 0, 15})).
				ForegroundColor(ui.Bind[color.Color](color.White)).
				BackgroundColor(ui.Bind[color.Color](color.Gray{128})),
			ui.Text(view.content),
			ui.Spacer(),
		).Frame(ui.Bind(200.0), nil),
		ui.Spacer(),
	).
		ForegroundColor(ui.Bind[color.Color](color.RGBA{200, 200, 200, 255})).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{255, 0, 0, 255})).
		Padding(ui.Bind(ui.Inset{5, 5, 5, 5}))
}
```

### 作為獨立應用運行

```go
import (
	ui "github.com/yanun0323/ebui"
	"image/color"
	"log"
)

func main() {
	title := ui.Bind("title")
	content := ui.Bind("content")
	contentView := ContentView(title, content)

	app := ui.NewApplication(contentView)
	app.SetWindowBackgroundColor(color.RGBA{100, 0, 0, 0})
	app.SetWindowTitle("EBUI Demo")
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(ui.WindowResizingModeEnabled)
	app.SetResourceFolder("resource")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### 在 Ebitengine 中運行

```go
import (
	ui "github.com/yanun0323/ebui"
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"image/color"
)

func main() {
	contentView := ui.ContentView("title", "content")
	g := NewGame(contentView)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("run game", "error", err)
	}
}

func NewGame(contentView ui.View) *Game {
	return &Game{
		contentView: contentView.Body(),
	}
}

type Game struct {
	contentView ui.SomeView
}

func (g *Game) Update() error {
	ui.EbitenUpdate(g.contentView)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ui.EbitenDraw(screen, g.contentView)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ui.EbitenLayout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
```

## 可用元件

### 視圖 (✓ = 已實現)

| 基礎元件                     | 佈局元件              | 輸入元件                    | 進階元件                 |
| ---------------------------- | --------------------- | --------------------------- | ------------------------ |
| ✓ Text（文字）               | ✓ VStack（垂直堆疊）  | ✓ Button（按鈕）            | ✓ ScrollView（滾動視圖） |
| ✓ Image（圖片）              | ✓ HStack（水平堆疊）  | ✓ Toggle（開關）            | ❑ List（列表）           |
| ✓ Rectangle（矩形）          | ✓ ZStack（疊層堆疊）  | ✓ Slider（滑桿）            | ❑ TableView（表格視圖）  |
| ✓ Circle（圓形）             | ✓ Spacer（間隔器）    | ❑ TextField（文字輸入框）   | ❑ Navigation（導航）     |
| ✓ Divider（分隔線）          | ✓ EmptyView（空視圖） | ❑ TextEditor（文字編輯器）  | ❑ Sheet（底部彈出）      |
| ✓ ViewModifier（視圖修飾符） |                       | ❑ Picker（選擇器）          | ❑ Menu（選單）           |
|                              |                       | ❑ Radio（單選）             | ❑ ProgressView（進度條） |
|                              |                       | ❑ DatePicker（日期選擇器）  |                          |
|                              |                       | ❑ TimePicker（時間選擇器）  |                          |
|                              |                       | ❑ ColorPicker（顏色選擇器） |                          |

### 功能 (✓ = 已實現)

- ✓ Modifier Stack（修飾符堆疊）
- ✓ CornerRadius（圓角）
- ✓ Animation（動畫）
- ✓ Alignment（對齊）
- ❑ Gesture（手勢）
- ❑ Overlay（覆蓋層）
- ❑ Mask（遮罩）
- ❑ Clip（裁剪）

## VSCode 即時預覽

EBUI 提供類似 SwiftUI 的開發體驗，讓你能在 VSCode 中直接進行 UI 的即時預覽。

### 使用 VSCode 擴充功能

[ebui-vscode 擴充功能](https://github.com/yanun0323/ebui-vscode) 提供熱重載功能，自動偵測以 `Preview_` 開頭的函數。

1. 從 VSCode 市集安裝擴充功能
2. 在 Go 檔案中創建以 `Preview_` 為前綴的函數
3. 儲存後即可在預覽視窗中即時查看 UI 更新

### 範例

```go
package mypackage

import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

// 此函數將自動預覽
func Preview_MyButton() ui.View {
	return ui.Button(ui.Const("點擊我")).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{200, 100, 100, 255})).
		Padding(ui.Bind(ui.Inset{10, 10, 10, 10})).
		Center()
}
```

### 工作原理

擴充功能的工作流程：

1. 監控 Go 檔案的變更
2. 儲存時運行 EBUI 預覽工具
3. 偵測代碼中任何 `Preview_` 函數
4. 在即時預覽視窗中呈現

這為 UI 開發提供了快速反饋循環，無需離開編輯器。
