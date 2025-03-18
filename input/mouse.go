package input

// MousePhase represents the phase of a touch event.
type MousePhase int

const (
	MousePhaseBegan MousePhase = iota + 1
	MousePhaseMoved
	MousePhaseEnded
	MousePhaseCancelled
)

type MouseEvent struct {
	Phase    MousePhase
	Position Vector
}
