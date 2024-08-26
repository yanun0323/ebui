package ebui

/* Check Interface Implementation */
var _ SomeView = (*emptyView)(nil)

func Empty() *emptyView {
	v := &emptyView{}
	v.view = newView(typeEmpty, v)
	v.view.w, v.view.h = 0, 0
	return v
}

type emptyView struct {
	*view
}
