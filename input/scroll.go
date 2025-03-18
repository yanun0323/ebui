package input

type ScrollPhase int

const (
	ScrollPhaseBegan ScrollPhase = iota + 1
	ScrollPhaseChanged
	ScrollPhaseEnded
)

type ScrollEvent struct {
	Phase ScrollPhase
	Delta Vector
}
