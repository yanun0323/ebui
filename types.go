package ebui

import (
	"image/color"
	"math"
)

const Inf = math.MaxFloat64

type Point struct {
	X float64
	Y float64
}

func P(x, y float64) Point {
	return Point{X: x, Y: y}
}

func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

type Size struct {
	Width  float64
	Height float64
}

func S(width, height float64) Size {
	return Size{Width: width, Height: height}
}

func (s Size) IsZero() bool {
	return s.Width == 0 && s.Height == 0
}

func (s Size) AddInsets(insets Insets) Size {
	return Size{
		Width:  s.Width + insets.Left + insets.Right,
		Height: s.Height + insets.Top + insets.Bottom,
	}
}

type Rect struct {
	Min Point
	Max Point
}

func R(x, y, width, height float64) Rect {
	return Rect{
		Min: P(x, y),
		Max: P(x+width, y+height),
	}
}

func (r Rect) Width() float64 {
	return r.Max.X - r.Min.X
}

func (r Rect) Height() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Size() Size {
	return S(r.Width(), r.Height())
}

func (r Rect) Inset(insets Insets) Rect {
	return Rect{
		Min: P(r.Min.X+insets.Left, r.Min.Y+insets.Top),
		Max: P(r.Max.X-insets.Right, r.Max.Y-insets.Bottom),
	}
}

func (r Rect) Offset(offset Point) Rect {
	return Rect{
		Min: r.Min.Add(offset),
		Max: r.Max.Add(offset),
	}
}

func (r Rect) Contains(point Point) bool {
	return point.X >= r.Min.X && point.X <= r.Max.X && point.Y >= r.Min.Y && point.Y <= r.Max.Y
}

func (r Rect) Intersects(other Rect) bool {
	return r.Min.X < other.Max.X && r.Max.X > other.Min.X && r.Min.Y < other.Max.Y && r.Max.Y > other.Min.Y
}

type Insets struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func All(value float64) Insets {
	return Insets{Top: value, Right: value, Bottom: value, Left: value}
}

func Symmetric(vertical, horizontal float64) Insets {
	return Insets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

func TRBL(top, right, bottom, left float64) Insets {
	return Insets{Top: top, Right: right, Bottom: bottom, Left: left}
}

func (in Insets) Add(other Insets) Insets {
	return Insets{
		Top:    in.Top + other.Top,
		Right:  in.Right + other.Right,
		Bottom: in.Bottom + other.Bottom,
		Left:   in.Left + other.Left,
	}
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func RGBA(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

func (c Color) NRGBA() color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

var (
	White       = RGBA(255, 255, 255, 255)
	Black       = RGBA(0, 0, 0, 255)
	Transparent = RGBA(0, 0, 0, 0)
)

type Constraints struct {
	Min Size
	Max Size
}

func Flexible(max Size) Constraints {
	return Constraints{Max: max}
}

func Fixed(size Size) Constraints {
	return Constraints{Min: size, Max: size}
}

func (c Constraints) Shrink(insets Insets) Constraints {
	width := math.Max(0, c.Max.Width-insets.Left-insets.Right)
	height := math.Max(0, c.Max.Height-insets.Top-insets.Bottom)
	minWidth := math.Max(0, c.Min.Width-insets.Left-insets.Right)
	minHeight := math.Max(0, c.Min.Height-insets.Top-insets.Bottom)
	return Constraints{
		Min: S(minWidth, minHeight),
		Max: S(width, height),
	}
}

func (c Constraints) Constrain(size Size) Size {
	return Size{
		Width:  clampFloat(size.Width, c.Min.Width, c.Max.Width),
		Height: clampFloat(size.Height, c.Min.Height, c.Max.Height),
	}
}

type Alignment uint8

const (
	AlignTopLeading Alignment = iota
	AlignTop
	AlignTopTrailing
	AlignLeading
	AlignCenter
	AlignTrailing
	AlignBottomLeading
	AlignBottom
	AlignBottomTrailing
)

type Direction uint8

const (
	DirectionVertical Direction = iota
	DirectionHorizontal
)

type FontWeight uint16

const (
	WeightLight   FontWeight = 300
	WeightRegular FontWeight = 400
	WeightBold    FontWeight = 700
)

type TextStyle struct {
	Family     string
	Size       float64
	Weight     FontWeight
	LineHeight float64
	Kerning    float64
	Italic     bool
}

func DefaultTextStyle() TextStyle {
	return TextStyle{
		Size:       16,
		Weight:     WeightRegular,
		LineHeight: 4,
	}
}

type Stroke struct {
	Insets Insets
	Color  Color
}

func UniformStroke(width float64, color Color) Stroke {
	return Stroke{
		Insets: All(width),
		Color:  color,
	}
}

type Shadow struct {
	Length float64
	Color  Color
}

type ImageFit uint8

const (
	ImageFitNatural ImageFit = iota
	ImageFitFill
)

type ImageSource struct {
	Path string
}

func ImageFromPath(path string) ImageSource {
	return ImageSource{Path: path}
}

func clampFloat(value, min, max float64) float64 {
	if max == 0 {
		max = Inf
	}
	if value < min {
		return min
	}
	if max != Inf && value > max {
		return max
	}
	return value
}
