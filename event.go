package ebui

import (
	"image"

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
	Position image.Point
}

type KeyEvent struct {
	Key     ebiten.Key
	Pressed bool
}

// 事件管理器
type EventManager struct {
	handlers []EventHandler
}

func (em *EventManager) DispatchTouchEvent(event TouchEvent) bool {
	// 從上到下傳遞事件
	for i := len(em.handlers) - 1; i >= 0; i-- {
		if em.handlers[i].HandleTouchEvent(event) {
			return true
		}
	}
	return false
}
