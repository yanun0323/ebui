package ebui

import (
	"image"

	"github.com/yanun0323/ebui/direction"
)

type numberable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type CGColor struct {
	R, G, B, A uint8
}

func (c CGColor) RGBA() (r, g, b, a uint32) {
	return uint32(c.R) * 256, uint32(c.G) * 256, uint32(c.B) * 256, uint32(c.A) * 256
}

// NewColor uses a traditional 32-bit alpha-premultiplied color, having 8 bits for each of red, green, blue and alpha.
//
// An alpha-premultiplied color component C has been scaled by alpha (A), so has valid values 0 <= C <= A.
func NewColor[Number numberable](r, g, b, a Number) CGColor {
	return CGColor{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// CGPoint represents a coordinate in 2D space.
type CGPoint struct {
	X float64
	Y float64
}

// NewPoint creates a CGPoint from any numberable type.
func NewPoint[Number numberable](x, y Number) CGPoint {
	return CGPoint{X: float64(x), Y: float64(y)}
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

// CGSize represents a size including width and height in 2D space.
type CGSize struct {
	Width  float64
	Height float64
}

// NewSize creates a CGSize from any numberable type.
func NewSize[Number numberable](width, height Number) CGSize {
	return CGSize{Width: max(float64(width), 0), Height: max(float64(height), 0)}
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

// CGRect represents a rectangle including a start point and an end point in 2D space.
type CGRect struct {
	Start CGPoint
	End   CGPoint
}

// NewRect creates a CGRect from any numberable type.
func NewRect[Number numberable](minX, minY, maxX, maxY Number) CGRect {
	return CGRect{
		Start: CGPoint{X: float64(min(minX, maxX)), Y: float64(min(minY, maxY))},
		End:   CGPoint{X: float64(max(minX, maxX)), Y: float64(max(minY, maxY))},
	}
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

// CGInset represents a padding including top, right, bottom and left in 2D space.
type CGInset struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// NewInset creates an Inset from any numberable type.
func NewInset[Number numberable](top, right, bottom, left Number) CGInset {
	return CGInset{Top: float64(top), Right: float64(right), Bottom: float64(bottom), Left: float64(left)}
}

func (i CGInset) IsZero() bool {
	return i == CGInset{}
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

// flexibleSize 表示一個可能包含無限維度的彈性尺寸
type flexibleSize struct {
	Frame    CGSize // 實際的有限尺寸
	IsInfX   bool   // X 軸是否無限
	IsInfY   bool   // Y 軸是否無限
	IsSpacer bool
}

func newFlexibleSize(width, height float64, isSpacer ...bool) flexibleSize {
	isInfX := isInf(width) || width < 0
	isInfY := isInf(height) || height < 0

	if isInfX {
		width = 0
	}

	if isInfY {
		height = 0
	}

	return flexibleSize{
		Frame:    NewSize(width, height),
		IsInfX:   isInfX,
		IsInfY:   isInfY,
		IsSpacer: len(isSpacer) != 0 && isSpacer[0],
	}
}
