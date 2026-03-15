package input

type ScrollPhase int

const (
	ScrollPhaseNone ScrollPhase = iota
	ScrollPhaseBegan
	ScrollPhaseChanged
	ScrollPhaseEnded
)

type ScrollEvent struct {
	Phase ScrollPhase
	Delta Vector
}
