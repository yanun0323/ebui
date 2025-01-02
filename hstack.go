package ebui

// HStack 水平排列
func HStack(views ...SomeView) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			children := make([]View, len(views))
			for i, v := range views {
				children[i] = v.Build()
			}
			return &hStackImpl{children: children}
		},
	}
}
