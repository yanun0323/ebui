package ebui

import (
	"time"

	"github.com/yanun0323/ebui/input"
	"github.com/yanun0323/ebui/layout"
)

const (
	_sliderSizeIndicator       = 22.0
	_sliderSizeBar             = 5.0
	_sliderSizeIndicatorRadius = _sliderSizeIndicator / 2
	_sliderIndicatorShadow     = 3.0
)

type sliderImpl struct {
	*stackImpl

	value      *Binding[float64]
	max        *Binding[float64]
	min        *Binding[float64]
	leftOffset *Binding[CGPoint]
}

func Slider(value, minimum, maximum *Binding[float64]) SomeView {
	isDragging := Bind(false)
	hovered := Bind(false)

	colorBlockScale := Bind(NewPoint(0, 1))
	leftOffset := Bind(NewPoint(0, 0))
	hs := stack(stackTypeZStack, false,
		HStack(
			Rectangle().
				Frame(Const(NewSize(Inf, _sliderSizeBar))).
				Fill(AccentColor).
				Scale(colorBlockScale),
		).Frame(Const(NewSize(Inf, _sliderSizeBar))).
			Fill(Const(NewColor(128))).
			RoundCorner(),
		Circle().
			Frame(Const(NewSize(_sliderSizeIndicator))).
			Fill(Const(white)).Offset(leftOffset).
			Shadow(Const(_sliderIndicatorShadow)),
	)

	hs.OnAppear(func() {
		var (
			minValue       = minimum.Get()
			maxValue       = maximum.Get()
			currValue      = value.Get()
			ratio          = (currValue - minValue) / (maxValue - minValue)
			frameSizeWidth = hs.systemSetFrame().Size().Width - _sliderSizeIndicatorRadius
			anim           = value.animStyle
		)

		if anim != nil {
			anim = anim.Delay(100 * time.Millisecond)
		}

		colorBlockScale.Set(NewPoint(ratio, 1), anim)
		leftOffset.Set(NewPoint(ratio*frameSizeWidth, 0), anim)
	})

	hs.OnHover(func(b bool) {
		hovered.Set(b)
	})

	hs.OnMouse(func(phase input.MousePhase, offset input.Vector) {
		switch phase {
		case input.MousePhaseNone:
			return
		case input.MousePhaseBegan:
			if !hovered.Get() {
				return
			}

			isDragging.Set(true)
		case input.MousePhaseMoved:
			if !isDragging.Get() {
				return
			}
		case input.MousePhaseEnded:
			isDragging.Set(false)
			if !isDragging.Get() {
				return
			}
		}

		frame := hs.systemSetFrame()
		offset.X = min(offset.X, frame.Dx()-_sliderSizeIndicatorRadius)
		offset.X = max(offset.X, _sliderSizeIndicatorRadius)
		ratio := offset.X / (frame.Dx() - _sliderSizeIndicatorRadius)
		valueRatio := (offset.X - _sliderSizeIndicatorRadius) / (frame.Dx() - _sliderSizeIndicator)
		minValue := minimum.Value()
		maxValue := maximum.Value()
		value.Set(minValue+valueRatio*(maxValue-minValue), nil)

		offset.X -= _sliderSizeIndicatorRadius
		leftOffset.Set(NewPoint(offset.X, 0))
		colorBlockScale.Set(NewPoint(ratio, 1))
	})

	return (&sliderImpl{
		stackImpl:  hs,
		value:      value,
		max:        maximum,
		min:        minimum,
		leftOffset: leftOffset,
	}).Spacing(Const(0.0)).
		Padding(Const(NewInset(7))).
		Border(Const(NewInset(1)), Const(transparent)).
		Align(Const(layout.AlignLeading | layout.AlignTop | layout.AlignBottom))
}
