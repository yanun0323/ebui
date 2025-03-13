package ebui

import (
	"sync"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
)

// Const creates a binding that always returns the same value.
func Const[T bindable](value T) *Binding[T] {
	return Bind(value)
}

// Bind creates a binding that can be used to bind a value to a UI element.
func Bind[T bindable](initialValue ...T) *Binding[T] {
	var value T
	if len(initialValue) != 0 {
		value = initialValue[0]
	}

	return BindFunc(
		func() T { return value },
		func(v T) { value = v },
	)
}

// BindFunc creates a binding that can be used to bind a value to a UI element.
func BindFunc[T bindable](get func() T, set func(T)) *Binding[T] {
	return &Binding[T]{
		getter:    get,
		setter:    set,
		listeners: make([]func(T, T, animation.Style, bool), 0),
	}
}

// BindForward binds a source binding to a target binding.
//   - The target binding will be updated with the value of the source binding.
//   - The source binding will not be updated when the target binding is updated.
func BindForward[T, F bindable](source *Binding[T], forward func(T) F) *Binding[F] {
	var (
		sv T
		fv F
	)
	return BindFunc(func() F {
		s := source.Get()
		if s == sv {
			return fv
		}

		sv = s
		fv = forward(s)
		return fv
	}, func(v F) {})
}

// BindCombineForward binds two bindings and returns a new binding that combines the two bindings.
//   - The new binding will be updated with the value of the two bindings.
//   - The two bindings will not be updated when the new binding is updated.
func BindCombineForward[T bindable](a, b *Binding[T], forward func(a, b T) T) *Binding[T] {
	var (
		av T
		bv T
		cv T
	)

	return BindFunc(func() T {
		a := a.Get()
		b := b.Get()
		if a == av && b == bv {
			return cv
		}

		av = a
		bv = b
		cv = forward(av, bv)
		return cv
	}, func(v T) {})
}

/*
	######## 		####		##    ##		########
	##     ##		 ## 		###   ##		##     ##
	##     ##		 ## 		####  ##		##     ##
	######## 		 ## 		## ## ##		##     ##
	##     ##		 ## 		##  ####		##     ##
	##     ##		 ## 		##   ###		##     ##
	######## 		####		##    ##		########
*/

// bindable is a type that can be bound to a UI element.
type bindable interface {
	numberable | ~string | ~bool | CGPoint | CGSize | CGRect | CGInset | CGColor | *ebiten.Image
}

// Binding is a binding that can be used to bind a value to a UI element.
//
// Binding is thread-safe.
type Binding[T bindable] struct {
	mu sync.RWMutex

	getter    func() T
	setter    func(T)
	listeners []func(T, T, animation.Style, bool)

	animStyle  animation.Style // set default animation
	animResult *T
}

// Get returns the current value of the binding.
//
// When ignoreAnim is true, it returns the final target value of the binding, ignoring any intermediate values during animation.
//
// Get is thread-safe.
func (b *Binding[T]) Get(ignoreAnim ...bool) T {
	if b == nil {
		return *new(T)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.getValue(len(ignoreAnim) != 0 && ignoreAnim[0])
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

	b.mu.Lock()
	defer b.mu.Unlock()

	b.set(newVal, with...)
}

// Update returns the current value of the binding without considering the animation value offset, and sets the new value.
//
// Update is thread-safe.
func (b *Binding[T]) Update(fn func(oldVal T) (newVal T), with ...animation.Style) {
	if b == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.set(fn(b.getValue(true)), with...)
}

// addListener adds a listener to the binding.
//
// The listener will be called when the binding is updated.
//
// addListener is thread-safe.
func (b *Binding[T]) addListener(listener func(oldVal, newVal T, animStyle animation.Style, isAnimating bool)) {
	if b == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.listeners = append(b.listeners, listener)
}

func (b *Binding[T]) getValue(ignoreAnim bool) T {
	if ignoreAnim && b.animResult != nil {
		return *b.animResult
	}

	return b.getter()
}

func (b *Binding[T]) set(newVal T, with ...animation.Style) {
	oldVal := b.getValue(true)
	if oldVal != newVal {
		// check if there is an animation style
		var animStyle animation.Style
		if len(with) > 0 {
			animStyle = with[0]
		} else {
			animStyle = globalContext.CurrentStyle()

			if animStyle == nil {
				animStyle = b.animStyle
			}
		}

		if animStyle == nil {
			animStyle = animation.None()
		}

		// if there is an animation style and the duration is greater than 0, create an animation
		if animStyle != nil && animStyle.Duration() > 0 {
			id := animationID(unsafe.Pointer(b))
			executor, ok := globalAnimationManager.RemoveExecutor(id)
			if ok {
				executor.onCancel()
			}

			var (
				startTime = time.Now().UnixMilli()
				duration  = animStyle.Duration().Milliseconds()
			)

			globalAnimationManager.AddExecutor(
				id,
				animationExecutor{
					onUpdate: func(now int64) bool {
						elapsed := now - startTime
						if elapsed >= duration {
							return true
						}

						progress := float64(elapsed) / float64(duration)
						progress = animStyle.Value(progress)
						value := animateValue(oldVal, newVal, progress)
						b.setValue(oldVal, value, animStyle, true)
						return value == newVal
					},
					onCancel: func() {
						b.setValue(oldVal, newVal, animStyle, false)
						b.animResult = nil
					},
				},
			)

			b.animResult = &newVal

			return
		}

		// when there is no animation, set the value directly
		b.setValue(oldVal, newVal, animStyle, false)
	}
}

func (b *Binding[T]) setValue(oldVal, newVal T, animStyle animation.Style, isAnimating bool) {
	b.setter(newVal)
	go b.notifyListeners(oldVal, newVal, animStyle, isAnimating)
	globalStateManager.markDirty()
}

func (b *Binding[T]) notifyListeners(oldVal, newVal T, animStyle animation.Style, isAnimating bool) {
	if b == nil {
		return
	}

	for _, listener := range b.listeners {
		listener(oldVal, newVal, animStyle, isAnimating)
	}
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
