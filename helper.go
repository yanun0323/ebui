package ebui

import (
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
)

func removeLastChar(s string) string {
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

func onHover(bounds CGRect) bool {
	x, y := ebiten.CursorPosition()
	return bounds.Contains(NewPoint(x, y))
}

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
