package ebui

// Size represents a size.
type Size struct {
	W, H int
}

func (s *Size) Reset() {
	s.W = 0
	s.H = 0
}

func (s Size) IsZero() bool {
	return s.W == 0 && s.H == 0
}

func (s Size) Add(s2 Size) Size {
	return Size{
		W: s.W + s2.W,
		H: s.H + s2.H,
	}
}

func (s Size) Eq(s2 Size) bool {
	return s.W == s2.W && s.H == s2.H
}
