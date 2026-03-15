# EBUI Legacy to v2

## Public naming

- `Component` -> `View`
- `BuildContext` -> `ViewContext`
- `Build(ctx)` -> `Body(ctx)`
- `Use(view)` -> `Render(view)` for explicit lowering only

## Modifier model

v2 uses function-based modifiers:

```go
ui.Apply(
    ui.Text("Hello"),
    ui.Foreground(ui.White),
    ui.Padding(ui.All(8)),
)
```

Instead of method chaining:

```go
ui.Text("Hello").Foreground(...).Padding(...)
```

## State

- Dynamic values live in `State[T]` / `Signal[T]`
- Views read them with `ui.Read(ctx, signal)`
- `ViewContext` remains opaque and runtime-managed
