package ebui

/* Check Interface Implementation */
var _ SomeView = (*emptyView)(nil)

func Empty() *emptyView {
	v := &emptyView{}
	v.uiView = newUIView(typesEmpty, v)
	v.uiView.size.w, v.uiView.size.h = 0, 0
	return v
}

type emptyView struct {
	*uiView
}
