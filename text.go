package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
)

type textImpl struct {
	*viewContext

	content *Binding[string]
	frame   image.Rectangle
	cache   *ViewCache
	face    text.Face
}

func Text[T string | *Binding[string]](content T) SomeView {
	switch content := any(content).(type) {
	case string:
		return Text(NewBinding(content))
	case *Binding[string]:
		v := &textImpl{
			content: content,
			cache:   NewViewCache(),
			face:    createDefaultFace(),
		}
		v.viewContext = NewViewContext(v)
		// 添加監聽器，當內容改變時標記需要更新
		content.addListener(func() {
			defaultStateManager.markDirty()
		})

		return v
	}

	return nil
}

func createDefaultFace() text.Face {
	face := &text.GoTextFace{
		Source: _defaultFontResource,
		Size:   font.Body.F64(),
	}
	face.SetVariation(_fontTagWeight, font.Normal.F32())
	face.SetVariation(_fontTagItalic, 0)
	return face
}

func (t *textImpl) layout(bounds image.Rectangle) image.Rectangle {
	w, h := text.Measure(t.content.Get(), t.face, t.viewContext.fontLineHeight.Get())

	t.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X+int(w),
		bounds.Min.Y+int(h),
	)

	return t.frame
}

func (t *textImpl) draw(screen *ebiten.Image) {
	content := t.content.Get()
	if content == "" {
		return
	}

	// 計算文字位置
	x := float64(t.frame.Min.X)
	y := float64(t.frame.Min.Y)

	// 處理對齊
	w, _ := text.Measure(content, t.face, t.viewContext.fontLineHeight.Get())
	switch t.viewContext.fontAlignment.Get() {
	case font.AlignCenter:
		x += float64(t.frame.Dx()-int(w)) / 2
	case font.AlignRight:
		x += float64(t.frame.Dx() - int(w))
	}

	// 繪製文字
	for i, gl := range text.AppendGlyphs(nil, content, t.face, &text.LayoutOptions{
		LineSpacing: t.viewContext.fontLineHeight.Get(),
	}) {
		if gl.Image == nil {
			continue
		}

		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(t.viewContext.foregroundColor.Get())

		opt.GeoM.Translate(float64(i)*t.viewContext.fontLetterSpacing.Get(), 0)
		opt.GeoM.Translate(x+gl.X, y+gl.Y)

		screen.DrawImage(gl.Image, opt)
	}
}
