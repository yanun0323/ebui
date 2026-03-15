# EBUI v2

EBUI v2 is a declarative UI framework for Go with a `View`-first public API, runtime-managed invalidation, and an Ebiten backend.

## Core API

- `type View interface { Body(ctx ViewContext) Element }`
- `State[T]`, `Signal[T]`, `Read(ctx, signal)`
- `Element` as the concrete rendering and layout IR
- `Apply(view, modifiers...)` for function-based modifiers

## Example

```go
type HelloScreen struct{}

func (HelloScreen) Body(ctx ui.ViewContext) ui.Element {
    titleStyle := ui.DefaultTextStyle()
    titleStyle.Size = 28
    titleStyle.Weight = ui.WeightBold

    return ui.Apply(
        ui.ZStack(
            ui.Apply(
                ui.Rectangle(),
                ui.Fill(ui.RGBA(30, 34, 45, 255)),
                ui.Frame(ui.S(640, 360)),
            ),
            ui.Apply(
                ui.Text("哈囉，EBUI v2"),
                ui.Foreground(ui.White),
                ui.Font(titleStyle),
                ui.Padding(ui.Symmetric(12, 18)),
                ui.Background(ui.RGBA(64, 124, 255, 255)),
            ),
        ),
        ui.Frame(ui.S(640, 360)),
        ui.Align(ui.AlignCenter),
    )
}
```

## Run

```bash
go run ./examples/hello
go run ./examples/playground
```

## Verify

```bash
go test ./...
go test -race ./...
go build ./examples/...
```

## Notes

- `backup/` remains the archived legacy implementation.
- The root package is the rebuilt v2 codebase.
- The renderer prefers `resource/Taipei Sans TC Regular.ttf`, then `resource/NotoSansTC.ttf`.
