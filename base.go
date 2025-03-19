package ebui

import (
	"image"
	"unsafe"

	"github.com/yanun0323/ebui/input"
)

type numberable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// NewColor uses a traditional 32-bit alpha-premultiplied color, having 8 bits for each of red, green, blue and alpha.
//
// An alpha-premultiplied color component C has been scaled by alpha (A), so has valid values 0 <= C <= A.
//
// # Usage:
//
//	NewColor()  /* transparent color */
//	NewColor(gray)  /* color with gray */
//	NewColor(r, g, b)  /* color with alpha 255 */
//	NewColor(r, g, b, a)  /* color with alpha */
func NewColor[Number numberable](val ...Number) CGColor {
	switch len(val) {
	case 1:
		return CGColor{uint8(val[0]), uint8(val[0]), uint8(val[0]), 255}
	case 3:
		return CGColor{uint8(val[0]), uint8(val[1]), uint8(val[2]), 255}
	case 4:
		return CGColor{uint8(val[0]), uint8(val[1]), uint8(val[2]), uint8(val[3])}
	default:
		return CGColor{}
	}
}

// CGColor represents a color with red, green, blue and alpha components.
type CGColor struct {
	R, G, B, A uint8
}

func (c CGColor) RGBA() (r, g, b, a uint32) {
	return uint32(c.R) * 256, uint32(c.G) * 256, uint32(c.B) * 256, uint32(c.A) * 256
}

func (c CGColor) Bytes() []byte {
	return (*[4]byte)(unsafe.Pointer(&c))[:]
}

func (c CGColor) IsZero() bool {
	return c == transparent
}

// NewPoint creates a CGPoint from any numberable type.
//
// # Usage:
//
//	CGPoint{} /* zero value */
//	NewPoint(v) /* x = v, y = v */
//	NewPoint(x, y) /* x, y */
func NewPoint[Number numberable](val ...Number) CGPoint {
	switch len(val) {
	case 1:
		return CGPoint{X: float64(val[0]), Y: float64(val[0])}
	case 2:
		return CGPoint{X: float64(val[0]), Y: float64(val[1])}
	default:
		return CGPoint{}
	}
}

// CGPoint represents a coordinate in 2D space.
type CGPoint struct {
	X float64
	Y float64
}

func (p CGPoint) Bytes() []byte {
	return (*[16]byte)(unsafe.Pointer(&p))[:]
}

var zeroPoint CGPoint

func (p CGPoint) IsZero() bool {
	return p == zeroPoint
}

func (p CGPoint) Add(other CGPoint) CGPoint {
	return CGPoint{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p CGPoint) Sub(other CGPoint) CGPoint {
	return CGPoint{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p CGPoint) In(r CGRect) bool {
	return p.X >= r.Start.X && p.X <= r.End.X && p.Y >= r.Start.Y && p.Y <= r.End.Y
}

func (p CGPoint) Max(other CGPoint) CGPoint {
	if p.X > other.X && p.Y > other.Y {
		return p
	}

	return other
}

func (p CGPoint) Min(other CGPoint) CGPoint {
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

func (p CGPoint) vector() input.Vector {
	return newVector(p.X, p.Y)
}

// NewSize creates a CGSize from any numberable type.
//
// # Usage:
//
//	CGSize{} /* zero value */
//	NewSize(size) /* width, height = size */
//	NewSize(width, height) /* width, height */
func NewSize[Number numberable](val ...Number) CGSize {
	switch len(val) {
	case 1:
		return CGSize{Width: float64(val[0]), Height: float64(val[0])}
	case 2:
		return CGSize{Width: float64(val[0]), Height: float64(val[1])}
	default:
		return CGSize{}
	}
}

// CGSize represents a size including width and height in 2D space.
type CGSize struct {
	Width  float64
	Height float64
}

func (s CGSize) Bytes() []byte {
	return (*[16]byte)(unsafe.Pointer(&s))[:]
}

var zeroSize CGSize

func (s CGSize) drawable() bool {
	return s.Width > 0 && s.Height > 0 && !isInf(s.Width) && !isInf(s.Height)
}

func (s CGSize) IsZero() bool {
	return s == zeroSize
}

func (s CGSize) Empty() bool {
	return s.Width == 0 && s.Height == 0
}

func (s CGSize) IsInfWidth() bool {
	return isInf(s.Width) || s.Width < 0
}

func (s CGSize) IsInfHeight() bool {
	return isInf(s.Height) || s.Height < 0
}

func (s CGSize) NoneInfSize() CGSize {
	if s.IsInfWidth() {
		s.Width = 0
	}

	if s.IsInfHeight() {
		s.Height = 0
	}

	return s
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

func (s CGSize) Expand(inset CGInset) CGSize {
	return CGSize{Width: s.Width + inset.Left + inset.Right, Height: s.Height + inset.Top + inset.Bottom}
}

func (s CGSize) Shrink(inset CGInset) CGSize {
	return CGSize{Width: s.Width - inset.Left - inset.Right, Height: s.Height - inset.Top - inset.Bottom}
}

// NewRect creates a CGRect from any numberable type.
//
// # Usage:
//
//	CGRect{} /* zero value */
//	NewRect(xy) /* start = xy, end = xy */
//	NewRect(minXY, maxXY) /* start = minXY, end = maxXY */
//	NewRect(minX, minY, maxX, maxY) /* minX, minY, maxX, maxY */
func NewRect[Number numberable](val ...Number) CGRect {
	switch len(val) {
	case 1:
		return CGRect{
			Start: CGPoint{X: float64(val[0]), Y: float64(val[0])},
			End:   CGPoint{X: float64(val[0]), Y: float64(val[0])},
		}
	case 2:
		return CGRect{
			Start: CGPoint{X: float64(val[0]), Y: float64(val[0])},
			End:   CGPoint{X: float64(val[1]), Y: float64(val[1])},
		}
	case 4:
		return CGRect{
			Start: CGPoint{X: float64(min(val[0], val[2])), Y: float64(min(val[1], val[3]))},
			End:   CGPoint{X: float64(max(val[0], val[2])), Y: float64(max(val[1], val[3]))},
		}
	default:
		return CGRect{}
	}
}

// CGRect represents a rectangle including a start point and an end point in 2D space.
type CGRect struct {
	Start CGPoint
	End   CGPoint
}

func (r CGRect) Bytes() []byte {
	return (*[32]byte)(unsafe.Pointer(&r))[:]
}

var zeroRect CGRect

func (r CGRect) IsZero() bool {
	return r == zeroRect
}

func (r CGRect) Contains(p input.Vector) bool {
	return p.X >= r.Start.X && p.X <= r.End.X && p.Y >= r.Start.Y && p.Y <= r.End.Y
}

func (r CGRect) Contain(p CGPoint) bool {
	return p.X >= r.Start.X && p.X <= r.End.X && p.Y >= r.Start.Y && p.Y <= r.End.Y
}

func (r CGRect) Move(offset CGPoint) CGRect {
	return NewRect(r.Start.X+offset.X, r.Start.Y+offset.Y, r.End.X+offset.X, r.End.Y+offset.Y)
}

func (r CGRect) Empty() bool {
	return r.Start.X == 0 && r.Start.Y == 0 && r.End.X == 0 && r.End.Y == 0
}

func (r CGRect) drawable() bool {
	w, h := r.Dx(), r.Dy()
	return int(w) > 0 && int(h) > 0 && !isInf(w) && !isInf(h)
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

func (r CGRect) Expand(inset CGInset) CGRect {
	return NewRect(r.Start.X, r.Start.Y, r.End.X+inset.Left+inset.Right, r.End.Y+inset.Top+inset.Bottom)
}

// NewInset creates an Inset from any numberable type.
//
// # Usage:
//
//	CGInset{} /* zero value */
//	NewInset(all)
//	NewInset(horizontal, vertical)
//	NewInset(top, right, bottom, left)
func NewInset[Number numberable](inset ...Number) CGInset {
	switch len(inset) {
	case 1:
		return CGInset{Top: float64(inset[0]), Right: float64(inset[0]), Bottom: float64(inset[0]), Left: float64(inset[0])}
	case 2:
		return CGInset{Top: float64(inset[0]), Right: float64(inset[1]), Bottom: float64(inset[0]), Left: float64(inset[1])}
	case 4:
		return CGInset{Top: float64(inset[0]), Right: float64(inset[1]), Bottom: float64(inset[2]), Left: float64(inset[3])}
	default:
		return CGInset{}
	}
}

// CGInset represents a padding including top, right, bottom and left in 2D space.
type CGInset struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func (i CGInset) Bytes() []byte {
	return (*[32]byte)(unsafe.Pointer(&i))[:]
}

var zeroInset CGInset

func (i CGInset) IsZero() bool {
	return i == zeroInset
}

func (i CGInset) Add(other CGInset) CGInset {
	return CGInset{
		Top:    i.Top + other.Top,
		Right:  i.Right + other.Right,
		Bottom: i.Bottom + other.Bottom,
		Left:   i.Left + other.Left,
	}
}
func (i CGInset) MaxBounds(other CGInset) CGInset {
	return CGInset{
		Top:    max(i.Top, other.Top),
		Right:  max(i.Right, other.Right),
		Bottom: max(i.Bottom, other.Bottom),
		Left:   max(i.Left, other.Left),
	}
}
