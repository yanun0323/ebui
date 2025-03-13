package ebui

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/cespare/xxhash/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/internal/helper"
)

var faceTable = sync.Map{}

type textImpl struct {
	*viewCtx

	content *Binding[string]
	cache   *helper.HashCache[*ebiten.Image]
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
		cache:   helper.NewHashCache[*ebiten.Image](),
	}
	v.viewCtx = newViewContext(v)
	return v
}

func (textImpl) faceKey(size font.Size, weight font.Weight, italic bool) string {
	return fmt.Sprintf("%d-%d-%t", size, weight, italic)
}

func (t *textImpl) face(s ...font.Size) text.Face {
	size := t.fontSize.Get()
	if len(s) != 0 {
		size = s[0]
	}

	if size == 0 {
		size = font.Body
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
		Source:    defaultFontResource,
		Direction: text.DirectionLeftToRight,
		Size:      size.F64(),
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
	kerning := t.fontKerning.Get()
	w, h = text.Measure(content, t.face(), kerning)
	w += float64(utf8.RuneCountInString(content)-1) * kerning
	return w, h
}

func (t *textImpl) preload(parent *viewCtxEnv, _ ...formulaType) (preloadData, layoutFunc) {
	data, layoutFn := t.viewCtx.preload(parent)
	return data, func(start CGPoint, flexBoundsSize CGSize) (CGRect, alignFunc) {
		flexFrameSize := flexBoundsSize.Shrink(data.Padding).Shrink(data.Border)
		if isInf(flexFrameSize.Width) {
			flexFrameSize.Width = data.FrameSize.Width
		}

		if isInf(flexFrameSize.Height) {
			flexFrameSize.Height = data.FrameSize.Height
		}

		result, _ := layoutFn(start, flexFrameSize)
		t.cache.SetNextHash(t.Hash())
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
	if t.cache.IsNextHashCached() {
		screen.DrawImage(t.cache.Get(), op)
		return
	}

	println(time.Now().UnixMilli(), "text redraw", content)

	bounds := t.systemSetBounds()
	if !bounds.drawable() {
		return
	}

	textBase := ebiten.NewImage(int(bounds.Dx()), int(bounds.Dy()))

	face := t.face()
	foregroundColor := t.foregroundColor.Get()
	kerning := t.fontKerning.Get()
	for i, gl := range text.AppendGlyphs(nil, content, face, &text.LayoutOptions{
		LineSpacing: t.fontLineHeight.Get(),
	}) {
		if gl.Image == nil {
			continue
		}

		opt := &ebiten.DrawImageOptions{}
		opt.ColorScale.ScaleWithColor(foregroundColor)
		opt.GeoM.Translate(float64(i)*kerning, 0)
		opt.GeoM.Translate(gl.X, gl.Y)
		opt.ColorScale.ScaleWithColorScale(op.ColorScale)

		textBase.DrawImage(gl.Image, opt)
	}

	t.cache.Update(textBase)
	screen.DrawImage(textBase, op)
}

func (t *textImpl) Hash() string {
	h := xxhash.New()
	h.Write(t.viewCtx.Bytes(true))
	h.Write(helper.BytesString(t.content.Get()))
	return strconv.FormatUint(h.Sum64(), 16)
}
