package runtime

import ui "github.com/yanun0323/ebui"

type EventKind uint8

const (
	EventPointerMove EventKind = iota
	EventPointerDown
	EventPointerUp
	EventWheel
)

type Event struct {
	Kind   EventKind
	Point  ui.Point
	Scroll ui.Point
}
