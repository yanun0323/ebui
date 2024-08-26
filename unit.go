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

type frame struct {
	w, h int
}

func (f frame) Add(w, h int) frame {
	return frame{f.w + w, f.h + h}
}

func (f frame) Sub(w, h int) frame {
	return frame{f.w - w, f.h - h}
}

type bounds struct {
	top, right, bottom, left int
}
