package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EventHandler interface {
	HandleTouchEvent(event TouchEvent) bool
	HandleKeyEvent(event KeyEvent) bool
}

type TouchPhase int

const (
	TouchPhaseNone TouchPhase = iota
	TouchPhaseBegan
	TouchPhaseMoved
	TouchPhaseEnded
	TouchPhaseCancelled
)

type TouchEvent struct {
	Phase    TouchPhase
	Position CGPoint
}

type KeyEvent struct {
	Key     ebiten.Key
	Pressed bool
}

var globalEventManager = &EventManager{
	handlers:   make([]EventHandler, 0),
	isTracking: false,
}

// 事件管理器
type EventManager struct {
	handlers   []EventHandler
	isTracking bool
}

func (em *EventManager) DispatchTouchEvent(event TouchEvent) bool {
	if event.Phase == TouchPhaseBegan {
		em.isTracking = true
	} else if event.Phase == TouchPhaseEnded || event.Phase == TouchPhaseCancelled {
		em.isTracking = false
	}

	for i := len(em.handlers) - 1; i >= 0; i-- {
		if em.handlers[i].HandleTouchEvent(event) {
			return true
		}
	}
	return false
}

func (em *EventManager) RegisterHandler(handler EventHandler) {
	em.handlers = append(em.handlers, handler)
}

func RegisterEventHandler(handler EventHandler) {
	globalEventManager.RegisterHandler(handler)
}
