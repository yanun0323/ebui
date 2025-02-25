package ebui

type rectangleImpl struct {
	*viewCtx
}

func Rectangle() SomeView {
	rect := &rectangleImpl{}
	rect.viewCtx = newViewContext(rect)
	return rect
}
