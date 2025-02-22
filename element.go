package ebui

import (
	"image"

	"github.com/yanun0323/ebui/direction"
)

var (
	ptZero = pt(0, 0)
)

// CGPoint 是 Core Graphics 的 Point
type CGPoint struct {
	X float64
	Y float64
}

func pt(x, y float64) CGPoint {
	return CGPoint{X: x, Y: y}
}

func (p CGPoint) Add(other CGPoint) CGPoint {
	return CGPoint{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p CGPoint) Sub(other CGPoint) CGPoint {
	return CGPoint{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p CGPoint) In(r CGRect) bool {
	return p.X >= r.Start.X && p.X < r.End.X && p.Y >= r.Start.Y && p.Y < r.End.Y
}

func (p CGPoint) Max(other CGPoint, d direction.D) CGPoint {
	if p.X > other.X && p.Y > other.Y {
		return p
	}

	return other
}

func (p CGPoint) Min(other CGPoint, d direction.D) CGPoint {
	if p.X < other.X && p.Y < other.Y {
		return p
	}

	return other
}

func (p CGPoint) MaxXY(other CGPoint) CGPoint {
	return CGPoint{X: max(p.X, other.X), Y: max(p.Y, other.Y)}
}

func (p CGPoint) MinXY(other CGPoint) CGPoint {
	return CGPoint{X: min(p.X, other.X), Y: min(p.Y, other.Y)}
}

func (p CGPoint) Gt(other CGPoint) bool {
	return p.X > other.X && p.Y > other.Y
}

func (p CGPoint) Lt(other CGPoint) bool {
	return p.X < other.X && p.Y < other.Y
}

// CGSize 是 Core Graphics 的 Size
type CGSize struct {
	Width  float64
	Height float64
}

func sz(width, height float64) CGSize {
	return CGSize{Width: max(width, 0), Height: max(height, 0)}
}

func (s CGSize) Empty() bool {
	return s.Width == 0 && s.Height == 0
}

func (s CGSize) Eq(other CGSize) bool {
	return s.Width == other.Width && s.Height == other.Height
}

func (s CGSize) Add(other CGSize) CGSize {
	return CGSize{Width: s.Width + other.Width, Height: s.Height + other.Height}
}

func (s CGSize) Sub(other CGSize) CGSize {
	return CGSize{Width: s.Width - other.Width, Height: s.Height - other.Height}
}

func (s CGSize) Max(other CGSize) CGSize {
	sArea := s.Width * s.Height
	otherArea := other.Width * other.Height
	if sArea > otherArea {
		return s
	}

	return other
}

func (s CGSize) Min(other CGSize) CGSize {
	sArea := s.Width * s.Height
	otherArea := other.Width * other.Height
	if sArea < otherArea {
		return s
	}

	return other
}

func (s CGSize) MaxWH(other CGSize) CGSize {
	return CGSize{Width: max(s.Width, other.Width), Height: max(s.Height, other.Height)}
}

func (s CGSize) MinWH(other CGSize) CGSize {
	return CGSize{Width: min(s.Width, other.Width), Height: min(s.Height, other.Height)}
}

func (s CGSize) ToCGPoint() CGPoint {
	return CGPoint{X: s.Width, Y: s.Height}
}

func (s CGSize) Expand(inset Inset) CGSize {
	return CGSize{Width: s.Width + inset.Left + inset.Right, Height: s.Height + inset.Top + inset.Bottom}
}

func (s CGSize) Shrink(inset Inset) CGSize {
	return CGSize{Width: s.Width - inset.Left - inset.Right, Height: s.Height - inset.Top - inset.Bottom}
}

// CGRect 是 Core Graphics 的 Rectangle
type CGRect struct {
	Start CGPoint
	End   CGPoint
}

func rect(minX, minY, maxX, maxY float64) CGRect {
	return CGRect{
		Start: CGPoint{X: minX, Y: minY},
		End:   CGPoint{X: max(maxX, minX), Y: max(maxY, minY)},
	}
}

func (r CGRect) Move(offset CGPoint) CGRect {
	return rect(r.Start.X+offset.X, r.Start.Y+offset.Y, r.End.X+offset.X, r.End.Y+offset.Y)
}

func (r CGRect) Empty() bool {
	return r.Start.X == 0 && r.Start.Y == 0 && r.End.X == 0 && r.End.Y == 0
}

func (r CGRect) drawable() bool {
	return !isInf(r.End.X) && !isInf(r.End.Y) && r.Dx() > 0 && r.Dy() > 0
}

func (r CGRect) Dx() float64 {
	return max(r.End.X-r.Start.X, 0)
}

func (r CGRect) Dy() float64 {
	return max(r.End.Y-r.Start.Y, 0)
}

func (r CGRect) MaxStartEnd(other CGRect) CGRect {
	return CGRect{
		Start: r.Start.MaxXY(other.Start),
		End:   r.End.MaxXY(other.End),
	}
}

func (r CGRect) MinStartEnd(other CGRect) CGRect {
	return CGRect{
		Start: r.Start.MinXY(other.Start),
		End:   r.End.MinXY(other.End),
	}
}

func (r CGRect) Size() CGSize {
	return CGSize{
		Width:  r.Dx(),
		Height: r.Dy(),
	}
}

func (r CGRect) Rect() image.Rectangle {
	return image.Rect(int(r.Start.X), int(r.Start.Y), int(r.End.X), int(r.End.Y))
}

func (r CGRect) Shrink(inset Inset) CGRect {
	return rect(r.Start.X+inset.Left, r.Start.Y+inset.Top, r.End.X-inset.Right, r.End.Y-inset.Bottom)
}

func (r CGRect) Expand(inset Inset) CGRect {
	return rect(r.Start.X, r.Start.Y, r.End.X+inset.Left+inset.Right, r.End.Y+inset.Top+inset.Bottom)
}

type Inset struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func ins(top, right, bottom, left float64) Inset {
	return Inset{Top: top, Right: right, Bottom: bottom, Left: left}
}

func (i Inset) MaxBounds(other Inset) Inset {
	return Inset{
		Top:    max(i.Top, other.Top),
		Right:  max(i.Right, other.Right),
		Bottom: max(i.Bottom, other.Bottom),
		Left:   max(i.Left, other.Left),
	}
}
