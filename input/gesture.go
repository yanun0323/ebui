package input

type GesturePhase int

const (
	GesturePhaseNone GesturePhase = iota
	GesturePhaseBegan
	GesturePhaseMoved
	GesturePhaseEnded
)

type GestureEvent struct {
	Position Vector
	// TODO: Implement me
}
