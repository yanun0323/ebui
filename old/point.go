package ebui

type Point struct {
	X, Y int
}

func (p *Point) Eq(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

func (p Point) Add(p2 Point) Point {
	return Point{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
	}
}
