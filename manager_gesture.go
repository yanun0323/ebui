package ebui

import (
	"math"
	"time"
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
	Location  CGPoint
	Delta     CGPoint
	Scale     float64 // for pinch gesture
	Rotation  float64 // for rotation gesture
	Velocity  CGPoint
	Timestamp time.Time
}

type gestureRecognizer struct {
	onGesture func(gestureEvent)

	// internal state
	startTime   time.Time
	startPos    CGPoint
	lastPos     CGPoint
	lastTime    time.Time
	touchPoints []CGPoint
	isTracking  bool
}

func newGestureRecognizer(handler func(gestureEvent)) *gestureRecognizer {
	return &gestureRecognizer{
		onGesture: handler,
	}
}

func (gr *gestureRecognizer) HandleTouchEvent(event touchEvent) bool {
	switch event.Phase {
	case touchPhaseBegan:
		gr.startTracking(event)

	case touchPhaseMoved:
		if gr.isTracking {
			gr.updateTracking(event)
		}

	case touchPhaseEnded:
		if gr.isTracking {
			gr.endTracking(event)
		}

	case touchPhaseCancelled:
		gr.cancelTracking()
	}

	return gr.isTracking
}

func (gr *gestureRecognizer) startTracking(event touchEvent) {
	gr.isTracking = true
	gr.startTime = time.Now()
	gr.startPos = event.Position
	gr.lastPos = event.Position
	gr.lastTime = gr.startTime
	gr.touchPoints = []CGPoint{event.Position}
}

func (gr *gestureRecognizer) updateTracking(event touchEvent) {
	now := time.Now()
	delta := event.Position.Sub(gr.lastPos)

	// calculate velocity
	duration := now.Sub(gr.lastTime).Seconds()
	velocity := CGPoint{
		X: delta.X / duration,
		Y: delta.Y / duration,
	}

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

func (gr *gestureRecognizer) endTracking(event touchEvent) {
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
