package input

// TouchPhase represents the phase of a touch event.
type TouchPhase int

const (
	TouchPhaseBegan TouchPhase = iota + 1
	TouchPhaseMoved
	TouchPhaseEnded
	TouchPhaseCancelled
)

type TouchEvent struct {
	Phase    TouchPhase
	Position Vector
}
