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
		listeners: make([]func(T, T), 0),
	}
}

// BindOneWay binds a source binding to a target binding.
//   - The target binding will be updated with the value of the source binding.
//   - The source binding will not be updated when the target binding is updated.
func BindOneWay[T, F bindable](source *Binding[T], forward func(T) F) *Binding[F] {
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

// BindTwoWay binds a source binding to a target binding.
//   - The target binding will be updated with the value of the source binding.
//   - The source binding will be updated when the target binding is updated.
func BindTwoWay[T, F bindable](source *Binding[T], forward func(T) F, backward func(F) T) *Binding[F] {
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
	}, func(v F) {
		if v == fv {
			return
		}

		fv = v
		sv = backward(v)
		source.Set(sv)
	})
}

func BindCombineOneWay[T bindable](a, b *Binding[T], forward func(a, b T) T) *Binding[T] {
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

func BindCombineTwoWay[T bindable](a, b *Binding[T], forward func(a, b T) T, backward func(T) (a, b T)) *Binding[T] {
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
	}, func(v T) {
		if v == cv {
			return
		}

		cv = v
		av, bv = backward(v)
		a.Set(av)
		b.Set(bv)
	})
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
	listeners []func(T, T)

	animStyle  animation.Style // 設定的預設動畫
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

// AddListener adds a listener to the binding.
//
// The listener will be called when the binding is updated.
//
// AddListener is thread-safe.
func (b *Binding[T]) AddListener(listener func(oldVal, newVal T)) {
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
		// 檢查是否有動畫風格
		var animStyle animation.Style
		if len(with) > 0 && with[0].Duration() > 0 {
			animStyle = with[0]
		}

		if animStyle == nil {
			animStyle = globalContext.CurrentStyle()
		}

		if animStyle == nil {
			animStyle = b.animStyle
		}

		// 如果有動畫風格且持續時間大於0，創建動畫
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

						value := animateValue(oldVal, newVal, float64(elapsed)/float64(duration))
						b.setValue(oldVal, value)
						return false
					},
					onCancel: func() {
						b.setValue(oldVal, newVal)
						b.animResult = nil
					},
				},
			)

			b.animResult = &newVal

			return
		}

		// 無動畫時直接設置值
		b.setValue(oldVal, newVal)
	}
}

func (b *Binding[T]) setValue(oldVal, newVal T) {
	b.setter(newVal)
	b.notifyListeners(oldVal, newVal)
	globalStateManager.markDirty()
}

func (b *Binding[T]) notifyListeners(oldVal, newVal T) {
	if b == nil {
		return
	}

	for _, listener := range b.listeners {
		listener(oldVal, newVal)
	}
}

// Animated sets the animated style as the default animated style for the binding.
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
