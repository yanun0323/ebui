package ebui

import "github.com/yanun0323/ebui/input"

type Hashable interface {
	Hash() []byte
}

type eventHandler interface {
	onAppearEvent()
	onHoverEvent(input.Vector)
	onScrollEvent(input.Vector, input.ScrollEvent) bool
	onMouseEvent(input.MouseEvent)
	onKeyEvent(input.KeyEvent)
	onTypeEvent(input.TypeEvent)
	onGestureEvent(input.GestureEvent)
	onTouchEvent(input.TouchEvent)
}
