package ebui

type identity uint8

const (
	idNone identity = iota
	idEmpty
	idSpacer
	idHStack
	idVStack
	idZStack
	idButton
	idText
	idImage
	idCircle
	idRectangle
)

type paramIdentity uint8

const (
	paramIDFrame paramIdentity = iota
	paramIDOffset
	paramIDForegroundColor
	paramIDBackgroundColor
	paramIDOpacity
)
