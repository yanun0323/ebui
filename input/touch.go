package input

type TouchPhase int

const (
	TouchPhaseBegan TouchPhase = iota + 1
	TouchPhaseMoved
	TouchPhaseEnded
)

type TouchEvent struct {
	// TODO: Implement me
}
