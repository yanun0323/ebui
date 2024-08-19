package font

type Size int

const (
	Caption    Size = 12
	Body       Size = 24
	Headline   Size = 28
	Title3     Size = 32
	Title2     Size = 40
	Title      Size = 48
	LargeTitle Size = 56
)

func NewSize(size int) Size {
	if size <= 0 {
		return Body
	}

	return Size(size)
}
