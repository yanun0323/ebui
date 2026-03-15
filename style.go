package ebui

type Style struct {
	Padding Insets

	HasFrame       bool
	Frame          Size
	HasBackground  bool
	Background     Color
	HasFill        bool
	Fill           Color
	HasForeground  bool
	Foreground     Color
	HasBorder      bool
	Border         Stroke
	HasOpacity     bool
	Opacity        float64
	HasFont        bool
	Font           TextStyle
	HasLineLimit   bool
	LineLimit      int
	HasAlign       bool
	Align          Alignment
	HasSpacing     bool
	Spacing        float64
	HasShadow      bool
	Shadow         Shadow
	HasDisabled    bool
	Disabled       bool
	HasOffset      bool
	Offset         Point
	HasImageFit    bool
	ImageFit       ImageFit
	PreserveAspect bool
}

type Modifier func(Element) Element

func Apply(view View, modifiers ...Modifier) Element {
	element := Render(view)
	for _, modifier := range modifiers {
		if modifier == nil {
			continue
		}
		element = modifier(element)
	}
	return element
}

func Padding(insets Insets) Modifier {
	return func(element Element) Element {
		element.style.Padding = insets
		return element
	}
}

func Frame(size Size) Modifier {
	return func(element Element) Element {
		element.style.HasFrame = true
		element.style.Frame = size
		return element
	}
}

func Background(color Color) Modifier {
	return func(element Element) Element {
		element.style.HasBackground = true
		element.style.Background = color
		return element
	}
}

func Fill(color Color) Modifier {
	return func(element Element) Element {
		element.style.HasFill = true
		element.style.Fill = color
		return element
	}
}

func Foreground(color Color) Modifier {
	return func(element Element) Element {
		element.style.HasForeground = true
		element.style.Foreground = color
		return element
	}
}

func Border(stroke Stroke) Modifier {
	return func(element Element) Element {
		element.style.HasBorder = true
		element.style.Border = stroke
		return element
	}
}

func Opacity(opacity float64) Modifier {
	return func(element Element) Element {
		element.style.HasOpacity = true
		element.style.Opacity = clampFloat(opacity, 0, 1)
		return element
	}
}

func Font(style TextStyle) Modifier {
	return func(element Element) Element {
		element.style.HasFont = true
		element.style.Font = style
		return element
	}
}

func LineLimit(limit int) Modifier {
	return func(element Element) Element {
		element.style.HasLineLimit = true
		element.style.LineLimit = limit
		return element
	}
}

func Align(alignment Alignment) Modifier {
	return func(element Element) Element {
		element.style.HasAlign = true
		element.style.Align = alignment
		return element
	}
}

func Spacing(spacing float64) Modifier {
	return func(element Element) Element {
		element.style.HasSpacing = true
		element.style.Spacing = spacing
		return element
	}
}

func ShadowStyle(shadow Shadow) Modifier {
	return func(element Element) Element {
		element.style.HasShadow = true
		element.style.Shadow = shadow
		return element
	}
}

func Disabled(disabled bool) Modifier {
	return func(element Element) Element {
		element.style.HasDisabled = true
		element.style.Disabled = disabled
		return element
	}
}

func Offset(offset Point) Modifier {
	return func(element Element) Element {
		element.style.HasOffset = true
		element.style.Offset = offset
		return element
	}
}

func ImageFitMode(fit ImageFit) Modifier {
	return func(element Element) Element {
		element.style.HasImageFit = true
		element.style.ImageFit = fit
		return element
	}
}

func PreserveAspectRatio(enabled bool) Modifier {
	return func(element Element) Element {
		element.style.PreserveAspect = enabled
		return element
	}
}
