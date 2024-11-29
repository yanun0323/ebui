package ebui

import "sync/atomic"

type value[T any] struct {
	val atomic.Value
}

func NewValue[T any](val ...T) value[T] {
	v := atomic.Value{}
	if len(val) != 0 {
		v.Store(val[0])
	}

	return value[T]{
		val: v,
	}
}

// CompareAndSwap executes the compare-and-swap operation for the Value.
func (v *value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	return v.val.CompareAndSwap(old, new)
}

// Load returns the value set by the most recent Store.
// It returns zero value if there has been no call to Store for this Value.
func (v *value[T]) Load() T {
	if vv, ok := v.val.Load().(T); ok {
		return vv
	}

	var zero T
	return zero
}

func (v *value[T]) TryLoad() (T, bool) {
	raw := v.val.Load()
	if raw == nil {
		return *new(T), false
	}

	vv, ok := raw.(T)
	return vv, ok
}

// Store sets the value of the Value v to val.
func (v *value[T]) Store(val T) {
	v.val.Store(val)
}

// Swap stores new into Value and returns the previous value. It returns zero value if
// the Value is empty.
func (v *value[T]) Swap(new T) (old T) {
	if vv, ok := v.val.Swap(new).(T); ok {
		return vv
	}

	var zero T
	return zero
}
