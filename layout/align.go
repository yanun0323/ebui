package layout

import "github.com/yanun0323/ebui/internal/helper"

type Align int8

const (
	// AlignDefault represents top leading alignment
	AlignDefault Align = 0

	// AlignLeading aligns the item to the leading edge of the container
	AlignLeading Align = 1 << 0
	// AlignTop aligns the item to the top edge of the container
	AlignTop Align = 1 << 1
	// AlignTrailing aligns the item to the trailing edge of the container
	AlignTrailing Align = 1 << 2
	// AlignBottom aligns the item to the bottom edge of the container
	AlignBottom Align = 1 << 3
)

const (
	// AlignCenter aligns the item to the center of the container
	AlignCenter Align = AlignLeading | AlignTrailing | AlignTop | AlignBottom
	// AlignCenterHorizontal aligns the item to the horizontal center of the container
	AlignCenterHorizontal Align = AlignLeading | AlignTrailing
	// AlignCenterVertical aligns the item to the vertical center of the container
	AlignCenterVertical Align = AlignTop | AlignBottom
)

func (a Align) Contain(other Align) bool {
	if a == AlignDefault {
		return (AlignTop|AlignLeading)&other == other
	}

	return a&other == other
}

func (a Align) Hash() []byte {
	return helper.BytesInt8(int8(a))
}
