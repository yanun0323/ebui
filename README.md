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
type contentView struct {
	title   string
	content string
}

func ContentView(title, content string) ebui.View {
	return &contentView{
		title:   title,
		content: content,
	}
}

func (view *contentView) Body() ebui.SomeView {
	return ebui.HStack(
		ebui.Spacer(),
		ebui.VStack(
			ebui.Spacer(),
			ebui.Text(view.title).
				Padding(0, 15, 0, 15).
				ForegroundColor(color.White).
				BackgroundColor(color.Gray{128}),
			ebui.Text(view.content),
			ebui.Spacer(),
		).Frame(200, -1),
		ebui.Spacer(),
	).
		ForegroundColor(color.RGBA{200, 200, 200, 255}).
		BackgroundColor(color.RGBA{255, 0, 0, 255}).
		Padding(5, 5, 5, 5)
}
```

### Run an App

```go
func main() {
	contentView := ContentView("title", "content")
	ebui.Run("Windows Title", contentView)
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

func NewGame(contentView ebui.View) *Game {
	return &Game{
		contentView: contentView.Body(),
	}
}

type Game struct {
	contentView ebui.SomeView
}

func (g *Game) Update() error {
	ebui.EbitenUpdate(g.contentView)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebui.EbitenDraw(screen, g.contentView)
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
- [x] EmptyView
- [x] Text
- [x] Image
- [x] Button
- [x] Spacer
- [ ] Circle
- [ ] Rectangle
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

- [ ] Animation
- [ ] Gesture
- [ ] Overlay
- [ ] Mask
- [ ] Clip
