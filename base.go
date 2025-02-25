package ebui

import (
	"image"
	"image/color"

	"github.com/yanun0323/ebui/direction"
)

type numberable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Color is a alias of image/color.Color.
type Color color.Color

// CGColor uses a traditional 32-bit alpha-premultiplied color, having 8 bits for each of red, green, blue and alpha.
//
// An alpha-premultiplied color component C has been scaled by alpha (A), so has valid values 0 <= C <= A.
func CGColor(r, g, b, a uint8) Color {
	return Color(color.RGBA{r, g, b, a})
}

// Point represents a coordinate in 2D space.
type Point struct {
	X float64
	Y float64
}

// CGPoint creates a Point from any numberable type.
func CGPoint[Number numberable](x, y Number) Point {
	return Point{X: float64(x), Y: float64(y)}
}

func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p Point) Sub(other Point) Point {
	return Point{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p Point) In(r Rect) bool {
	return p.X >= r.Start.X && p.X <= r.End.X && p.Y >= r.Start.Y && p.Y <= r.End.Y
}

func (p Point) Max(other Point, d direction.D) Point {
	if p.X > other.X && p.Y > other.Y {
		return p
	}

	return other
}

func (p Point) Min(other Point, d direction.D) Point {
	if p.X < other.X && p.Y < other.Y {
		return p
	}

	return other
}

func (p Point) MaxXY(other Point) Point {
	return Point{X: max(p.X, other.X), Y: max(p.Y, other.Y)}
}

func (p Point) MinXY(other Point) Point {
	return Point{X: min(p.X, other.X), Y: min(p.Y, other.Y)}
}

func (p Point) Gt(other Point) bool {
	return p.X > other.X && p.Y > other.Y
}

func (p Point) Lt(other Point) bool {
	return p.X < other.X && p.Y < other.Y
}

// Size represents a size including width and height in 2D space.
type Size struct {
	Width  float64
	Height float64
}

// CGSize creates a Size from any numberable type.
func CGSize[Number numberable](width, height Number) Size {
	return Size{Width: max(float64(width), 0), Height: max(float64(height), 0)}
}

func (s Size) Empty() bool {
	return s.Width == 0 && s.Height == 0
}

func (s Size) Eq(other Size) bool {
	return s.Width == other.Width && s.Height == other.Height
}

func (s Size) Add(other Size) Size {
	return Size{Width: s.Width + other.Width, Height: s.Height + other.Height}
}

func (s Size) Sub(other Size) Size {
	return Size{Width: s.Width - other.Width, Height: s.Height - other.Height}
}

func (s Size) Max(other Size) Size {
	sArea := s.Width * s.Height
	otherArea := other.Width * other.Height
	if sArea > otherArea {
		return s
	}

	return other
}

func (s Size) Min(other Size) Size {
	sArea := s.Width * s.Height
	otherArea := other.Width * other.Height
	if sArea < otherArea {
		return s
	}

	return other
}

func (s Size) MaxWH(other Size) Size {
	return Size{Width: max(s.Width, other.Width), Height: max(s.Height, other.Height)}
}

func (s Size) MinWH(other Size) Size {
	return Size{Width: min(s.Width, other.Width), Height: min(s.Height, other.Height)}
}

func (s Size) ToCGPoint() Point {
	return Point{X: s.Width, Y: s.Height}
}

func (s Size) Expand(inset Inset) Size {
	return Size{Width: s.Width + inset.Left + inset.Right, Height: s.Height + inset.Top + inset.Bottom}
}

func (s Size) Shrink(inset Inset) Size {
	return Size{Width: s.Width - inset.Left - inset.Right, Height: s.Height - inset.Top - inset.Bottom}
}

// Rect represents a rectangle including a start point and an end point in 2D space.
type Rect struct {
	Start Point
	End   Point
}

// CGRect creates a Rect from any numberable type.
func CGRect[Number numberable](minX, minY, maxX, maxY Number) Rect {
	return Rect{
		Start: Point{X: float64(minX), Y: float64(minY)},
		End:   Point{X: float64(max(maxX, minX)), Y: float64(max(maxY, minY))},
	}
}

func (r Rect) Move(offset Point) Rect {
	return CGRect(r.Start.X+offset.X, r.Start.Y+offset.Y, r.End.X+offset.X, r.End.Y+offset.Y)
}

func (r Rect) Empty() bool {
	return r.Start.X == 0 && r.Start.Y == 0 && r.End.X == 0 && r.End.Y == 0
}

func (r Rect) drawable() bool {
	w, h := r.Dx(), r.Dy()
	return int(w) > 0 && int(h) > 0 && !isInf(w) && !isInf(h)
}

func (r Rect) Dx() float64 {
	return max(r.End.X-r.Start.X, 0)
}

func (r Rect) Dy() float64 {
	return max(r.End.Y-r.Start.Y, 0)
}

func (r Rect) MaxStartEnd(other Rect) Rect {
	return Rect{
		Start: r.Start.MaxXY(other.Start),
		End:   r.End.MaxXY(other.End),
	}
}

func (r Rect) MinStartEnd(other Rect) Rect {
	return Rect{
		Start: r.Start.MinXY(other.Start),
		End:   r.End.MinXY(other.End),
	}
}

func (r Rect) Size() Size {
	return Size{
		Width:  r.Dx(),
		Height: r.Dy(),
	}
}

func (r Rect) Rect() image.Rectangle {
	return image.Rect(int(r.Start.X), int(r.Start.Y), int(r.End.X), int(r.End.Y))
}

func (r Rect) Expand(inset Inset) Rect {
	return CGRect(r.Start.X, r.Start.Y, r.End.X+inset.Left+inset.Right, r.End.Y+inset.Top+inset.Bottom)
}

// Inset represents a padding including top, right, bottom and left in 2D space.
type Inset struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// CGInset creates an Inset from any numberable type.
func CGInset[Number numberable](top, right, bottom, left Number) Inset {
	return Inset{Top: float64(top), Right: float64(right), Bottom: float64(bottom), Left: float64(left)}
}

func (i Inset) MaxBounds(other Inset) Inset {
	return Inset{
		Top:    max(i.Top, other.Top),
		Right:  max(i.Right, other.Right),
		Bottom: max(i.Bottom, other.Bottom),
		Left:   max(i.Left, other.Left),
	}
}

// flexibleSize 表示一個可能包含無限維度的彈性尺寸
type flexibleSize struct {
	Frame    Size // 實際的有限尺寸
	IsInfX   bool // X 軸是否無限
	IsInfY   bool // Y 軸是否無限
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
		Frame:    CGSize(width, height),
		IsInfX:   isInfX,
		IsInfY:   isInfY,
		IsSpacer: len(isSpacer) != 0 && isSpacer[0],
	}
}
