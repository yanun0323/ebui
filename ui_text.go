package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
)

type textImpl struct {
	*ctx

	content *Binding[string]
	face    text.Face
}

func Text[T string | *Binding[string]](content T) SomeView {
	switch content := any(content).(type) {
	case string:
		return Text(Bind(content))
	case *Binding[string]:
		v := &textImpl{
			content: content,
			face:    createDefaultFace(),
		}
		v.ctx = newViewContext(v)
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

func (t *textImpl) userSetFrameSize() flexibleCGSize {
	ctxUserSetFrameSize := t.ctx.userSetFrameSize()
	w, h := text.Measure(t.content.Get(), t.face, t.ctx.fontLineHeight.Get())

	if ctxUserSetFrameSize.IsInfX {
		ctxUserSetFrameSize.Frame.Width = w
		ctxUserSetFrameSize.IsInfX = false
	}

	if ctxUserSetFrameSize.IsInfY {
		ctxUserSetFrameSize.Frame.Height = h
		ctxUserSetFrameSize.IsInfY = false
	}

	return ctxUserSetFrameSize
}

func (t *textImpl) preload() (flexibleCGSize, Inset, layoutFunc) {
	frameSize, padding, layoutFn := t.ctx.preload()
	w, h := text.Measure(t.content.Get(), t.face, t.ctx.fontLineHeight.Get())
	return frameSize, padding, func(start CGPoint, flexFrameSize CGSize) CGRect {
		if isInf(flexFrameSize.Width) {
			flexFrameSize.Width = w
		}

		if isInf(flexFrameSize.Height) {
			flexFrameSize.Height = h
		}

		result := layoutFn(start, flexFrameSize)
		t.ctx.debugPrint(result)
		return result
	}
}

func (t *textImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := t.ctx.draw(screen, hook...)

	content := t.content.Get()
	if content == "" {
		return &ebiten.DrawImageOptions{}
	}

	// 計算文字位置
	// frame := t.systemSetFrame()
	// x := frame.Start.X
	// y := frame.Start.Y

	// 處理對齊
	// w, _ := text.Measure(content, t.face, t.ctx.fontLineHeight.Get())
	// switch t.ctx.fontAlignment.Get() {
	// case font.AlignCenter:
	// 	x += float64(frame.Dx()-w) / 2
	// case font.AlignRight:
	// 	x += float64(frame.Dx() - w)
	// }

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
		opt.GeoM.Translate(gl.X, gl.Y)
		opt.GeoM.Concat(op.GeoM)
		opt.ColorScale.ScaleWithColorScale(op.ColorScale)

		screen.DrawImage(gl.Image, opt)
	}

	return op
}
