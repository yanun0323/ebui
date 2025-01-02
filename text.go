package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"github.com/your-project/ebui/font"
	"golang.org/x/image/font/opentype"
)

type TextStyle struct {
	Size          font.Size
	Weight        font.Weight
	Color         color.Color
	LineHeight    float64
	LetterSpacing float64
	Alignment     TextAlignment
}

type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
)

type textImpl struct {
	content string
	style   TextStyle
	frame   image.Rectangle
	cache   *ViewCache
}

func Text(content string) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &textImpl{
				content: content,
				style: TextStyle{
					Size:          font.Body,
					Weight:        font.Normal,
					Color:         color.Black,
					LineHeight:    1.2,
					LetterSpacing: 0,
					Alignment:     TextAlignLeft,
				},
				cache: NewViewCache(),
			}
		},
	}
}

func (t *textImpl) WithStyle(style TextStyle) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			t.style = style
			return t
		},
	}
}

func (t *textImpl) Layout(bounds image.Rectangle) image.Rectangle {
	face := getFontFace(t.style.Size, t.style.Weight)
	width := font.MeasureString(face, t.content).Round()
	height := t.style.Size.Int()
	
	// 考慮行高
	height = int(float64(height) * t.style.LineHeight)
	
	t.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X + width,
		bounds.Min.Y + height,
	)
	
	return t.frame
}

func (t *textImpl) Draw(screen *ebiten.Image) {
	if t.content == "" {
		return
	}

	face := getFontFace(t.style.Size, t.style.Weight)
	
	// 計算文字位置
	x := t.frame.Min.X
	switch t.style.Alignment {
	case TextAlignCenter:
		width := t.frame.Max.X - t.frame.Min.X
		bounds, _ := font.BoundString(face, t.content)
		textWidth := bounds.Max.X.Round() - bounds.Min.X.Round()
		x += (width - textWidth) / 2
	case TextAlignRight:
		bounds, _ := font.BoundString(face, t.content)
		textWidth := bounds.Max.X.Round() - bounds.Min.X.Round()
		x = t.frame.Max.X - textWidth
	}

	// 繪製文字
	text.Draw(
		screen,
		t.content,
		face,
		x,
		t.frame.Min.Y + int(float64(t.style.Size.Int())*t.style.LineHeight),
		t.style.Color,
	)
}
