package ebui

type Binding[T comparable] struct {
	value     T
	listeners []func()
}

func NewBinding[T comparable](initialValue T) Binding[T] {
	return Binding[T]{
		value:     initialValue,
		listeners: make([]func(), 0),
	}
}

func (b *Binding[T]) Get() T {
	return b.value
}

func (b *Binding[T]) Set(v T) {
	if b.value != v {
		b.value = v
		b.notifyListeners()
		defaultStateManager.MarkDirty()
	}
}

func (b *Binding[T]) AddListener(listener func()) {
	b.listeners = append(b.listeners, listener)
}

func (b *Binding[T]) notifyListeners() {
	for _, listener := range b.listeners {
		listener()
	}
}
