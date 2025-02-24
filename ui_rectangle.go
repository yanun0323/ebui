package ebui

type rectangleImpl struct {
	*ctx
}

func Rectangle() SomeView {
	rect := &rectangleImpl{}
	rect.ctx = newViewContext(rect)
	return rect
}
