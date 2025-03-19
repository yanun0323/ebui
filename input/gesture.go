package input

type GesturePhase int

const (
	GesturePhaseBegan GesturePhase = iota + 1
	GesturePhaseMoved
	GesturePhaseEnded
)

type GestureEvent struct {
	// TODO: Implement me
}
