package ebui

type rectangle struct {
	*view
}

func Rectangle() SomeView {
	v := &rectangle{}
	v.view = newView(idRectangle, v)
	return v
}
