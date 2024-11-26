package ebui

/* Check Interface Implementation */
var _ SomeView = (*spacerView)(nil)

func Spacer() *spacerView {
	v := &spacerView{}
	v.uiView = newView(typesSpacer, v)
	return v
}

type spacerView struct {
	*uiView
}

func (p *spacerView) getSize() size {
	return _zeroSize
}

func (p *spacerView) Frame(w, h int) SomeView {
	return p
}
