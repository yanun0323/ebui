package ebui

type rectangle struct {
	view
}

func Rectangle() SomeView {
	v := &rectangle{}
	v.view = newView(v)
	return v
}
