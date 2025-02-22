package ebui

import (
	"math"
	"time"
)

type GestureType int

const (
	GestureTap GestureType = iota
	GestureDoubleTap
	GestureSwipe
	GesturePinch
	GestureRotation
)

type GestureState int

const (
	GestureStateBegan GestureState = iota
	GestureStateChanged
	GestureStateEnded
	GestureStateCancelled
)

type GestureEvent struct {
	Type      GestureType
	State     GestureState
	Location  CGPoint
	Delta     CGPoint
	Scale     float64 // 用於捏合手勢
	Rotation  float64 // 用於旋轉手勢
	Velocity  CGPoint
	Timestamp time.Time
}

type GestureRecognizer struct {
	onGesture func(GestureEvent)

	// 內部狀態
	startTime   time.Time
	startPos    CGPoint
	lastPos     CGPoint
	lastTime    time.Time
	touchPoints []CGPoint
	isTracking  bool
}

func NewGestureRecognizer(handler func(GestureEvent)) *GestureRecognizer {
	return &GestureRecognizer{
		onGesture: handler,
	}
}

func (gr *GestureRecognizer) HandleTouchEvent(event TouchEvent) bool {
	switch event.Phase {
	case TouchPhaseBegan:
		gr.startTracking(event)

	case TouchPhaseMoved:
		if gr.isTracking {
			gr.updateTracking(event)
		}

	case TouchPhaseEnded:
		if gr.isTracking {
			gr.endTracking(event)
		}

	case TouchPhaseCancelled:
		gr.cancelTracking()
	}

	return gr.isTracking
}

func (gr *GestureRecognizer) startTracking(event TouchEvent) {
	gr.isTracking = true
	gr.startTime = time.Now()
	gr.startPos = event.Position
	gr.lastPos = event.Position
	gr.lastTime = gr.startTime
	gr.touchPoints = []CGPoint{event.Position}
}

func (gr *GestureRecognizer) updateTracking(event TouchEvent) {
	now := time.Now()
	delta := event.Position.Sub(gr.lastPos)

	// 計算速度
	duration := now.Sub(gr.lastTime).Seconds()
	velocity := CGPoint{
		X: delta.X / duration,
		Y: delta.Y / duration,
	}

	gr.onGesture(GestureEvent{
		Type:      GestureSwipe,
		State:     GestureStateChanged,
		Location:  event.Position,
		Delta:     delta,
		Velocity:  velocity,
		Timestamp: now,
	})

	gr.lastPos = event.Position
	gr.lastTime = now
	gr.touchPoints = append(gr.touchPoints, event.Position)
}

func (gr *GestureRecognizer) endTracking(event TouchEvent) {
	duration := time.Since(gr.startTime)

	// 檢測點擊
	if duration < 300*time.Millisecond &&
		math.Abs(float64(event.Position.X-gr.startPos.X)) < 10 &&
		math.Abs(float64(event.Position.Y-gr.startPos.Y)) < 10 {
		gr.onGesture(GestureEvent{
			Type:      GestureTap,
			State:     GestureStateEnded,
			Location:  event.Position,
			Timestamp: time.Now(),
		})
	}

	gr.isTracking = false
}

func (gr *GestureRecognizer) cancelTracking() {
	gr.isTracking = false
	gr.touchPoints = nil
}
