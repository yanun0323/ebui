package ebui

import "fmt"

type CounterView struct {
	count Binding[int]
}

func (v *CounterView) Body() SomeView {
	return VStack(
		Text(fmt.Sprintf("Count: %d", v.count.Get())),
		Button(func() {
			v.count.Set(v.count.Get() + 1)
		}, Text("Increment")),
		Button(func() {
			v.count.Set(v.count.Get() - 1)
		}, Text("Decrement")),
	).padding(10).frame(200, 100)
}
