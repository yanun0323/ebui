package ebui

func If(condition bool, then SomeView, els SomeView) SomeView {
	if condition {
		return then
	}
	return els
}
