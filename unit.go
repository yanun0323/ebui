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

func (p point) Adds(a point) point {
	return point{p.x + a.x, p.y + a.y}
}

func (p point) Subs(a point) point {
	return point{p.x - a.x, p.y - a.y}
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

func (f size) Adds(a size) size {
	return size{f.w + a.w, f.h + a.h}
}

func (f size) Subs(a size) size {
	return size{f.w - a.w, f.h - a.h}
}

type bounds struct {
	top, right, bottom, left int
}

func (b bounds) Add(top, right, bottom, left int) bounds {
	return bounds{b.top + top, b.right + right, b.bottom + bottom, b.left + left}
}

func (b bounds) Sub(top, right, bottom, left int) bounds {
	return bounds{b.top - top, b.right - right, b.bottom - bottom, b.left - left}
}

func (b bounds) Adds(a bounds) bounds {
	return bounds{b.top + a.top, b.right + a.right, b.bottom + a.bottom, b.left + a.left}
}

func (b bounds) Subs(a bounds) bounds {
	return bounds{b.top - a.top, b.right - a.right, b.bottom - a.bottom, b.left - a.left}
}

func (b bounds) IsZero() bool {
	return b.top == 0 && b.right == 0 && b.bottom == 0 && b.left == 0
}
