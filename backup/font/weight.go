package font

import "github.com/yanun0323/ebui/internal/helper"

type Weight int

const (
	Thin       Weight = 100
	ExtraLight Weight = 200
	Light      Weight = 300
	Normal     Weight = 400
	Medium     Weight = 500
	SemiBold   Weight = 600
	Bold       Weight = 700
	ExtraBold  Weight = 800
	Black      Weight = 900
)

func NewWeight(weight int) Weight {
	if weight <= 0 {
		return Normal
	}

	return Weight(weight)
}

func (w Weight) Int() int {
	return int(w)
}

func (w Weight) F32() float32 {
	return float32(w)
}

func (w Weight) Bytes() []byte {
	return helper.BytesInt(int(w))
}
