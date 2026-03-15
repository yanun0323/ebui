package runtime

import (
	"fmt"
	"testing"

	ui "github.com/yanun0323/ebui"
)

type fakeMeasurer struct{}

func (fakeMeasurer) MeasureText(content string, style ui.TextStyle, lineLimit int) ui.Size {
	return ui.S(float64(len(content))*style.Size*0.5, style.Size+style.LineHeight)
}

func (fakeMeasurer) ImageSize(source ui.ImageSource) ui.Size {
	return ui.S(100, 100)
}

type staticChild struct{}

func (staticChild) Body(ctx ui.ViewContext) ui.Element {
	return ui.Text("static")
}

type dynamicChild struct {
	value *ui.State[int]
}

func (c *dynamicChild) Body(ctx ui.ViewContext) ui.Element {
	value := ui.Read(ctx, c.value)
	return ui.Text(fmt.Sprintf("value:%d", value))
}

type rootView struct {
	static  staticChild
	dynamic *dynamicChild
}

func (c *rootView) Body(ctx ui.ViewContext) ui.Element {
	return ui.VStack(c.static, c.dynamic)
}

func TestSelectiveViewRebuild(t *testing.T) {
	value := ui.NewState(1)
	root := &rootView{dynamic: &dynamicChild{value: value}}
	rt := New(root, fakeMeasurer{})
	dynamicNode := FindFirstView(rt.root, root.dynamic)
	staticNode := FindFirstView(rt.root, root.static)
	if dynamicNode == nil || staticNode == nil {
		t.Fatalf("expected mounted view nodes")
	}

	value.Set(2)
	rt.Update(nil)

	if dynamicNode.buildCount < 2 {
		t.Fatalf("expected dynamic child rebuild, got %d", dynamicNode.buildCount)
	}
	if staticNode.buildCount != 1 {
		t.Fatalf("expected static child to remain at one build, got %d", staticNode.buildCount)
	}
}

type buttonView struct {
	count *ui.State[int]
}

func (c *buttonView) Body(ctx ui.ViewContext) ui.Element {
	return ui.Apply(
		ui.Button(ui.Text("Tap"), func() {
			c.count.Update(func(v int) int { return v + 1 })
		}),
		ui.Frame(ui.S(120, 44)),
	)
}

func TestButtonClick(t *testing.T) {
	count := ui.NewState(0)
	rt := New(&buttonView{count: count}, fakeMeasurer{})
	rt.Resize(ui.S(160, 80))
	rt.Update([]Event{
		{Kind: EventPointerDown, Point: ui.P(10, 10)},
		{Kind: EventPointerUp, Point: ui.P(10, 10)},
	})
	if got := count.Value(); got != 1 {
		t.Fatalf("expected one click, got %d", got)
	}
}

type toggleView struct {
	state *ui.State[bool]
}

func (c *toggleView) Body(ctx ui.ViewContext) ui.Element {
	return ui.Apply(ui.Toggle(c.state), ui.Frame(ui.S(60, 32)))
}

func TestToggleClick(t *testing.T) {
	state := ui.NewState(false)
	rt := New(&toggleView{state: state}, fakeMeasurer{})
	rt.Resize(ui.S(100, 60))
	rt.Update([]Event{
		{Kind: EventPointerDown, Point: ui.P(20, 20)},
		{Kind: EventPointerUp, Point: ui.P(20, 20)},
	})
	if !state.Value() {
		t.Fatal("expected toggle to switch on")
	}
}

type sliderView struct {
	state *ui.State[float64]
}

func (c *sliderView) Body(ctx ui.ViewContext) ui.Element {
	return ui.Apply(ui.Slider(c.state, 0, 100), ui.Frame(ui.S(300, 32)))
}

func TestSliderClamp(t *testing.T) {
	state := ui.NewState(0.0)
	rt := New(&sliderView{state: state}, fakeMeasurer{})
	rt.Resize(ui.S(320, 50))
	rt.Update([]Event{
		{Kind: EventPointerDown, Point: ui.P(10, 16)},
		{Kind: EventPointerMove, Point: ui.P(1000, 16)},
		{Kind: EventPointerUp, Point: ui.P(1000, 16)},
	})
	if got := state.Value(); got != 100 {
		t.Fatalf("expected clamped slider value, got %f", got)
	}
}

type scrollView struct{}

func (scrollView) Body(ctx ui.ViewContext) ui.Element {
	items := make([]ui.View, 0, 20)
	for i := 0; i < 20; i++ {
		items = append(items, ui.Apply(ui.Text("item"), ui.Frame(ui.S(100, 24))))
	}
	return ui.Apply(
		ui.ScrollView(ui.Apply(ui.VStack(items...), ui.Spacing(4)), ui.DirectionVertical),
		ui.Frame(ui.S(120, 120)),
	)
}

func TestScrollTargetReceivesWheel(t *testing.T) {
	rt := New(scrollView{}, fakeMeasurer{})
	rt.Resize(ui.S(200, 200))
	node := rt.hitScroll(rt.root, ui.P(10, 10))
	if node == nil {
		t.Fatal("expected scroll target")
	}
	before := node.scrollOffset.Y
	rt.Update([]Event{{Kind: EventWheel, Point: ui.P(10, 10), Scroll: ui.P(0, -1)}})
	if node.scrollOffset.Y <= before {
		t.Fatalf("expected scroll offset to increase, got %f", node.scrollOffset.Y)
	}
}
