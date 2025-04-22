> [!IMPORTANT]
> 此包目前处于 alpha 开发阶段

# EBUI

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

EBUI 是一个受 [SwiftUI](https://developer.apple.com/documentation/swiftui) 启发，建立在 [Ebitengine](https://github.com/hajimehoshi/ebiten) 框架之上的声明式 UI 框架，专为 Go 语言打造。它让开发者能够使用简洁、功能性的语法创建交互式应用程序。

## 特色功能

- **声明式语法**：使用类似 SwiftUI 的简洁声明式语法构建 UI
- **数据绑定**：响应式编程模型，自动更新 UI
- **组件化架构**：创建可重用的 UI 组件
- **布局系统**：灵活的堆栈式布局系统（VStack、HStack、ZStack）
- **修饰符链**：串联修饰符以自定义样式和行为
- **实时预览**：通过 VSCode 集成实现实时预览变更
- **跨平台**：在任何 Ebitengine 支持的平台上运行
- **内置动画**：支持流畅的 UI 动画
- **Ebitengine 集成**：可作为独立应用或嵌入现有 Ebitengine 项目中运行

## 安装

```
# 即将推出至包管理器
```

## 快速入门

### 定义内容视图

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

### 作为独立应用运行

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

### 在 Ebitengine 中运行

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

## 可用组件

### 视图 (✓ = 已实现)

| 基础组件                     | 布局组件              | 输入组件                    | 高级组件                 |
| ---------------------------- | --------------------- | --------------------------- | ------------------------ |
| ✓ Text（文本）               | ✓ VStack（垂直堆栈）  | ✓ Button（按钮）            | ✓ ScrollView（滚动视图） |
| ✓ Image（图片）              | ✓ HStack（水平堆栈）  | ✓ Toggle（开关）            | ❑ List（列表）           |
| ✓ Rectangle（矩形）          | ✓ ZStack（层叠堆栈）  | ✓ Slider（滑块）            | ❑ TableView（表格视图）  |
| ✓ Circle（圆形）             | ✓ Spacer（间隔器）    | ❑ TextField（文本输入框）   | ❑ Navigation（导航）     |
| ✓ Divider（分隔线）          | ✓ EmptyView（空视图） | ❑ TextEditor（文本编辑器）  | ❑ Sheet（底部弹出）      |
| ✓ ViewModifier（视图修饰符） |                       | ❑ Picker（选择器）          | ❑ Menu（菜单）           |
|                              |                       | ❑ Radio（单选）             | ❑ ProgressView（进度条） |
|                              |                       | ❑ DatePicker（日期选择器）  |                          |
|                              |                       | ❑ TimePicker（时间选择器）  |                          |
|                              |                       | ❑ ColorPicker（颜色选择器） |                          |

### 功能 (✓ = 已实现)

- ✓ Modifier Stack（修饰符堆栈）
- ✓ CornerRadius（圆角）
- ✓ Animation（动画）
- ✓ Alignment（对齐）
- ❑ Gesture（手势）
- ❑ Overlay（覆盖层）
- ❑ Mask（遮罩）
- ❑ Clip（裁剪）

## VSCode 实时预览

EBUI 提供类似 SwiftUI 的开发体验，让你能在 VSCode 中直接进行 UI 的实时预览。

### 使用 VSCode 扩展

[ebui-vscode 扩展](https://github.com/yanun0323/ebui-vscode) 提供热重载功能，自动检测以 `Preview_` 开头的函数。

1. 从 VSCode 市场安装扩展
2. 在 Go 文件中创建以 `Preview_` 为前缀的函数
3. 保存后即可在预览窗口中实时查看 UI 更新

### 示例

```go
package mypackage

import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

// 此函数将自动预览
func Preview_MyButton() ui.View {
	return ui.Button(ui.Const("点击我")).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{200, 100, 100, 255})).
		Padding(ui.Bind(ui.Inset{10, 10, 10, 10})).
		Center()
}
```

### 工作原理

扩展的工作流程：

1. 监控 Go 文件的变更
2. 保存时运行 EBUI 预览工具
3. 检测代码中任何 `Preview_` 函数
4. 在实时预览窗口中呈现

这为 UI 开发提供了快速反馈循环，无需离开编辑器。
