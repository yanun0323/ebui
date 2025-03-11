package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
)

// noCopy may be added to structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
//
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// bindable is a type that can be bound to a UI element.
type bindable interface {
	numberable | ~string | ~bool | CGPoint | CGSize | CGRect | CGInset | CGColor | *ebiten.Image
}

// Const creates a binding that always returns the same value.
func Const[T bindable](value T) *Binding[T] {
	return &Binding[T]{
		getter: func() T { return value },
		setter: func(T) {},
	}
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

// Binding is a binding that can be used to bind a value to a UI element.
type Binding[T bindable] struct {
	_ noCopy

	getter    func() T
	setter    func(T)
	listeners []func(T, T)
}

// Get returns the value of the binding.
func (b *Binding[T]) Get() T {
	if b == nil {
		return *new(T)
	}

	return b.getter()
}

// Set sets the value of the binding.
//
// The binding will notify its listeners when the value is updated.
func (b *Binding[T]) Set(newVal T, with ...animation.Style) {
	if b == nil {
		return
	}

	if b.getter() != newVal {
		oldVal := b.getter()

		// 檢查是否有動畫風格
		var animStyle animation.Style
		if len(with) > 0 && with[0].Duration() > 0 {
			animStyle = with[0]
		} else {
			// 檢查當前動畫上下文
			animStyle = animation.GetCurrentStyle()
		}

		// 如果有動畫風格且持續時間大於0，創建動畫
		if animStyle != nil && animStyle.Duration() > 0 {
			globalAnimationManager.CreateAnimatedExecutor(
				animStyle,
				func(progress float64) bool {
					// 計算插值
					currentVal := animateValue(oldVal, newVal, progress)
					b.setter(currentVal)
					b.notifyListeners(oldVal, newVal)
					return false // 繼續動畫，直到完成
				},
			)
			return
		}

		// 無動畫時直接設置值
		b.setter(newVal)
		b.notifyListeners(oldVal, newVal)
		globalStateManager.markDirty()
	}
}

// AddListener adds a listener to the binding.
//
// The listener will be called when the binding is updated.
func (b *Binding[T]) AddListener(listener func(oldVal, newVal T)) {
	if b == nil {
		return
	}

	b.listeners = append(b.listeners, listener)
}

func (b *Binding[T]) notifyListeners(oldVal, newVal T) {
	if b == nil {
		return
	}

	for _, listener := range b.listeners {
		listener(oldVal, newVal)
	}
}
