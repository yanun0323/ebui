package ebui

import (
	"time"

	"github.com/yanun0323/ebui/input"
)

type eventHandler interface {
	HandleWheelEvent(input.ScrollEvent)
	HandleTouchEvent(input.TouchEvent)
	HandleKeyEvent(input.KeyEvent)
	HandleInputEvent(input.TypeEvent)
}

var globalEventManager = &eventManager{
	keyPressThreshold:  500,
	keyPressInterval:   50,
	keyStatusUpdatedAt: make(map[input.Key]int64),
	handlers:           make([]eventHandler, 0),
	isTracking:         false,
}

// eventManager responsible for handling events
type eventManager struct {
	keyPressThreshold  int64               /* millisecond */
	keyPressInterval   int64               /* millisecond */
	keyStatusUpdatedAt map[input.Key]int64 /* millisecond */
	handlers           []eventHandler
	isTracking         bool
	lastScrollEvent    input.ScrollEvent
}

func (em *eventManager) DispatchWheelEvent(event input.ScrollEvent) {
	defer func() {
		em.lastScrollEvent = event
	}()

	if event.Delta.IsZero() {
		if !em.lastScrollEvent.Delta.IsZero() {
			return
		}
	}

	for _, handler := range em.handlers {
		handler.HandleWheelEvent(event)
	}
}

func (em *eventManager) DispatchTouchEvent(event input.TouchEvent) {
	if event.Phase == input.TouchPhaseBegan {
		em.isTracking = true
	} else if event.Phase == input.TouchPhaseEnded || event.Phase == input.TouchPhaseCancelled {
		em.isTracking = false
	}

	for _, handler := range em.handlers {
		handler.HandleTouchEvent(event)
	}
}

func (em *eventManager) DispatchKeyEvent(event input.KeyEvent) {
	now := time.Now().UnixMilli()
	switch event.Phase {
	case input.KeyPhaseJustPressed:
		em.keyStatusUpdatedAt[event.Key] = now + em.keyPressThreshold
	case input.KeyPhasePressing:
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

func (em *eventManager) DispatchInputEvent(event input.TypeEvent) {
	for _, handler := range em.handlers {
		handler.HandleInputEvent(event)
	}
}

func (em *eventManager) RegisterHandler(handler eventHandler) {
	em.handlers = append(em.handlers, handler)
}
