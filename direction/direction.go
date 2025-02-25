package direction

type D int8

const (
	Horizontal D = 1 << 0
	Vertical   D = 1 << 1
)

func (d D) Has(other D) bool {
	return d&other == other
}

func (d D) HasNo(other D) bool {
	return d&other == 0
}

func (d D) Include(other D) D {
	return d | other
}

func (d D) Exclude(other D) D {
	return d &^ other
}
