package ebui

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

func Const[T comparable](value T) *Binding[T] {
	return &Binding[T]{
		getter: func() T { return value },
		setter: func(T) {},
	}
}

type constant[T comparable] struct{ value T }

func (c constant[T]) Get() T           { return c.value }
func (constant[T]) Set(T)              {}
func (constant[T]) AddListener(func()) {}

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

func (b *Binding[T]) AddListener(listener func()) {
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

func (b *Binding[T]) Combine(other *Binding[T], combine func(T, T) T) *Binding[T] {
	if b == nil {
		return other
	}

	if other == nil {
		return b
	}

	return BindFunc(func() T {
		return combine(b.Get(), other.Get())
	}, func(v T) {
		// do nothing
	})
}
