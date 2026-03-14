package layout

import "github.com/yanun0323/ebui/internal/helper"

type Align int8

const (
	AlignDefault  Align = 0      // AlignDefault represents top leading alignment
	AlignTop      Align = 1 << 0 // AlignTop aligns the item to the top edge of the container
	AlignTrailing Align = 1 << 1 // AlignTrailing aligns the item to the trailing edge of the container
	AlignBottom   Align = 1 << 2 // AlignBottom aligns the item to the bottom edge of the container
	AlignLeading  Align = 1 << 3 // AlignLeading aligns the item to the leading edge of the container
)

const (
	// AlignCenter aligns the item to the center of the container
	AlignCenter         Align = AlignLeading | AlignTrailing | AlignTop | AlignBottom
	AlignTopCenter      Align = AlignTop | AlignLeading | AlignTrailing    // AlignTopCenter aligns the item to the top center of the container
	AlignTrailingCenter Align = AlignTrailing | AlignTop | AlignBottom     // AlignTrailingCenter aligns the item to the trailing center of the container
	AlignBottomCenter   Align = AlignBottom | AlignLeading | AlignTrailing // AlignBottomCenter aligns the item to the bottom center of the container
	AlignLeadingCenter  Align = AlignLeading | AlignTop | AlignBottom      // AlignLeadingCenter aligns the item to the leading center of the container
)

const (
	AlignTopLeading     Align = AlignTop | AlignLeading     // AlignTopLeading aligns the item to the top leading edge of the container
	AlignTopTrailing    Align = AlignTop | AlignTrailing    // AlignTopTrailing aligns the item to the top trailing edge of the container
	AlignBottomLeading  Align = AlignBottom | AlignLeading  // AlignBottomLeading aligns the item to the bottom leading edge of the container
	AlignBottomTrailing Align = AlignBottom | AlignTrailing // AlignBottomTrailing aligns the item to the bottom trailing edge of the container
)

func (a Align) Vector() (float64, float64) {
	x, y := 0.0, 0.0
	if a.Contain(AlignBottom) {
		y = 1
	}

	if a.Contain(AlignTrailing) {
		x = 1
	}

	if a.Contain(AlignTop) {
		y /= 2
	}

	if a.Contain(AlignLeading) {
		x /= 2
	}

	return x, y
}

func (a Align) Contain(other Align) bool {
	if a == AlignDefault {
		return AlignTopLeading&other == other
	}

	return a&other == other
}

func (a Align) Hash() []byte {
	return helper.BytesInt8(int8(a))
}

func (a Align) String() string {
	switch a {
	case AlignDefault:
		return "TopLeading"
	case AlignTop:
		return "Top"
	case AlignTrailing:
		return "Trailing"
	case AlignBottom:
		return "Bottom"
	case AlignLeading:
		return "Leading"
	case AlignCenter:
		return "Center"
	case AlignTopCenter:
		return "TopCenter"
	case AlignTrailingCenter:
		return "TrailingCenter"
	case AlignBottomCenter:
		return "BottomCenter"
	case AlignLeadingCenter:
		return "LeadingCenter"
	case AlignTopLeading:
		return "TopLeading"
	case AlignTopTrailing:
		return "TopTrailing"
	case AlignBottomLeading:
		return "BottomLeading"
	case AlignBottomTrailing:
		return "BottomTrailing"
	default:
		return "Unknown Align"
	}
}
