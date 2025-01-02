package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// 字體相關常數
const (
	FontSizeBody     = 16.0
	FontSizeTitle    = 24.0
	FontWeightNormal = 400.0
	FontWeightBold   = 700.0

	defaultFontWeight  = FontWeightNormal
	defaultLineSpacing = 1.2
	defaultKerning     = 1
)

type TextStyle struct {
	Size          float64
	Weight        float32
	Color         color.Color
	LineHeight    float64
	LetterSpacing float64
	Alignment     TextAlignment
	Italic        bool
}

type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
)

type textImpl struct {
	content       string
	style         TextStyle
	frame         image.Rectangle
	cache         *ViewCache
	face          text.Face
	contentGetter func() string
}

func Text(content string) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &textImpl{
				content: content,
				style: TextStyle{
					Size:          FontSizeBody,
					Weight:        FontWeightNormal,
					Color:         color.Black,
					LineHeight:    defaultLineSpacing,
					LetterSpacing: defaultKerning,
					Alignment:     TextAlignLeft,
					Italic:        false,
				},
				cache: NewViewCache(),
				face:  createDefaultFace(),
			}
		},
	}
}

func createDefaultFace() text.Face {
	face := &text.GoTextFace{
		Source: _defaultFontResource,
		Size:   FontSizeBody,
	}
	face.SetVariation(_fontTagWeight, float32(FontWeightNormal))
	face.SetVariation(_fontTagItalic, 0)
	return face
}

func (t *textImpl) Layout(bounds image.Rectangle) image.Rectangle {
	w, h := text.Measure(t.content, t.face, t.style.LineHeight)

	t.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X+int(w),
		bounds.Min.Y+int(h),
	)

	return t.frame
}

func (t *textImpl) Draw(screen *ebiten.Image) {
	if t.contentGetter != nil {
		t.content = t.contentGetter()
	}

	if t.content == "" {
		return
	}

	// 計算文字位置
	x := float64(t.frame.Min.X)
	y := float64(t.frame.Min.Y)

	// 處理對齊
	w, _ := text.Measure(t.content, t.face, t.style.LineHeight)
	switch t.style.Alignment {
	case TextAlignCenter:
		x += float64(t.frame.Dx()-int(w)) / 2
	case TextAlignRight:
		x += float64(t.frame.Dx() - int(w))
	}

	// 繪製文字
	for i, gl := range text.AppendGlyphs(nil, t.content, t.face, &text.LayoutOptions{
		LineSpacing: t.style.LineHeight,
	}) {
		if gl.Image == nil {
			continue
		}

		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(t.style.Color)

		// 處理字間距
		opt.GeoM.Translate(float64(i)*t.style.LetterSpacing, 0)
		opt.GeoM.Translate(x+gl.X, y+gl.Y)

		screen.DrawImage(gl.Image, opt)
	}
}

func (t *textImpl) Build() View {
	return t
}

func (v ViewBuilder) WithStyle(style TextStyle) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			t := v.Build().(*textImpl)

			// 創建新的字體
			face := &text.GoTextFace{
				Source: _defaultFontResource,
				Size:   style.Size,
			}
			face.SetVariation(_fontTagWeight, style.Weight)
			face.SetVariation(_fontTagItalic, boolToFloat32(style.Italic))

			return &textImpl{
				content: t.content,
				style:   style,
				cache:   NewViewCache(),
				face:    face,
			}
		},
	}
}

func boolToFloat32(b bool) float32 {
	if b {
		return 1.0
	}
	return 0.0
}

func DynamicText(getter func() string) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &textImpl{
				content: getter(),
				style: TextStyle{
					Size:          FontSizeBody,
					Weight:        FontWeightNormal,
					Color:         color.Black,
					LineHeight:    defaultLineSpacing,
					LetterSpacing: defaultKerning,
					Alignment:     TextAlignLeft,
					Italic:        false,
				},
				cache:         NewViewCache(),
				face:          createDefaultFace(),
				contentGetter: getter,
			}
		},
	}
}
