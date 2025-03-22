package ebui

import (
	"github.com/yanun0323/ebui/input"
	"github.com/yanun0323/ebui/layout"
)

const (
	_sliderSizeIndicator       = 20.0
	_sliderSizeBar             = 5.0
	_sliderSizeIndicatorRadius = _sliderSizeIndicator / 2
)

type sliderImpl struct {
	*stackImpl

	value      *Binding[float64]
	max        *Binding[float64]
	min        *Binding[float64]
	leftOffset *Binding[CGPoint]
}

func Slider(value, maximum, minimum *Binding[float64]) SomeView {
	isDragging := Bind(false)
	hovered := Bind(false)

	minValue := minimum.Get()
	maxValue := maximum.Get()
	currValue := value.Get()
	ratio := (currValue - minValue) / (maxValue - minValue)

	colorBlockScale := Bind(NewPoint(ratio, 1))
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
			Fill(Const(ivory)).Offset(leftOffset),
	)

	hs.Align(Const(layout.AlignLeading | layout.AlignTop | layout.AlignBottom)).
		Fill(Const(NewColor(128))).
		Spacing(Const(0.0))

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
		value.Set(minValue + valueRatio*(maxValue-minValue))

		offset.X -= _sliderSizeIndicatorRadius
		leftOffset.Set(NewPoint(offset.X, 0))
		colorBlockScale.Set(NewPoint(ratio, 1))
	})

	return &sliderImpl{
		stackImpl:  hs,
		value:      value,
		max:        maximum,
		min:        minimum,
		leftOffset: leftOffset,
	}
}

func (s *sliderImpl) preload(ctx *viewCtxEnv, types ...stackType) (preloadData, layoutFunc) {
	logf("preload")
	return preloadData{}, func(start CGPoint, childBoundsSize CGSize) (CGRect, alignFunc) {
		return CGRect{}, func(offset CGPoint) {}
	}
	data, layout := s.stackImpl.preload(ctx, types...)
	return data, func(start CGPoint, childBoundsSize CGSize) (CGRect, alignFunc) {
		bounds, alignFunc := layout(start, childBoundsSize)
		logf("bounds: %v", bounds)
		return bounds, alignFunc
	}
}
