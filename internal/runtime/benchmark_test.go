package runtime

import (
	"fmt"
	"testing"

	ui "github.com/yanun0323/ebui"
)

type listView struct {
	items int
}

func (c listView) Body(ctx ui.ViewContext) ui.Element {
	children := make([]ui.View, 0, c.items)
	for i := 0; i < c.items; i++ {
		children = append(children, ui.Apply(ui.Text(fmt.Sprintf("item-%d", i)), ui.Frame(ui.S(80, 18))))
	}
	return ui.Apply(ui.VStack(children...), ui.Spacing(2))
}

type singleLeafView struct {
	value *ui.State[int]
}

func (c *singleLeafView) Body(ctx ui.ViewContext) ui.Element {
	children := make([]ui.View, 0, 1000)
	for i := 0; i < 999; i++ {
		children = append(children, ui.Apply(ui.Text("static"), ui.Frame(ui.S(80, 18))))
	}
	children = append(children, ui.Apply(ui.Text(fmt.Sprintf("dynamic:%d", ui.Read(ctx, c.value))), ui.Frame(ui.S(80, 18))))
	return ui.Apply(ui.VStack(children...), ui.Spacing(2))
}

func BenchmarkStaticTreeNoChange(b *testing.B) {
	rt := New(listView{items: 1000}, fakeMeasurer{})
	rt.Resize(ui.S(400, 600))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.Update(nil)
	}
}

func BenchmarkSingleLeafStateChange(b *testing.B) {
	state := ui.NewState(0)
	rt := New(&singleLeafView{value: state}, fakeMeasurer{})
	rt.Resize(ui.S(400, 600))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state.Set(i)
		rt.Update(nil)
	}
}

func BenchmarkTextHeavyListScroll(b *testing.B) {
	rt := New(ui.ViewFunc(func(ctx ui.ViewContext) ui.Element {
		return ui.Apply(
			ui.ScrollView(listView{items: 500}, ui.DirectionVertical),
			ui.Frame(ui.S(400, 300)),
		)
	}), fakeMeasurer{})
	rt.Resize(ui.S(400, 300))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.Update([]Event{{Kind: EventWheel, Point: ui.P(10, 10), Scroll: ui.P(0, -1)}})
	}
}

func BenchmarkSliderDrag(b *testing.B) {
	state := ui.NewState(0.0)
	rt := New(&sliderView{state: state}, fakeMeasurer{})
	rt.Resize(ui.S(320, 50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		point := ui.P(float64(i%300), 16)
		rt.Update([]Event{
			{Kind: EventPointerDown, Point: point},
			{Kind: EventPointerMove, Point: point},
			{Kind: EventPointerUp, Point: point},
		})
	}
}
