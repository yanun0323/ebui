package animation

import (
	"math"
	"time"
)

const (
	_defaultAnimationDuration = 300 * time.Millisecond
)

// Style is a curve that describes the progress of an animation over time.
type Style interface {
	// Value returns the value of the animation curve at the given duration.
	//
	// The duration is a value between 0 and 1, where 0 is the start of the animation and 1 is the end.
	Value(duration float64) float64

	// Duration returns the duration of the animation.
	Duration() time.Duration

	// Strength returns a new style with the given strength.
	//
	// The strength is a multiplier for the animation curve.
	//
	// The default strength is 1. More than 1 is stronger, less than 1 is weaker.
	Strengthen(float64) Style

	// Delay returns a new style with the given delay.
	Delay(time.Duration) Style

	// GetDelay returns the delay of the style.
	GetDelay() time.Duration
}

var (
	// None returns a none animation curve.
	None = func(duration ...time.Duration) Style {
		return &stylePrototype{
			formula:  func(float64, float64) float64 { return 1 },
			duration: 0,
		}
	}

	// Linear returns a linear animation curve.
	Linear = NewStyleBuilder(func(t, _ float64) float64 { return t })

	// EaseInOut returns an ease-in-out cubic animation curve.
	EaseInOut = NewStyleBuilder(func(t, s float64) float64 {
		s += 2
		if t < 0.5 {
			return 4 * t * t * t
		} else {
			return 1 - math.Pow(-2*t+2, s)/2
		}
	})

	// EaseIn returns an ease-in animation curve.
	EaseIn = NewStyleBuilder(func(t, s float64) float64 {
		s += 2
		return math.Pow(t, s)
	})

	// EaseOut returns an ease-out animation curve.
	EaseOut = NewStyleBuilder(func(t, s float64) float64 {
		s += 2
		return 1 - math.Pow(1-t, s)
	})

	// Spring returns a spring animation curve.
	Spring = NewStyleBuilder(func(t, s float64) float64 {
		var (
			c1 = s * 1.5
			c2 = c1 * 1.5
		)

		if t < 0.5 {
			return (2 * t * t * ((c2+1)*2*t - c2)) / 2
		} else {
			return (2 * ((t-1)*(t-1)*((c2+1)*(t*2-2)+c2) + 1)) / 2
		}
	})
)

// stylePrototype represents a style that can be used to animate a value over time.
type stylePrototype struct {
	formula  func(t, strength float64) float64
	duration time.Duration
	strength float64
	delay    time.Duration
}

// NewStyleBuilder returns a new style builder.
func NewStyleBuilder(formula func(t, strength float64) float64) func(duration ...time.Duration) Style {
	return func(duration ...time.Duration) Style {
		d := _defaultAnimationDuration
		if len(duration) != 0 {
			d = duration[0]
		}

		return &stylePrototype{
			formula:  formula,
			duration: d,
			strength: 1,
		}
	}
}

func (c *stylePrototype) Value(t float64) float64 {
	return c.formula(max(min(1, t), 0), c.strength)
}

func (c *stylePrototype) Duration() time.Duration {
	return c.duration
}

func (c stylePrototype) Strengthen(strength float64) Style {
	c.strength = strength
	return &c
}

func (c *stylePrototype) Delay(delay time.Duration) Style {
	c.delay = delay
	return c
}

func (c *stylePrototype) GetDelay() time.Duration {
	return c.delay
}
