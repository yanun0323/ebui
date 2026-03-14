package font

import "github.com/yanun0323/ebui/internal/helper"

type Size int

const (
	Caption     Size = 12
	Footnote    Size = 14
	SubHeadline Size = 16
	Body        Size = 18
	Headline    Size = 20
	Title3      Size = 22
	Title2      Size = 26
	Title       Size = 30
	LargeTitle  Size = 34
)

func NewSize(size int) Size {
	if size <= 0 {
		return Body
	}

	return Size(size)
}

func (s Size) Int() int {
	return int(s)
}

func (s Size) F64() float64 {
	return float64(s)
}

func (s Size) Bytes() []byte {
	return helper.BytesInt(int(s))
}
