package ebui

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
)

// Const creates a binding that always returns the same value.
func Const[T bindable](value T) *Binding[T] {
	return BindFunc(
		func() T { return value },
		nil,
	)
}

// Bind creates a binding that can be used to bind a value to a UI element.
func Bind[T bindable](initialValue ...T) *Binding[T] {
	var value T
	if len(initialValue) != 0 {
		value = initialValue[0]
	}

	return BindFunc(
		func() T { return value },
		func(v T, with ...animation.Style) { value = v },
	)
}

// BindFunc creates a binding that can be used to bind a value to a UI element.
func BindFunc[T bindable](get func() T, set func(T, ...animation.Style)) *Binding[T] {
	return &Binding[T]{
		getter:    get,
		setter:    set,
		listeners: make([]func(T, T, ...animation.Style), 0),
	}
}

// BindOneWay creates a one-way binding from source to target.
//
// When proactive is true, target updates immediately on source changes.
// Otherwise, target updates only when its value is accessed.
func BindOneWay[T, F bindable](source *Binding[T], forward func(T) F) *Binding[F] {
	return BindTwoWay(source, forward, nil)
}

func BindTwoWay[T, F bindable](source *Binding[T], forward func(T) F, backward func(F) T) *Binding[F] {
	var (
		sv T
		fv F
	)

	sv = source.Get()
	fv = forward(sv)
	b := BindFunc(func() F {
		return fv
	}, func(v F, with ...animation.Style) {
		fv = v
		if backward != nil {
			ssv := backward(v)
			if ssv != sv {
				sv = ssv
				source.Set(sv, with...)
			}
		}
	})

	source.AddListener(func(oldVal, newVal T, animStyle ...animation.Style) {
		if newVal != sv {
			sv = newVal
			b.Set(forward(sv), animStyle...)
		}
	})

	return b
}

// bindCombineOneWay creates a one-way binding from two sources to a target.
//
// When proactive is true, target updates immediately on source changes.
// Otherwise, target updates only when its value is accessed.
func bindCombineOneWay[T bindable](a, b *Binding[T], forward func(a, b T) T) *Binding[T] {
	var (
		av T
		bv T
	)

	av = a.Get()
	bv = b.Get()
	c := Bind(forward(av, bv))
	a.AddListener(func(oldVal, newVal T, animStyle ...animation.Style) {
		if newVal != av {
			av = newVal
			c.Set(forward(av, bv), animStyle...)
		}
	})

	b.AddListener(func(oldVal, newVal T, animStyle ...animation.Style) {
		if newVal != bv {
			bv = newVal
			c.Set(forward(av, bv), animStyle...)
		}
	})

	return c
}

/*
	######## 	####	##    ##	########
	##     ##	 ## 	###   ##	##     ##
	##     ##	 ## 	####  ##	##     ##
	######## 	 ## 	## ## ##	##     ##
	##     ##	 ## 	##  ####	##     ##
	##     ##	 ## 	##   ###	##     ##
	######## 	####	##    ##	########
*/

// bindable is a type that can be bound to a UI element.
// type bindable comparable
type bindable interface {
	numberable | ~string | ~bool | CGPoint | CGSize | CGRect | CGInset | CGColor | *ebiten.Image
}

// Binding is a binding that can be used to bind a value to a UI element.
//
// Binding is thread-safe.
type Binding[T bindable] struct {
	mu sync.RWMutex

	getter    func() T
	setter    func(T, ...animation.Style)
	listeners []func(T, T, ...animation.Style)
	notifying atomic.Bool

	animStyle  animation.Style // set default animation
	animResult atomic.Pointer[T]
}

func (b *Binding[T]) id() animationID {
	return animationID(unsafe.Pointer(b))
}

// Description returns the description of the binding.
func (b *Binding[T]) Description() string {
	if b == nil {
		return "<nil>"
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	type Info struct {
		IsConstant        bool            `json:"isConstant"`
		CurrentValue      T               `json:"currentValue"`
		CurrentFinalValue T               `json:"currentFinalValue"`
		ListenerCount     int             `json:"listenerCount"`
		AnimStyle         animation.Style `json:"animStyle"`
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	info, err := json.MarshalIndent(&Info{
		IsConstant:        b.getter == nil,
		CurrentValue:      b.getValue(false),
		CurrentFinalValue: b.getValue(true),
		ListenerCount:     len(b.listeners),
		AnimStyle:         b.animStyle,
	}, "", "  ")
	if err != nil {
		return "<unknown>"
	}

	return string(info)
}

// Value returns the current value of the binding. Current value may be an intermediate value during animation.
//
// Value is thread-safe.
func (b *Binding[T]) Value() T {
	if b == nil {
		return *new(T)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.getValue(false)
}

// Get returns the final value of the binding. Final value is the value after animation is completed.
//
// Get is thread-safe.
func (b *Binding[T]) Get() T {
	if b == nil {
		return *new(T)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.getValue(true)
}

// Set sets the value of the binding.
//
// The binding will notify its listeners when the value is updated.
//
// Set is thread-safe.
func (b *Binding[T]) Set(newVal T, with ...animation.Style) {
	if b == nil {
		return
	}

	if b.setter == nil {
		return
	}

	if b.notifying.Load() {
		return
	}

	executor, ok := globalAnimationManager.RemoveExecutor(b.id())
	if ok {
		executor.onCancel()
	}

	b.set(newVal, false, with...)
}

// AddListener adds a listener to the binding.
//
// The listener will be called when the binding is updated.
//
// AddListener is thread-safe.
func (b *Binding[T]) AddListener(listener func(oldVal, newVal T, animStyle ...animation.Style)) *Binding[T] {
	if b == nil {
		return b
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.listeners = append(b.listeners, listener)

	return b
}

// Animated sets the animated style as the default animated style for the binding.
//
// If not provided any animation style, it will use the default animation style.
//
// Animated is thread-safe.
func (b *Binding[T]) Animated(with ...animation.Style) *Binding[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(with) == 0 {
		b.animStyle = animation.EaseInOut()
	} else {
		b.animStyle = with[0]
	}

	return b
}

func (b *Binding[T]) getValue(ignoreAnim bool) T {
	if ignoreAnim {
		if val := b.animResult.Load(); val != nil {
			return *val
		}
	}

	return b.getter()
}

func (b *Binding[T]) getAnimationStyle(with ...animation.Style) animation.Style {
	if len(with) > 0 {
		return with[0]
	}

	if s := globalContext.CurrentStyle(); s != nil {
		return s
	}

	return b.animStyle
}

func (b *Binding[T]) set(newVal T, rmAnimResult bool, with ...animation.Style) {
	if b.setter == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.setValue(newVal, with...)
	if rmAnimResult {
		b.animResult.Store(nil)
	}
}

func (b *Binding[T]) setValue(newVal T, with ...animation.Style) {
	if b.setter == nil {
		return
	}

	var (
		animStyle = b.getAnimationStyle(with...)
		oldVal    = b.getValue(true)
	)

	if oldVal == newVal {
		return
	}

	animatableValue := animStyle != nil && animStyle.Duration() > 0 && animatable(newVal)
	if animatableValue {
		var (
			startTime = time.Now().Add(animStyle.GetDelay()).UnixMilli()
			duration  = animStyle.Duration().Milliseconds()
		)

		cancel, ok := globalAnimationManager.RemoveExecutor(b.id())
		if ok {
			cancel.onCancel()
		}

		globalAnimationManager.AddExecutor(
			b.id(),
			animationExecutor{
				onUpdate: func(now int64) bool {
					elapsed := now - startTime
					if elapsed >= duration {
						b.set(newVal, true, nil)
						return true
					}

					progress := float64(elapsed) / float64(duration)
					progress = animStyle.Value(progress)
					value := animateValue(oldVal, newVal, progress)
					b.set(value, false, nil)
					return false
				},
				onCancel: func() {
					b.set(newVal, true, nil)
				},
			},
		)

		b.animResult.Store(&newVal)

		return
	}

	setNewVal := func(newVal T, animStyle animation.Style) {
		if animStyle == nil {
			animStyle = animation.None()
			b.setter(newVal, with...)
		} else {
			b.setter(newVal, animStyle)
		}

		globalStateManager.markDirty()
		go func() { // notifyListeners
			b.notifying.Store(true)
			defer b.notifying.Store(false)

			for _, listener := range b.listeners {
				listener(oldVal, newVal, animStyle)
			}
		}()
	}

	// when there is no animation, set the value directly
	setNewVal(newVal, animStyle)
}
