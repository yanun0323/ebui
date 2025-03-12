package ebui

import layout "github.com/yanun0323/ebui/layout"

func alignContains(a layout.Align, fn func(containLeading, containTop, containTrailing, containBottom bool)) {
	containLeading := a.Contain(layout.AlignLeading)
	containTop := a.Contain(layout.AlignTop)
	containTrailing := a.Contain(layout.AlignTrailing)
	containBottom := a.Contain(layout.AlignBottom)

	fn(containLeading, containTop, containTrailing, containBottom)
}
