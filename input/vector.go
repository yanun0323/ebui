package input

type Vector struct {
	X float64
	Y float64
}

func (v Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vector) Add(x, y float64) Vector {
	return Vector{X: v.X + x, Y: v.Y + y}
}

func (v Vector) Sub(x, y float64) Vector {
	return Vector{X: v.X - x, Y: v.Y - y}
}

func (v Vector) Reverse() Vector {
	return Vector{X: -v.X, Y: -v.Y}
}
