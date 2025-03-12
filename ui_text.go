package ebui

import (
	"fmt"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
)

var faceTable = sync.Map{}

type textImpl struct {
	*viewCtx

	content *Binding[string]
}

func Text[T string | *Binding[string]](content T) SomeView {
	switch content := any(content).(type) {
	case string:
		return newText(Const(content))
	case *Binding[string]:
		return newText(content)
	}

	return nil
}

func newText(content *Binding[string]) SomeView {
	v := &textImpl{
		content: content,
	}
	v.viewCtx = newViewContext(v)
	return v
}

func (textImpl) faceKey(size font.Size, weight font.Weight, italic bool) string {
	return fmt.Sprintf("%d-%d-%t", size, weight, italic)
}

func (t *textImpl) face() text.Face {
	size := font.Body
	if t.fontSize != nil {
		size = t.fontSize.Get()
	}

	weight := font.Normal
	if t.fontWeight != nil {
		weight = t.fontWeight.Get()
	}

	italic := false
	if t.fontItalic != nil {
		italic = t.fontItalic.Get()
	}

	key := t.faceKey(size, weight, italic)
	{
		face, ok := faceTable.Load(key)
		if ok {
			return face.(*text.GoTextFace)
		}
	}

	face := &text.GoTextFace{
		Source: defaultFontResource,
		Size:   size.F64(),
	}
	face.SetVariation(fontTagWeight, weight.F32())
	if italic {
		face.SetVariation(fontTagItalic, 1)
	}

	faceTable.Store(key, face)
	return face
}

func (t *textImpl) userSetFrameSize() CGSize {
	ctxUserSetFrameSize := t.viewCtx.userSetFrameSize()
	w, h := t.measure(t.content.Get())

	if ctxUserSetFrameSize.IsInfWidth() {
		ctxUserSetFrameSize.Width = w
	}

	if ctxUserSetFrameSize.IsInfHeight() {
		ctxUserSetFrameSize.Height = h
	}

	return ctxUserSetFrameSize
}

func (t *textImpl) measure(content string) (w, h float64) {
	if len(content) == 0 {
		content = " "
	}
	w, h = text.Measure(content, t.face(), t.fontLineHeight.Get())
	return
}

func (t *textImpl) preload(parent *viewCtxEnv) (preloadData, layoutFunc) {
	data, layoutFn := t.viewCtx.preload(parent)
	w, h := text.Measure(t.content.Get(), t.face(), t.fontLineHeight.Get())
	return data, func(start CGPoint, flexBoundsSize CGSize) (CGRect, alignFunc) {
		flexFrameSize := flexBoundsSize.Shrink(data.Padding).Shrink(data.Border)
		if isInf(flexFrameSize.Width) {
			flexFrameSize.Width = w
		}

		if isInf(flexFrameSize.Height) {
			flexFrameSize.Height = h
		}

		result, _ := layoutFn(start, flexFrameSize)
		t.viewCtx.debugPrintPreload(result, flexFrameSize, data)
		return result, t.viewCtx.align
	}
}

func (t *textImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	t.viewCtx.draw(screen, hook...)

	content := t.content.Get()
	if content == "" {
		return
	}

	op := t.viewCtx.drawOption(t.systemSetBounds(), hook...)

	// 繪製文字
	face := t.face()
	for i, gl := range text.AppendGlyphs(nil, content, face, &text.LayoutOptions{
		LineSpacing: t.fontLineHeight.Get(),
	}) {
		if gl.Image == nil {
			continue
		}

		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(t.foregroundColor.Get())

		opt.GeoM.Translate(float64(i)*t.fontLetterSpacing.Get(), 0)
		opt.GeoM.Translate(gl.X, gl.Y)
		opt.GeoM.Concat(op.GeoM)
		opt.ColorScale.ScaleWithColorScale(op.ColorScale)

		screen.DrawImage(gl.Image, opt)
	}
}
