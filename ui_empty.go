package ebui

/* Check Interface Implementation */
var _ SomeView = (*emptyView)(nil)

func Empty() *emptyView {
	v := &emptyView{}
	v.uiView = newView(typesEmpty, v)
	return v
}

type emptyView struct {
	*uiView
}

func (p *emptyView) getSize() size {
	return size{}
}

func (p *emptyView) Frame(w, h int) SomeView {
	return p
}
