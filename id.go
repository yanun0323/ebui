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
