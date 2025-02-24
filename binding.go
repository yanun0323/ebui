package ebui

func Bind[T comparable](initialValue ...T) *Binding[T] {
	var value T
	if len(initialValue) != 0 {
		value = initialValue[0]
	}

	return BindFunc(
		func() T { return value },
		func(v T) { value = v },
	)
}

func BindFunc[T comparable](get func() T, set func(T)) *Binding[T] {
	return &Binding[T]{
		getter:    get,
		setter:    set,
		listeners: make([]func(), 0),
	}
}

type Binding[T comparable] struct {
	_         noCopy
	getter    func() T
	setter    func(T)
	listeners []func()
}

func (b *Binding[T]) Get() T {
	if b == nil {
		return *new(T)
	}

	return b.getter()
}

func (b *Binding[T]) Set(v T) {
	if b == nil {
		return
	}

	if b.getter() != v {
		b.setter(v)
		b.notifyListeners()
		globalStateManager.markDirty()
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
