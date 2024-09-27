package ebui

/* Check Interface Implementation */
var _ SomeView = (*hstackView)(nil)

func HStack(views ...View) *hstackView {
	v := &hstackView{}
	v.uiView = newView(typesHStack, v, views...)
	return v
}

type hstackView struct {
	*uiView
}
