package ebui

/* Check Interface Implementation */
var _ SomeView = (*vstackView)(nil)

func VStack(views ...View) *vstackView {
	v := &vstackView{}
	v.uiView = newView(typesVStack, v, views...)
	return v
}

type vstackView struct {
	*uiView
}
