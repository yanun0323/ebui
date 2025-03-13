package ebui

import "time"

var (
	defaultAccentColor = NewColor(0, 0, 255, 255)
)

var (
	AccentColor        = Bind(defaultAccentColor)
	DefaultPadding     = Bind(NewInset(15))
	DefaultBorderColor = Bind(defaultAccentColor)
	DefaultRoundCorner = Bind(8.0)
	DefaultSpacing     = Bind(15.0)
)

var (
	SlowDrawThreshold = 100 * time.Millisecond
)
