package ebui

import (
	"math"
	"sync/atomic"

	"github.com/yanun0323/ebui/input"
)

func removeLastRune(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	return string(runes[:len(runes)-1])
}

func getTypes(types ...stackType) stackType {
	if len(types) == 0 {
		return stackTypeZStack
	}
	return types[0]
}

func ceil(v float64) float64 {
	return math.Ceil(v)
}

func newVector[T numberable](x, y T) input.Vector {
	return input.Vector{X: float64(x), Y: float64(y)}
}

func abs[T numberable](v T) T {
	if v >= 0 {
		return v
	}
	return -v
}

// value is a helper struct for atomic value
type value[T any] struct {
	val atomic.Value
}

func newValue[T any]() *value[T] {
	return &value[T]{
		val: atomic.Value{},
	}
}

func (v *value[T]) Load() T {
	e, ok := v.val.Load().(T)
	if ok {
		return e
	}

	return *new(T)
}

func (v *value[T]) Store(val T) {
	v.val.Store(val)
}

func (v *value[T]) Swap(val T) T {
	e, ok := v.val.Swap(val).(T)
	if ok {
		return e
	}

	return *new(T)
}

func (v *value[T]) IsNil() bool {
	return v.val.Load() == nil
}
