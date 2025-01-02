package ebui

import (
	"image"
	"math"
)

type LayoutConstraints struct {
	MinWidth, MaxWidth   float64
	MinHeight, MaxHeight float64
}

type LayoutContext struct {
	Constraints  LayoutConstraints
	ParentBounds image.Rectangle
}

func (ctx *LayoutContext) FitWithin(size image.Point) image.Rectangle {
	width := float64(size.X)
	height := float64(size.Y)

	width = math.Min(math.Max(width, ctx.Constraints.MinWidth), ctx.Constraints.MaxWidth)
	height = math.Min(math.Max(height, ctx.Constraints.MinHeight), ctx.Constraints.MaxHeight)

	return image.Rect(0, 0, int(width), int(height))
}
