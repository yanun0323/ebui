package input

type Vector struct {
	X float64
	Y float64
}

func (v Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}
