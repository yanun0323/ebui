package ebui

type ViewContext interface{}

type dependencyRecorder interface {
	RegisterDependency(AnySignal)
}

func Read[T any](ctx ViewContext, signal Signal[T]) T {
	if recorder, ok := ctx.(dependencyRecorder); ok && signal != nil {
		recorder.RegisterDependency(signal)
	}
	if signal == nil {
		var zero T
		return zero
	}
	return signal.Value()
}

type View interface {
	Body(ctx ViewContext) Element
}

type ViewFunc func(ctx ViewContext) Element

func (fn ViewFunc) Body(ctx ViewContext) Element {
	return fn(ctx)
}

type ElementKind uint8

const (
	ElementKindEmpty ElementKind = iota
	ElementKindView
	ElementKindText
	ElementKindImage
	ElementKindRectangle
	ElementKindCircle
	ElementKindDivider
	ElementKindSpacer
	ElementKindVStack
	ElementKindHStack
	ElementKindZStack
	ElementKindScrollView
	ElementKindButton
	ElementKindToggle
	ElementKindSlider
)

type TextProps struct {
	Content string
}

type ImageProps struct {
	Source ImageSource
}

type StackProps struct {
	Direction Direction
}

type ScrollViewProps struct {
	Direction Direction
}

type ButtonProps struct {
	OnPress func()
}

type ToggleProps struct {
	State *State[bool]
}

type SliderProps struct {
	State *State[float64]
	Min   float64
	Max   float64
}

type SpacerProps struct {
	Grow float64
}

type Element struct {
	kind     ElementKind
	view     View
	props    any
	children []Element
	style    Style
	key      string
}

func Empty() Element {
	return Element{kind: ElementKindEmpty}
}

func (e Element) Body(ctx ViewContext) Element {
	return e
}

func Render(view View) Element {
	if view == nil {
		return Empty()
	}
	if element, ok := any(view).(Element); ok {
		return element
	}
	return Element{
		kind: ElementKindView,
		view: view,
	}
}

func Text(content string) Element {
	return Element{
		kind:  ElementKindText,
		props: TextProps{Content: content},
	}
}

func Image(source ImageSource) Element {
	return Element{
		kind:  ElementKindImage,
		props: ImageProps{Source: source},
	}
}

func Rectangle() Element {
	return Element{kind: ElementKindRectangle}
}

func Circle() Element {
	return Element{kind: ElementKindCircle}
}

func Divider() Element {
	return Element{kind: ElementKindDivider}
}

func Spacer(grow ...float64) Element {
	flexGrow := 1.0
	if len(grow) > 0 {
		flexGrow = grow[0]
	}
	return Element{
		kind:  ElementKindSpacer,
		props: SpacerProps{Grow: flexGrow},
	}
}

func VStack(children ...View) Element {
	return Element{
		kind:     ElementKindVStack,
		children: renderChildren(children...),
		props: StackProps{
			Direction: DirectionVertical,
		},
	}
}

func HStack(children ...View) Element {
	return Element{
		kind:     ElementKindHStack,
		children: renderChildren(children...),
		props: StackProps{
			Direction: DirectionHorizontal,
		},
	}
}

func ZStack(children ...View) Element {
	return Element{
		kind:     ElementKindZStack,
		children: renderChildren(children...),
	}
}

func ScrollView(content View, direction ...Direction) Element {
	dir := DirectionVertical
	if len(direction) > 0 {
		dir = direction[0]
	}
	return Element{
		kind:     ElementKindScrollView,
		children: []Element{Render(content)},
		props: ScrollViewProps{
			Direction: dir,
		},
	}
}

func Button(label View, onPress func()) Element {
	return Element{
		kind:     ElementKindButton,
		children: []Element{Render(label)},
		props: ButtonProps{
			OnPress: onPress,
		},
	}
}

func Toggle(state *State[bool], label ...View) Element {
	return Element{
		kind:     ElementKindToggle,
		children: renderChildren(label...),
		props: ToggleProps{
			State: state,
		},
	}
}

func Slider(state *State[float64], min, max float64) Element {
	if max < min {
		min, max = max, min
	}
	return Element{
		kind: ElementKindSlider,
		props: SliderProps{
			State: state,
			Min:   min,
			Max:   max,
		},
	}
}

func renderChildren(children ...View) []Element {
	rendered := make([]Element, 0, len(children))
	for _, child := range children {
		rendered = append(rendered, Render(child))
	}
	return rendered
}

func (e Element) Kind() ElementKind {
	return e.kind
}

func (e Element) View() View {
	return e.view
}

func (e Element) Props() any {
	return e.props
}

func (e Element) Children() []Element {
	return append([]Element(nil), e.children...)
}

func (e Element) Style() Style {
	return e.style
}

func (e Element) Key() string {
	return e.key
}

func (e Element) WithKey(key string) Element {
	e.key = key
	return e
}
