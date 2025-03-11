package ebui

import (
	"time"

	"github.com/yanun0323/ebui/animation"
)

// animator 是一個泛型動畫器，負責計算動畫過程中的值
type animator[T bindable] struct {
	style       animation.Style
	startValue  T
	endValue    T
	startTime   time.Time
	completed   bool
	onCompleted func()
}

// newAnimator 創建一個新的動畫器
func newAnimator[T bindable](style animation.Style, startValue, endValue T, onCompleted ...func()) animator[T] {
	var callback func()
	if len(onCompleted) > 0 {
		callback = onCompleted[0]
	}

	return animator[T]{
		startTime:   time.Now(),
		style:       style,
		startValue:  startValue,
		endValue:    endValue,
		completed:   false,
		onCompleted: callback,
	}
}

// Value 返回當前時間點的動畫值
func (a animator[T]) Value() T {
	if a.style == nil {
		return a.endValue
	}

	elapsed := time.Since(a.startTime)
	duration := a.style.Duration()

	// 如果動畫時間已經結束，返回終值
	if elapsed >= duration {
		if !a.completed && a.onCompleted != nil {
			a.onCompleted()
		}
		return a.endValue
	}

	// 計算進度值 (0.0 - 1.0)
	progress := float64(elapsed) / float64(duration)

	// 應用動畫曲線函數
	curveProgress := a.style.Value(progress)

	// 根據曲線進度插值計算當前值
	return animateValue(a.startValue, a.endValue, curveProgress)
}

// IsCompleted 檢查動畫是否已完成
func (a animator[T]) IsCompleted() bool {
	return time.Since(a.startTime) >= a.style.Duration()
}

// DurationLeft 返回剩餘的動畫時間
func (a animator[T]) DurationLeft() time.Duration {
	elapsed := time.Since(a.startTime)
	duration := a.style.Duration()

	if elapsed >= duration {
		return 0
	}

	return duration - elapsed
}
