package animation

import (
	"math"
	"time"
)

const (
	_defaultAnimationDuration = 500 * time.Millisecond
)

// Style is a curve that describes the progress of an animation over time.
type Style interface {
	// Value returns the value of the animation curve at the given duration.
	//
	// The duration is a value between 0 and 1, where 0 is the start of the animation and 1 is the end.
	Value(duration float64) float64

	// Duration returns the duration of the animation.
	Duration() time.Duration
}

// None returns a none animation curve.
func None() Style {
	return &noneAnimationCurve{}
}

type noneAnimationCurve struct{}

func (c *noneAnimationCurve) Value(duration float64) float64 {
	return 1
}

func (c *noneAnimationCurve) Duration() time.Duration {
	return 0
}

// Linear returns a linear animation curve.
//
// The duration is the duration of the animation.
// If no duration is provided, the default duration is used.
func Linear(duration ...time.Duration) Style {
	d := _defaultAnimationDuration
	if len(duration) != 0 {
		d = duration[0]
	}

	return &linearAnimationCurve{
		duration: d,
	}
}

type linearAnimationCurve struct {
	duration time.Duration
}

func (c *linearAnimationCurve) Value(duration float64) float64 {
	return duration
}

func (c *linearAnimationCurve) Duration() time.Duration {
	return c.duration
}

// EaseInOut returns an ease-in-out animation curve.
//
// The duration is the duration of the animation.
// If no duration is provided, the default duration is used.
func EaseInOut(duration ...time.Duration) Style {
	d := _defaultAnimationDuration
	if len(duration) != 0 {
		d = duration[0]
	}

	return &easeInOutAnimationCurve{
		duration: d,
	}
}

type easeInOutAnimationCurve struct {
	duration time.Duration
}

func (c *easeInOutAnimationCurve) Value(duration float64) float64 {
	if duration < 0.5 {
		// 前半段使用 EaseIn (加速)
		return math.Pow(duration*2, 2) / 2
	} else {
		// 後半段使用 EaseOut (減速)
		return 1 - math.Pow(-2*duration+2, 2)/2
	}
}

func (c *easeInOutAnimationCurve) Duration() time.Duration {
	return c.duration
}

// EaseIn returns an ease-in animation curve.
//
// The duration is the duration of the animation.
// If no duration is provided, the default duration is used.
func EaseIn(duration ...time.Duration) Style {
	d := _defaultAnimationDuration
	if len(duration) != 0 {
		d = duration[0]
	}

	return &easeInAnimationCurve{
		duration: d,
	}
}

type easeInAnimationCurve struct {
	duration time.Duration
}

func (c *easeInAnimationCurve) Value(duration float64) float64 {
	return math.Pow(duration, 2)
}

func (c *easeInAnimationCurve) Duration() time.Duration {
	return c.duration
}

// EaseOut returns an ease-out animation curve.
//
// The duration is the duration of the animation.
// If no duration is provided, the default duration is used.
func EaseOut(duration ...time.Duration) Style {
	d := _defaultAnimationDuration
	if len(duration) != 0 {
		d = duration[0]
	}

	return &easeOutAnimationCurve{
		duration: d,
	}
}

type easeOutAnimationCurve struct {
	duration time.Duration
}

func (c *easeOutAnimationCurve) Value(duration float64) float64 {
	return 1 - math.Pow(1-duration, 2)
}

func (c *easeOutAnimationCurve) Duration() time.Duration {
	return c.duration
}
