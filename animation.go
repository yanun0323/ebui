package ebui

import "time"

type Animatable interface {
	Animate(duration time.Duration) Animation
}

// 動畫管理器
type AnimationManager struct {
	running []Animation
}

func (am *AnimationManager) Update() {
	now := time.Now()
	for i := len(am.running) - 1; i >= 0; i-- {
		anim := am.running[i]
		if now.Sub(anim.startTime) >= anim.duration {
			// 動畫完成
			am.running = append(am.running[:i], am.running[i+1:]...)
		}
	}
}

func (b Binding[float64]) Animate(to float64, duration time.Duration) Animation {
	return Animation{
		startValue: b.value,
		duration:   duration,
		from:       b.value,
		to:         to,
		startTime:  time.Now(),
	}
}

type AnimationCurve interface {
	Value(progress float64) float64
}

// 預定義的動畫曲線
type (
	LinearCurve struct{}
	EaseInCurve struct{}
	EaseOutCurve struct{}
	EaseInOutCurve struct{}
)

func (LinearCurve) Value(p float64) float64     { return p }
func (EaseInCurve) Value(p float64) float64     { return p * p }
func (EaseOutCurve) Value(p float64) float64    { return -(p * (p - 2)) }
func (EaseInOutCurve) Value(p float64) float64 {
	p *= 2
	if p < 1 {
		return 0.5 * p * p
	}
	p--
	return -0.5 * (p*(p-2) - 1)
}

type Animation struct {
	startValue float64
	endValue   float64
	duration   time.Duration
	startTime  time.Time
	curve      AnimationCurve
	onComplete func()    // 新增：動畫完成回調
}

func NewAnimation(from, to float64, duration time.Duration) *Animation {
	return &Animation{
		startValue: from,
		endValue:   to,
		duration:   duration,
		startTime:  time.Now(),
		curve:      LinearCurve{},
	}
}

func (a *Animation) WithCurve(curve AnimationCurve) *Animation {
	a.curve = curve
	return a
}

func (a *Animation) WithCompletion(f func()) *Animation {
	a.onComplete = f
	return a
}

func (a *Animation) CurrentValue() float64 {
	elapsed := time.Since(a.startTime)
	progress := float64(elapsed) / float64(a.duration)

	if progress >= 1.0 {
		return a.endValue
	}

	// 使用動畫曲線計算當前值
	curveProgress := a.curve.Value(progress)
	return a.startValue + (a.endValue-a.startValue)*curveProgress
}
