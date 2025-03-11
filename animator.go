package ebui

import (
	"time"

	"github.com/yanun0323/ebui/animation"
)

type animator[T bindable] struct {
	style      animation.Style
	startValue T
	endValue   T
	startTime  time.Time
}

func newAnimator[T bindable](style animation.Style, startValue, endValue T) animator[T] {
	return animator[T]{
		startTime:  time.Now(),
		style:      style,
		startValue: startValue,
		endValue:   endValue,
	}
}

func (a animator[T]) Value() T {
	if a.style == nil {
		return a.endValue
	}

	elapsed := float64(time.Since(a.startTime).Milliseconds())
	duration := float64(a.style.Duration().Milliseconds())
	progress := elapsed / duration

	return animateValue(a.startValue, a.endValue, progress)
}
