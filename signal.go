package ebui

import (
	"reflect"
	"sync"
)

type AnySignal interface {
	ValueAny() any
	SubscribeNotify(func()) func()
}

type Signal[T any] interface {
	AnySignal
	Value() T
}

type State[T any] struct {
	mu       sync.RWMutex
	value    T
	watchers map[uint64]func()
	nextID   uint64
}

func NewState[T any](initial T) *State[T] {
	return &State[T]{
		value:    initial,
		watchers: make(map[uint64]func()),
	}
}

func (s *State[T]) Value() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *State[T]) ValueAny() any {
	return s.Value()
}

func (s *State[T]) Set(next T) {
	s.mu.Lock()
	if reflect.DeepEqual(s.value, next) {
		s.mu.Unlock()
		return
	}
	s.value = next
	watchers := make([]func(), 0, len(s.watchers))
	for _, watcher := range s.watchers {
		watchers = append(watchers, watcher)
	}
	s.mu.Unlock()

	for _, watcher := range watchers {
		watcher()
	}
}

func (s *State[T]) Update(update func(T) T) {
	s.Set(update(s.Value()))
}

func (s *State[T]) SubscribeNotify(fn func()) func() {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++
	s.watchers[id] = fn
	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.watchers, id)
	}
}

type computedSignal[T any] struct {
	mu         sync.RWMutex
	value      T
	compute    func() T
	watchers   map[uint64]func()
	nextID     uint64
	depCancels []func()
}

func Computed[T any](compute func() T, deps ...AnySignal) Signal[T] {
	signal := &computedSignal[T]{
		compute:  compute,
		watchers: make(map[uint64]func()),
	}
	signal.value = compute()
	for _, dep := range deps {
		if dep == nil {
			continue
		}
		signal.depCancels = append(signal.depCancels, dep.SubscribeNotify(signal.recompute))
	}
	return signal
}

func Map[A, B any](source Signal[A], transform func(A) B) Signal[B] {
	if source == nil {
		return Computed(func() B {
			var zero A
			return transform(zero)
		})
	}
	return Computed(func() B {
		return transform(source.Value())
	}, source)
}

func Combine[A, B, C any](left Signal[A], right Signal[B], combine func(A, B) C) Signal[C] {
	deps := make([]AnySignal, 0, 2)
	if left != nil {
		deps = append(deps, left)
	}
	if right != nil {
		deps = append(deps, right)
	}
	return Computed(func() C {
		var a A
		var b B
		if left != nil {
			a = left.Value()
		}
		if right != nil {
			b = right.Value()
		}
		return combine(a, b)
	}, deps...)
}

func (c *computedSignal[T]) Value() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *computedSignal[T]) ValueAny() any {
	return c.Value()
}

func (c *computedSignal[T]) SubscribeNotify(fn func()) func() {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID
	c.nextID++
	c.watchers[id] = fn
	return func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.watchers, id)
	}
}

func (c *computedSignal[T]) recompute() {
	next := c.compute()

	c.mu.Lock()
	if reflect.DeepEqual(c.value, next) {
		c.mu.Unlock()
		return
	}
	c.value = next
	watchers := make([]func(), 0, len(c.watchers))
	for _, watcher := range c.watchers {
		watchers = append(watchers, watcher)
	}
	c.mu.Unlock()

	for _, watcher := range watchers {
		watcher()
	}
}
