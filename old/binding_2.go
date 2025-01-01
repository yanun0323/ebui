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
	val *T
}

func (b *Binding[T]) Get() T {
	if b.val != nil {
		return *b.val
	}

	return *new(T)
}

func (b *Binding[T]) TryGet() (T, bool) {
	if b.val != nil {
		return *b.val, true
	}

	return *new(T), false
}

func (b *Binding[T]) Set(val T) {
	if b.val != nil {
		*b.val = val
	}
}

func (b *Binding[T]) TrySet(val T) bool {
	if b.val != nil {
		*b.val = val
		return true
	}

	return false
}

// New creates a new Binding with the given value.
func New[T any](val ...T) Binding[T] {
	b := Binding[T]{}
	if len(val) != 0 {
		b.val = &val[0]
	}

	return b
}

// NewBinding creates a new Binding with the given function.
// func NewBinding[T any](get func() T, set func(T)) Binding[T] {
// 	return Binding[T]{
// 		getter: get,
// 		setter: set,
// 	}
// }
