package input

type TouchPhase int

const (
	TouchPhaseNone TouchPhase = iota
	TouchPhaseBegan
	TouchPhaseMoved
	TouchPhaseEnded
)

type TouchEvent struct {
	Position Vector
	// TODO: Implement me
}
