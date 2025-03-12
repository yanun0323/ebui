package font

type TextAlign int

const (
	TextAlignNone   TextAlign = 0
	TextAlignLeft   TextAlign = 1 << 0
	TextAlignRight  TextAlign = 1 << 1
	TextAlignCenter TextAlign = TextAlignLeft | TextAlignRight
)
