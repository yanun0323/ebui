> [!IMPORTANT]
> This package is presently in its alpha stage of development

# EBUI

#### EBUI is a [SwiftUI](https://developer.apple.com/documentation/swiftui) likes UI framework build with golang and [Ebitengine](https://github.com/hajimehoshi/ebiten) framework.

For develop games, apps, websites, and other interactive applications.

It also works within a [Ebitengine](https://github.com/hajimehoshi/ebiten) application.

## Install

### Coming soon...

## Example

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

### Run an App

```go
import (
	ui "github.com/yanun0323/ebui"
	"image/color"
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

### Run in a Ebitengine Application

```go
import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

func main() {
	contentView := ui.ContentView("title", "content")
	g := NewGame(contentView)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("run game", "error", err)
	}
}

func NewGame(contentView View) *Game {
	return &Game{
		contentView: contentView.Body(),
	}
}

type Game struct {
	contentView SomeView
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

## Developing Status

### View

- [x] VStack
- [x] HStack
- [x] ZStack
- [x] Text
- [x] Image
- [x] Button
- [x] Spacer
- [x] EmptyView
- [x] Circle
- [x] Rectangle
- [x] ViewModifier
- [x] Divider
- [x] ScrollView
- [x] Toggle
- [x] Slider
- [ ] List
- [ ] TableView
- [ ] Menu
- [ ] Sheet
- [ ] Navigation
- [ ] Radio
- [ ] Picker
- [ ] DatePicker
- [ ] TimePicker
- [ ] ColorPicker
- [ ] ProgressView
- [ ] TextField
- [ ] TextEditor

### Feature

- [x] Modifier Stack
- [x] CornerRadius
- [x] Animation
- [x] Alignment
- [ ] Gesture
- [ ] Overlay
- [ ] Mask
- [ ] Clip

## Live Preview with VSCode

EBUI supports live preview functionality similar to SwiftUI, allowing you to see your UI changes in real-time without restarting your application.

### Using the VSCode Extension

You can use the [ebui-vscode extension](https://github.com/yanun0323/ebui-vscode) to automatically hot-reload any function that starts with `Preview_`. This enables SwiftUI-like instant previews of your UI components.

1. Install the extension from the VSCode marketplace or from the repository
2. Open your project using EBUI as a framework
3. Create functions that start with `Preview_` to define your preview components
4. Save your file to automatically see changes in a preview window

### Example

Create a file with a `Preview_` function:

```go
package mypackage

import (
	ui "github.com/yanun0323/ebui"
	"image/color"
)

// This function will be automatically detected and previewed
func Preview_MyButton() ui.View {
	return ui.Button(ui.Const("Click Me")).
		BackgroundColor(ui.Bind[color.Color](color.RGBA{200, 100, 100, 255})).
		Padding(ui.Bind(ui.Inset{10, 10, 10, 10})).
		Center()
}
```

Similarly, as shown in the `ui_text.go` file:

```go
func Preview_Text() ui.View {
	return ui.Text(ui.Const("Hello, World!")).Center()
}
```

When you save the file, the VSCode extension will automatically run the EBUI tool and display a preview window with your component.

### How It Works

The extension:

1. Monitors changes to your Go files
2. Automatically runs the EBUI tool (`ebui -f ${FILEPATH}`) when saving
3. Detects functions starting with `Preview_`
4. Generates a temporary application to render the preview
5. Shows the preview in a floating window

This workflow allows you to iteratively design and test your UI components without leaving your editor.
