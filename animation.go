package ebui

import "time"

type Animatable interface {
	Animate(duration time.Duration) Animation
}

// 動畫管理器
type AnimationManager struct {
	running []*Animation
}

var defaultAnimationManager = NewAnimationManager()

func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		running: make([]*Animation, 0),
	}
}

type Animation struct {
	startValue float64
	endValue   float64
	duration   time.Duration
	startTime  time.Time
	curve      AnimationCurve
	onComplete func()
}

func (am *AnimationManager) Update() {
	now := time.Now()
	for i := len(am.running) - 1; i >= 0; i-- {
		anim := am.running[i]
		if now.Sub(anim.startTime) >= anim.duration {
			if anim.onComplete != nil {
				anim.onComplete()
			}
			am.running = append(am.running[:i], am.running[i+1:]...)
		}
	}
}

func (am *AnimationManager) Add(anim *Animation) {
	am.running = append(am.running, anim)
}

// 修改為返回指標
func NewAnimation(from, to float64, duration time.Duration) *Animation {
	return &Animation{
		startValue: from,
		endValue:   to,
		duration:   duration,
		startTime:  time.Now(),
		curve:      LinearCurve{},
	}
}

// 修改 Binding 的 animate 方法
func (b *Binding[T]) animate(to T, duration time.Duration) *Animation {
	// 將 T 轉換為 float64
	toFloat64 := func(v T) float64 {
		switch val := any(v).(type) {
		case int:
			return float64(val)
		case int8:
			return float64(val)
		case int16:
			return float64(val)
		case int32:
			return float64(val)
		case int64:
			return float64(val)
		case uint:
			return float64(val)
		case uint8:
			return float64(val)
		case uint16:
			return float64(val)
		case uint32:
			return float64(val)
		case uint64:
			return float64(val)
		case float32:
			return float64(val)
		case float64:
			return val
		default:
			return 0
		}
	}

	return NewAnimation(
		toFloat64(b.Get()),
		toFloat64(to),
		duration,
	)
}

type AnimationCurve interface {
	Value(progress float64) float64
}

// 預定義的動畫曲線
type (
	LinearCurve    struct{}
	EaseInCurve    struct{}
	EaseOutCurve   struct{}
	EaseInOutCurve struct{}
)

func (LinearCurve) Value(p float64) float64  { return p }
func (EaseInCurve) Value(p float64) float64  { return p * p }
func (EaseOutCurve) Value(p float64) float64 { return -(p * (p - 2)) }
func (EaseInOutCurve) Value(p float64) float64 {
	p *= 2
	if p < 1 {
		return 0.5 * p * p
	}
	p--
	return -0.5 * (p*(p-2) - 1)
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

	curveProgress := a.curve.Value(progress)
	return a.startValue + (a.endValue-a.startValue)*curveProgress
}
