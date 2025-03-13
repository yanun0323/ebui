package font

import "github.com/yanun0323/ebui/internal/helper"

type TextAlign int

const (
	TextAlignNone   TextAlign = 0
	TextAlignLeft   TextAlign = 1 << 0
	TextAlignRight  TextAlign = 1 << 1
	TextAlignCenter TextAlign = TextAlignLeft | TextAlignRight
)

func (a TextAlign) Bytes() []byte {
	return helper.BytesInt(int(a))
}
