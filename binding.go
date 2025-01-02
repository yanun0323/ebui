package ebui

type bindable interface {
}

type Binding[T comparable] struct {
	value     T
	listeners []func()
}

func (b *Binding[T]) Set(v T) {
	if b.value != v { // 值變化時才更新
		b.value = v
		b.notifyListeners()
	}
}

func (b *Binding[T]) notifyListeners() {
	for _, listener := range b.listeners {
		listener()
	}
}
