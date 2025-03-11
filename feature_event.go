package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type eventHandler interface {
	HandleTouchEvent(event touchEvent) bool
	HandleKeyEvent(event keyEvent) bool
}

type touchPhase int

const (
	touchPhaseNone touchPhase = iota
	touchPhaseBegan
	touchPhaseMoved
	touchPhaseEnded
	touchPhaseCancelled
)

type touchEvent struct {
	Phase    touchPhase
	Position CGPoint
}

type keyEvent struct {
	Key     ebiten.Key
	Pressed bool

	Shift   bool
	Control bool
	Alt     bool // Option
	Meta    bool // Windows or Command
}

var globalEventManager = &eventManager{
	handlers:   make([]eventHandler, 0),
	isTracking: false,
}

// 事件管理器
type eventManager struct {
	handlers   []eventHandler
	isTracking bool
}

func (em *eventManager) DispatchTouchEvent(event touchEvent) bool {
	if event.Phase == touchPhaseBegan {
		em.isTracking = true
	} else if event.Phase == touchPhaseEnded || event.Phase == touchPhaseCancelled {
		em.isTracking = false
	}

	for i := len(em.handlers) - 1; i >= 0; i-- {
		if em.handlers[i].HandleTouchEvent(event) {
			return true
		}
	}
	return false
}

func (em *eventManager) DispatchKeyEvent(event keyEvent) bool {
	for i := len(em.handlers) - 1; i >= 0; i-- {
		if em.handlers[i].HandleKeyEvent(event) {
			return true
		}
	}

	return false
}

func (em *eventManager) RegisterHandler(handler eventHandler) {
	em.handlers = append(em.handlers, handler)
}
