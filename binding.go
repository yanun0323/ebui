package ebui

func NewBinding[T comparable](initialValue ...T) *Binding[T] {
	var value T
	if len(initialValue) != 0 {
		value = initialValue[0]
	}

	return &Binding[T]{
		value:     value,
		listeners: make([]func(), 0),
	}
}

type Binding[T comparable] struct {
	_         noCopy
	value     T
	listeners []func()
}

func (b *Binding[T]) Get() T {
	if b == nil {
		return *new(T)
	}

	return b.value
}

func (b *Binding[T]) Set(v T) {
	if b.value != v {
		b.value = v
		b.notifyListeners()
		defaultStateManager.markDirty()
	}
}

func (b *Binding[T]) addListener(listener func()) {
	if b == nil {
		return
	}

	b.listeners = append(b.listeners, listener)
}

func (b *Binding[T]) notifyListeners() {
	if b == nil {
		return
	}

	for _, listener := range b.listeners {
		listener()
	}
}
