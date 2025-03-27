package layout

type Direction int8

const (
	DirectionVertical   Direction = 0
	DirectionHorizontal Direction = 1 << 0
)

func (d Direction) Has(other Direction) bool {
	return d&other == other
}
