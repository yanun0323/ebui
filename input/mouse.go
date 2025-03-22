package input

// MousePhase represents the phase of a touch event.
type MousePhase int

const (
	MousePhaseNone MousePhase = iota
	MousePhaseBegan
	MousePhaseMoved
	MousePhaseEnded
	MousePhaseCancelled
)

type MouseEvent struct {
	Phase    MousePhase
	Position Vector
}
