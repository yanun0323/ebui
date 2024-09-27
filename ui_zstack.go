package ebui

/* Check Interface Implementation */
var _ SomeView = (*zstackView)(nil)

func ZStack(views ...View) SomeView {
	v := &zstackView{}
	v.uiView = newView(typesZStack, v, views...)
	return v
}

type zstackView struct {
	*uiView
}
