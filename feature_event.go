package ebui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type eventHandler interface {
	HandleTouchEvent(event touchEvent)
	HandleKeyEvent(event keyEvent)
	HandleInputEvent(event inputEvent)
}

// touchEvent presents the touch event of the user
type touchEvent struct {
	Phase    touchPhase
	Position CGPoint
}

type touchPhase int

const (
	touchPhaseNone touchPhase = iota
	touchPhaseBegan
	touchPhaseMoved
	touchPhaseEnded
	touchPhaseCancelled
)

// keyEvent presents the key event of the user
type keyEvent struct {
	Key   ebiten.Key
	Phase keyPhase

	Shift   bool
	Control bool
	Alt     bool // Option
	Meta    bool // Windows or Command
}

type keyPhase int

const (
	keyPhaseJustPressed keyPhase = iota
	keyPhasePressing
	keyPhaseJustReleased
)

// inputEvent presents the input event of the user
type inputEvent struct {
	Char rune
}

var globalEventManager = &eventManager{
	keyPressThreshold:  500,
	keyPressInterval:   50,
	keyStatusUpdatedAt: make(map[ebiten.Key]int64),
	handlers:           make([]eventHandler, 0),
	isTracking:         false,
}

// eventManager responsible for handling events
type eventManager struct {
	keyPressThreshold  int64                /* millisecond */
	keyPressInterval   int64                /* millisecond */
	keyStatusUpdatedAt map[ebiten.Key]int64 /* millisecond */
	handlers           []eventHandler
	isTracking         bool
}

func (em *eventManager) DispatchTouchEvent(event touchEvent) {
	if event.Phase == touchPhaseBegan {
		em.isTracking = true
	} else if event.Phase == touchPhaseEnded || event.Phase == touchPhaseCancelled {
		em.isTracking = false
	}

	for _, handler := range em.handlers {
		handler.HandleTouchEvent(event)
	}
}

func (em *eventManager) DispatchKeyEvent(event keyEvent) {
	now := time.Now().UnixMilli()
	switch event.Phase {
	case keyPhaseJustPressed:
		em.keyStatusUpdatedAt[event.Key] = now + em.keyPressThreshold
	case keyPhasePressing:
		t, ok := em.keyStatusUpdatedAt[event.Key]
		if ok {
			if now <= t {
				return
			}
			em.keyStatusUpdatedAt[event.Key] = now + em.keyPressInterval
		}
	}

	for _, handler := range em.handlers {
		handler.HandleKeyEvent(event)
	}
}

func (em *eventManager) DispatchInputEvent(event inputEvent) {
	for _, handler := range em.handlers {
		handler.HandleInputEvent(event)
	}
}

func (em *eventManager) RegisterHandler(handler eventHandler) {
	em.handlers = append(em.handlers, handler)
}
