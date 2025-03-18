package ebui

import (
	"math"
	"time"

	"github.com/yanun0323/ebui/input"
)

// TODO: Implement gesture recognizer manager

type gestureType int

const (
	gestureTap gestureType = iota
	gestureDoubleTap
	gestureSwipe
	gesturePinch
	gestureRotation
)

type gestureState int

const (
	gestureStateBegan gestureState = iota
	gestureStateChanged
	gestureStateEnded
	gestureStateCancelled
)

type gestureEvent struct {
	Type      gestureType
	State     gestureState
	Location  input.Vector
	Delta     input.Vector
	Scale     float64 // for pinch gesture
	Rotation  float64 // for rotation gesture
	Velocity  input.Vector
	Timestamp time.Time
}

type gestureRecognizer struct {
	onGesture func(gestureEvent)

	// internal state
	startTime   time.Time
	startPos    input.Vector
	lastPos     input.Vector
	lastTime    time.Time
	touchPoints []input.Vector
	isTracking  bool
}

func newGestureRecognizer(handler func(gestureEvent)) *gestureRecognizer {
	return &gestureRecognizer{
		onGesture: handler,
	}
}

func (gr *gestureRecognizer) HandleTouchEvent(event input.TouchEvent) bool {
	switch event.Phase {
	case input.TouchPhaseBegan:
		gr.startTracking(event)

	case input.TouchPhaseMoved:
		if gr.isTracking {
			gr.updateTracking(event)
		}

	case input.TouchPhaseEnded:
		if gr.isTracking {
			gr.endTracking(event)
		}

	case input.TouchPhaseCancelled:
		gr.cancelTracking()
	}

	return gr.isTracking
}

func (gr *gestureRecognizer) startTracking(event input.TouchEvent) {
	gr.isTracking = true
	gr.startTime = time.Now()
	gr.startPos = event.Position
	gr.lastPos = event.Position
	gr.lastTime = gr.startTime
	gr.touchPoints = []input.Vector{event.Position}
}

func (gr *gestureRecognizer) updateTracking(event input.TouchEvent) {
	now := time.Now()
	delta := newVector(event.Position.X-gr.lastPos.X, event.Position.Y-gr.lastPos.Y)

	// calculate velocity
	duration := now.Sub(gr.lastTime).Seconds()
	velocity := newVector(delta.X/duration, delta.Y/duration)

	gr.onGesture(gestureEvent{
		Type:      gestureSwipe,
		State:     gestureStateChanged,
		Location:  event.Position,
		Delta:     delta,
		Velocity:  velocity,
		Timestamp: now,
	})

	gr.lastPos = event.Position
	gr.lastTime = now
	gr.touchPoints = append(gr.touchPoints, event.Position)
}

func (gr *gestureRecognizer) endTracking(event input.TouchEvent) {
	duration := time.Since(gr.startTime)

	// detect tap
	if duration < 300*time.Millisecond &&
		math.Abs(float64(event.Position.X-gr.startPos.X)) < 10 &&
		math.Abs(float64(event.Position.Y-gr.startPos.Y)) < 10 {
		gr.onGesture(gestureEvent{
			Type:      gestureTap,
			State:     gestureStateEnded,
			Location:  event.Position,
			Timestamp: time.Now(),
		})
	}

	gr.isTracking = false
}

func (gr *gestureRecognizer) cancelTracking() {
	gr.isTracking = false
	gr.touchPoints = nil
}
