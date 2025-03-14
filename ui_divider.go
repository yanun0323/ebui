package ebui

type dividerImpl struct {
	*viewCtx

	length *Binding[float64]
}

func Divider(length ...*Binding[float64]) SomeView {
	l := Bind(1.0)
	if len(length) != 0 {
		l = length[0]
	}

	v := &dividerImpl{
		length: l,
	}
	v.viewCtx = newViewContext(v)
	return v
}

func (v *dividerImpl) preload(env *viewCtxEnv, stackTypes ...formulaType) (preloadData, layoutFunc) {
	types := getTypes(stackTypes...)

	length := v.length.Value()
	if v.frameSize == nil {
		v.frameSize = Bind(CGSize{})
	}

	switch types {
	case formulaHStack:
		v.frameSize.Set(NewSize(length, Inf))
	case formulaVStack:
		v.frameSize.Set(NewSize(Inf, length))
	default:
		v.frameSize.Set(CGSize{})
	}

	if v.viewCtx.backgroundColor == nil {
		v.viewCtx.backgroundColor = Bind(NewColor(128))
	}

	_, layout := v.viewCtx.preload(env)
	return newPreloadData(CGSize{}, v.padding(), v.border()), layout
}
