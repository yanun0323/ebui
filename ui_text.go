package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
)

type textImpl struct {
	*ctx

	content *Binding[string]
	cache   *ViewCache
	face    text.Face
}

func Text[T string | *Binding[string]](content T) SomeView {
	switch content := any(content).(type) {
	case string:
		return Text(Bind(content))
	case *Binding[string]:
		v := &textImpl{
			content: content,
			cache:   NewViewCache(),
			face:    createDefaultFace(),
		}
		v.ctx = newViewContext(tagText, v)
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

func (t *textImpl) getCtxFrame() CGRect {
	frame := t.ctx.systemSetFrame()
	ist := t.ctx.padding()
	w, h := text.Measure(t.content.Get(), t.face, t.ctx.fontLineHeight.Get())

	start := frame.Start.Add(pt(ist.Left, ist.Top))
	end := start.Add(frame.Size().ToCGPoint())
	if isInf(frame.End.X) {
		end.X = start.X + w
	}

	if isInf(frame.End.Y) {
		end.Y = start.Y + h
	}

	frame.Start = start
	frame.End = end

	return frame
}

func (t *textImpl) preload() (CGSize, Inset, func(CGPoint, CGSize) CGRect) {
	_, inset, layoutFn := t.ctx.preload()
	frameSize := t.getCtxFrame().Size()

	return frameSize, inset, func(start CGPoint, flexSize CGSize) CGRect {
		return layoutFn(start, frameSize)
	}
}

func (t *textImpl) draw(screen *ebiten.Image, bounds ...CGRect) {
	t.ctx.draw(screen, bounds...)

	content := t.content.Get()
	if content == "" {
		return
	}

	// 計算文字位置
	frame := t.getCtxFrame()
	x := frame.Start.X
	y := frame.Start.Y

	// 處理對齊
	w, _ := text.Measure(content, t.face, t.ctx.fontLineHeight.Get())
	switch t.ctx.fontAlignment.Get() {
	case font.AlignCenter:
		x += float64(frame.Dx()-w) / 2
	case font.AlignRight:
		x += float64(frame.Dx() - w)
	}

	// 繪製文字
	for i, gl := range text.AppendGlyphs(nil, content, t.face, &text.LayoutOptions{
		LineSpacing: t.ctx.fontLineHeight.Get(),
	}) {
		if gl.Image == nil {
			continue
		}

		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(t.ctx.foregroundColor.Get())

		opt.GeoM.Translate(float64(i)*t.ctx.fontLetterSpacing.Get(), 0)
		opt.GeoM.Translate(x+gl.X, y+gl.Y)

		screen.DrawImage(gl.Image, opt)
	}
}
