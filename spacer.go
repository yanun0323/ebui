package ebui

// Spacer 元件
type spacerImpl struct {
	size float64
}

func Spacer() ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &spacerImpl{size: 1.0}
		},
	}
}
