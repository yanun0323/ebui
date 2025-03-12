package layout

type Direction int8

const (
	Horizontal Direction = 1 << 0
	Vertical   Direction = 1 << 1
)

func (d Direction) Has(other Direction) bool {
	return d&other == other
}
