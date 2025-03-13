package ebui

import layout "github.com/yanun0323/ebui/layout"

func alignContains(a layout.Align, fn func(containLeading, containTop, containTrailing, containBottom bool)) {
	containLeading := a.Contain(layout.AlignLeading)
	containTop := a.Contain(layout.AlignTop)
	containTrailing := a.Contain(layout.AlignTrailing)
	containBottom := a.Contain(layout.AlignBottom)

	fn(containLeading, containTop, containTrailing, containBottom)
}

func alignToCGPoint(alignment layout.Align) CGPoint {
	p := CGPoint{}
	alignContains(alignment, func(containLeading, containTop, containTrailing, containBottom bool) {
		if containBottom {
			p.Y = 1
		}

		if containTrailing {
			p.X = 1
		}

		if containTop {
			p.Y /= 2
		}

		if containLeading {
			p.X /= 2
		}
	})

	return p
}
