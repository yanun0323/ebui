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
	. "github.com/yanun0323/ebui"
)

type contentView struct {
	title   *Binding[string]
	content *Binding[string]
}

func ContentView(title, content *Binding[string]) View {
	return &contentView{
		title:   title,
		content: content,
	}
}

func (view *contentView) Body() SomeView {
	return HStack(
		Spacer(),
		VStack(
			Spacer(),
			Text(view.title).
				Padding(Bind(Inset{0, 15, 0, 15})).
				ForegroundColor(Bind[color.Color](color.White)).
				BackgroundColor(Bind[color.Color](color.Gray{128})),
			Text(view.content),
			Spacer(),
		).Frame(Bind(200.0), nil),
		Spacer(),
	).
		ForegroundColor(Bind[color.Color](color.RGBA{200, 200, 200, 255})).
		BackgroundColor(Bind[color.Color](color.RGBA{255, 0, 0, 255})).
		Padding(Bind(Inset{5, 5, 5, 5}))
}
```

### Run an App

```go
func main() {
	title := Bind("title")
	content := Bind("content")
	contentView := ContentView(title, content)

	app := NewApplication(contentView)
	app.SetWindowBackgroundColor(color.RGBA{100, 0, 0, 0})
	app.SetWindowTitle("EBUI Demo")
	app.SetWindowSize(600, 500)
	app.SetWindowResizingMode(ebui.WindowResizingModeEnabled)
	app.SetResourceFolder("resource")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### Run in a Ebitengine Application

```go
func main() {
	contentView := component.ContentView("title", "content")
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
	EbitenUpdate(g.contentView)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	EbitenDraw(screen, g.contentView)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
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
- [] EmptyView
- [ ] Circle
- [x] Rectangle
- [ ] Divider
- [ ] Indicator
- [ ] Menu
- [ ] Sheet
- [ ] TextField
- [ ] TextEditor
- [ ] ScrollView
- [ ] ListView
- [ ] TableView
- [ ] Slider
- [ ] Toggle
- [ ] Navigation
- [ ] Checkbox
- [ ] Radio
- [ ] Picker
- [ ] DatePicker
- [ ] TimePicker
- [ ] ColorPicker
- [ ] ProgressView

### Feature

- [x] CornerRadius
- [ ] Animation
- [ ] Gesture
- [ ] Overlay
- [ ] Mask
- [ ] Clip
