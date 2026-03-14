package ebui

var (
	defaultAccentColor        = NewColor(64, 64, 255, 255)
	defaultIndicatorBaseColor = NewColor(32, 32, 32, 40)
	defaultIndicatorMainColor = NewColor(160, 160, 160, 255)
)

var (
	AccentColor         = Bind(defaultAccentColor)
	DefaultPadding      = Bind(NewInset(15))
	DefaultBorderColor  = Bind(defaultAccentColor)
	DefaultRoundCorner  = Bind(8.0)
	DefaultSpacing      = Bind(15.0)
	DefaultScrollSpeed  = Bind(5.0)
	DefaultShadowColor  = Bind(NewColor(0, 0, 0, 128))
	DefaultShadowLength = Bind(7.0)
)
