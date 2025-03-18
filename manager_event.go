package ebui

import (
	"time"

	"github.com/yanun0323/ebui/input"
)

type eventHandler interface {
	onScrollEvent(input.ScrollEvent)
	onMouseEvent(input.MouseEvent)
	onKeyEvent(input.KeyEvent)
	onTypeEvent(input.TypeEvent)
	onGestureEvent(input.GestureEvent)
	onTouchEvent(input.TouchEvent)

	processable() bool
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

func (em *eventManager) Update() {
	removable := make(map[int]bool, len(em.handlers))
	for i, handler := range em.handlers {
		if !handler.processable() {
			removable[i] = true
		}
	}

	if len(removable) == 0 {
		return
	}

	handlers := make([]eventHandler, 0, len(em.handlers)-len(removable))
	for i, handler := range em.handlers {
		if !removable[i] {
			handlers = append(handlers, handler)
		}
	}

	println("remain handlers: ", len(handlers))
	em.handlers = handlers
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
		handler.onScrollEvent(event)
	}
}

func (em *eventManager) DispatchTouchEvent(event input.MouseEvent) {
	if event.Phase == input.MousePhaseBegan {
		em.isTracking = true
	} else if event.Phase == input.MousePhaseEnded || event.Phase == input.MousePhaseCancelled {
		em.isTracking = false
	}

	for _, handler := range em.handlers {
		handler.onMouseEvent(event)
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
		handler.onKeyEvent(event)
	}
}

func (em *eventManager) DispatchInputEvent(event input.TypeEvent) {
	for _, handler := range em.handlers {
		handler.onTypeEvent(event)
	}
}

func (em *eventManager) RegisterHandler(handler eventHandler) {
	em.handlers = append(em.handlers, handler)
}
