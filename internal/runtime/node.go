package runtime

import (
	"reflect"

	ui "github.com/yanun0323/ebui"
)

type resolvedStyle struct {
	Foreground ui.Color
	Opacity    float64
	Font       ui.TextStyle
	LineLimit  int
	Disabled   bool
}

func defaultResolvedStyle() resolvedStyle {
	return resolvedStyle{
		Foreground: ui.Black,
		Opacity:    1,
		Font:       ui.DefaultTextStyle(),
	}
}

func (s resolvedStyle) with(style ui.Style) resolvedStyle {
	next := s
	if style.HasForeground {
		next.Foreground = style.Foreground
	}
	if style.HasOpacity {
		next.Opacity *= style.Opacity
	}
	if style.HasFont {
		next.Font = style.Font
	}
	if style.HasLineLimit {
		next.LineLimit = style.LineLimit
	}
	if style.HasDisabled {
		next.Disabled = style.Disabled
	}
	return next
}

type buildContext struct {
	owner *node
	seen  map[uintptr]struct{}
}

func (ctx *buildContext) RegisterDependency(signal ui.AnySignal) {
	if signal == nil {
		return
	}
	id := signalIdentity(signal)
	if id != 0 {
		if _, ok := ctx.seen[id]; ok {
			return
		}
		ctx.seen[id] = struct{}{}
	}
	cancel := signal.SubscribeNotify(func() {
		ctx.owner.markDirtyView()
	})
	ctx.owner.depCancels = append(ctx.owner.depCancels, cancel)
}

func signalIdentity(signal ui.AnySignal) uintptr {
	value := reflect.ValueOf(signal)
	if !value.IsValid() {
		return 0
	}
	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.UnsafePointer:
		return value.Pointer()
	default:
		return 0
	}
}

type node struct {
	runtime *Runtime
	parent  *node

	element ui.Element
	kind    ui.ElementKind
	view    ui.View
	props   any

	children []*node
	frame    ui.Rect
	content  ui.Rect
	measured ui.Size
	resolved resolvedStyle

	dirtyView   bool
	dirtyLayout bool
	depCancels  []func()

	scrollOffset ui.Point
	pressed      bool
	dragging     bool
	buildCount   int
}

func newNode(runtime *Runtime, parent *node, element ui.Element) *node {
	n := &node{
		runtime: runtime,
		parent:  parent,
	}
	n.updateFromElement(element)
	return n
}

func (n *node) updateFromElement(element ui.Element) {
	n.element = element
	n.kind = element.Kind()
	n.view = element.View()
	n.props = element.Props()
	n.markLayoutDirty()
	if n.kind == ui.ElementKindView {
		n.markDirtyView()
	}
}

func (n *node) markDirtyView() {
	n.dirtyView = true
	n.markLayoutDirty()
}

func (n *node) markLayoutDirty() {
	n.dirtyLayout = true
	n.runtime.dirtyLayout = true
	n.runtime.dirtyFrame = true
}

func (n *node) cleanupDependencies() {
	for _, cancel := range n.depCancels {
		cancel()
	}
	n.depCancels = nil
}

func (n *node) viewCompatible(next ui.View) bool {
	if n.view == nil || next == nil {
		return false
	}
	return reflect.TypeOf(n.view) == reflect.TypeOf(next)
}

func (n *node) localOffset() ui.Point {
	if !n.element.Style().HasOffset {
		return ui.Point{}
	}
	return n.element.Style().Offset
}

func (n *node) resolve(parent resolvedStyle) resolvedStyle {
	n.resolved = parent.with(n.element.Style())
	return n.resolved
}

func (n *node) borderInsets() ui.Insets {
	style := n.element.Style()
	if !style.HasBorder {
		return ui.Insets{}
	}
	return style.Border.Insets
}

func (n *node) chromeInsets() ui.Insets {
	return n.element.Style().Padding.Add(n.borderInsets())
}

func (n *node) childAlignment(defaultAlign ui.Alignment) ui.Alignment {
	style := n.element.Style()
	if style.HasAlign {
		return style.Align
	}
	return defaultAlign
}
