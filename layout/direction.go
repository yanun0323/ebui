package layout

type Direction int8

const (
	DirectionHorizontal Direction = 1 << 0
	DirectionVertical   Direction = 1 << 1
)

func (d Direction) Has(other Direction) bool {
	return d&other == other
}
