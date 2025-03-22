package ebui

import "github.com/yanun0323/ebui/input"

type Hashable interface {
	Hash() []byte
}

type eventHandler interface {
	onScrollEvent(input.Vector, input.ScrollEvent) bool
	onHoverEvent(input.Vector)
	onMouseEvent(input.MouseEvent)
	onKeyEvent(input.KeyEvent)
	onTypeEvent(input.TypeEvent)
	onGestureEvent(input.GestureEvent)
	onTouchEvent(input.TouchEvent)
}
