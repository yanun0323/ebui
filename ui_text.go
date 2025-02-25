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
		return Text(Bind(content))
	case *Binding[string]:
		v := &textImpl{
			content: content,
		}
		v.viewCtx = newViewContext(v)
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
		Source: _defaultFontResource,
		Size:   size.F64(),
	}
	face.SetVariation(_fontTagWeight, weight.F32())
	if italic {
		face.SetVariation(_fontTagItalic, 1)
	}

	faceTable.Store(key, face)
	return face
}

func (t *textImpl) userSetFrameSize() flexibleSize {
	ctxUserSetFrameSize := t.viewCtx.userSetFrameSize()
	w, h := text.Measure(t.content.Get(), t.face(), t.fontLineHeight.Get())

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

func (t *textImpl) preload(parent *viewCtxEnv) (flexibleSize, CGInset, layoutFunc) {
	frameSize, padding, layoutFn := t.viewCtx.preload(parent)
	w, h := text.Measure(t.content.Get(), t.face(), t.fontLineHeight.Get())
	return frameSize, padding, func(start CGPoint, flexFrameSize CGSize) CGRect {
		if isInf(flexFrameSize.Width) {
			flexFrameSize.Width = w
		}

		if isInf(flexFrameSize.Height) {
			flexFrameSize.Height = h
		}

		result := layoutFn(start, flexFrameSize)
		t.viewCtx.debugPrint(result)
		return result
	}
}

func (t *textImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := t.viewCtx.draw(screen, hook...)

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

	return op
}
