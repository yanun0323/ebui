package ebui

import (
	"math"
)

type cornerHandler struct {
	width, height int
	border        float64
	cRadius       float64
	bRadius       float64
	corners       []struct {
		condition        func(x, y int) bool
		CenterX, CenterY int
	}
}

func newCornerHandler(width, height int, radius float64, borderLength ...float64) *cornerHandler {
	border := 0.0
	if len(borderLength) != 0 {
		border = math.Abs(borderLength[0])
	}

	cornerCenter := int(radius + border)
	if cornerCenter > width/2 || cornerCenter > height/2 {
		cornerCenter = min(width/2, height/2)
		radius = float64(cornerCenter) - border
	}
	yCornerCenter := cornerCenter
	xCornerCenter := cornerCenter
	if xCornerCenter > width/2 {
		xCornerCenter = width / 2
	}

	if yCornerCenter > height/2 {
		yCornerCenter = height / 2
	}

	return &cornerHandler{
		width:   width,
		height:  height,
		border:  border,
		bRadius: radius,
		cRadius: radius + border,
		corners: []struct {
			condition        func(x, y int) bool
			CenterX, CenterY int
		}{
			{func(x, y int) bool { return x <= xCornerCenter && y <= yCornerCenter }, xCornerCenter, yCornerCenter},                               /* left top */
			{func(x, y int) bool { return x <= xCornerCenter && y >= height-yCornerCenter }, xCornerCenter, height - yCornerCenter},               /* left bottom */
			{func(x, y int) bool { return x >= width-xCornerCenter && y <= yCornerCenter }, width - xCornerCenter, yCornerCenter},                 /* right top */
			{func(x, y int) bool { return x >= width-xCornerCenter && y >= height-yCornerCenter }, width - xCornerCenter, height - yCornerCenter}, /* right bottom */
		},
	}
}

func (h *cornerHandler) Execute(fn func(isOutside, isBorder bool, x, y int)) {
	cRadiusSqr := int(h.cRadius * h.cRadius)
	bRadiusSqr := int(h.bRadius * h.bRadius)

	border := int(h.border)
	leftBounds := border
	rightBounds := h.width - border
	topBounds := border
	bottomBounds := h.height - border

	handleCorner := func(x, y int) bool {
		for _, corner := range h.corners {
			if corner.condition(x, y) {
				dx := x - corner.CenterX
				dy := y - corner.CenterY
				distance := dx*dx + dy*dy

				if distance >= cRadiusSqr {
					fn(true, false, x, y)
					return true
				}

				if distance >= bRadiusSqr {
					fn(false, true, x, y)
					return true
				}

				fn(false, false, x, y)
				return true
			}
		}

		return false
	}

	for x := range h.width {
		for y := range h.height {
			if handleCorner(x, y) {
				continue
			}

			if h.border == 0 {
				continue
			}

			between := x <= leftBounds || x >= rightBounds || y <= topBounds || y >= bottomBounds
			fn(false, between, x, y)
		}
	}
}
