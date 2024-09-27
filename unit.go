package ebui

type point struct {
	x, y int
}

func (p point) Add(x, y int) point {
	return point{p.x + x, p.y + y}
}

func (p point) Sub(x, y int) point {
	return point{p.x - x, p.y - y}
}

var _zeroSize = size{-1, -1}

type size struct {
	w, h int
}

func (f size) IsZero() bool {
	return f.w == -1 && f.h == -1
}

func (f size) Add(w, h int) size {
	return size{f.w + w, f.h + h}
}

func (f size) Sub(w, h int) size {
	return size{f.w - w, f.h - h}
}

type bounds struct {
	top, right, bottom, left int
}

func (b bounds) Add(top, right, bottom, left int) bounds {
	return bounds{b.top + top, b.right + right, b.bottom + bottom, b.left + left}
}

func (b bounds) IsZero() bool {
	return b.top == 0 && b.right == 0 && b.bottom == 0 && b.left == 0
}
