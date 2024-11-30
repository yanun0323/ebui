package ebui

// Binding represents a binding value.
//
// use `ebui.New(val)` or `ebui.NewBinding(get, set)` to create a new Binding.
//
// # example
//
//	b1 := ebui.New("Hello")
//	// or
//	val := "Hello"
//	b2 := ebui.NewBinding(func() string {
//		return val
//	}, func(v string) {
//		val = val
//	})
type Binding[T any] struct {
	getter func() T
	setter func(T)
}

func (b *Binding[T]) Get() T {
	if b.getter != nil {
		return b.getter()
	}

	return *new(T)
}

func (b *Binding[T]) TryGet() (T, bool) {
	if b.getter != nil {
		return b.getter(), true
	}

	return *new(T), false
}

func (b *Binding[T]) Set(val T) {
	if b.setter != nil {
		b.setter(val)
	}
}

func (b *Binding[T]) TrySet(val T) bool {
	if b.setter != nil {
		b.setter(val)
		return true
	}

	return false
}

// New creates a new Binding with the given value.
func New[T any](val ...T) Binding[T] {
	b := Binding[T]{}
	if len(val) != 0 {
		v := val[0]
		b.getter = func() T {
			return v
		}
		b.setter = func(vv T) {
			v = vv
		}
	}
	return b
}

// NewBinding creates a new Binding with the given function.
func NewBinding[T any](get func() T, set func(T)) Binding[T] {
	return Binding[T]{
		getter: get,
		setter: set,
	}
}
