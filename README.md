> [!IMPORTANT]
> This package is presently in its alpha stage of development

# EBUI

EBUI is a declarative UI framework for Go, inspired by [SwiftUI](https://developer.apple.com/documentation/swiftui) and built on top of the [Ebitengine](https://github.com/hajimehoshi/ebiten) framework. It enables developers to create interactive applications with a clean, functional syntax.

## Features

- **Declarative Syntax**: Build UIs with a clean, SwiftUI-like declarative syntax
- **Data Binding**: Reactive programming model with bindings to update UI automatically
- **Component-Based Architecture**: Create reusable UI components
- **Layout System**: Flexible stack-based layout system (VStack, HStack, ZStack)
- **Modifiers**: Chain modifiers for styling and behavior customization
- **Live Preview**: Preview changes in real-time with VSCode integration
- **Cross-Platform**: Run on any platform supported by Ebitengine
- **Integrated Animation**: Built-in support for smooth UI animations
- **Ebitengine Integration**: Works as a standalone app or inside existing Ebitengine projects

## Installation

```
# Coming soon to package managers
```

## Quick Start

### Define a Content View

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

### Run as Standalone App

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

### Run within Ebitengine

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

## Available Components

### Views (✓ = Implemented)

| Basic Components | Layout Components | Input Components | Advanced Components |
| ---------------- | ----------------- | ---------------- | ------------------- |
| ✓ Text           | ✓ VStack          | ✓ Button         | ✓ ScrollView        |
| ✓ Image          | ✓ HStack          | ✓ Toggle         | ❑ List              |
| ✓ Rectangle      | ✓ ZStack          | ✓ Slider         | ❑ TableView         |
| ✓ Circle         | ✓ Spacer          | ❑ TextField      | ❑ Navigation        |
| ✓ Divider        | ✓ EmptyView       | ❑ TextEditor     | ❑ Sheet             |
| ✓ ViewModifier   |                   | ❑ Picker         | ❑ Menu              |
|                  |                   | ❑ Radio          | ❑ ProgressView      |
|                  |                   | ❑ DatePicker     |                     |
|                  |                   | ❑ TimePicker     |                     |
|                  |                   | ❑ ColorPicker    |                     |

### Features (✓ = Implemented)

- ✓ Modifier Stack
- ✓ CornerRadius
- ✓ Animation
- ✓ Alignment
- ❑ Gesture
- ❑ Overlay
- ❑ Mask
- ❑ Clip

## Live Preview with VSCode

EBUI offers a seamless development experience with real-time UI previews directly in VSCode, similar to SwiftUI's preview functionality.

### Using the VSCode Extension

The [ebui-vscode extension](https://github.com/yanun0323/ebui-vscode) enables hot-reloading for any function that starts with `Preview_`.

1. Install the extension from the VSCode marketplace
2. Create functions prefixed with `Preview_` in your Go files
3. Save to see your UI update instantly in a preview window

### Example

```go
package mypackage

import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

// This function will be automatically previewed
func Preview_MyButton() ui.View {
	return ui.Button(ui.Const("Click Me")).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{200, 100, 100, 255})).
		Padding(ui.Bind(ui.Inset{10, 10, 10, 10})).
		Center()
}
```

### How It Works

The extension:

1. Watches for changes in your Go files
2. Runs the EBUI preview tool when you save
3. Detects any `Preview_` functions in your code
4. Renders them in a live preview window

This creates a fast feedback loop for UI development without leaving your editor.
