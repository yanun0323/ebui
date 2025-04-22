package ebui

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
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

func newText(content *Binding[string]) *textImpl {
	v := &textImpl{
		content: content,
		cache:   helper.NewHashCache[*ebiten.Image](),
	}
	v.viewCtx = newViewContext(v)
	return v
}

func (t *textImpl) getContent() []string {
	lines := strings.Split(t.content.Value(), "\n")
	if lineLimit := t.lineLimit.Value(); lineLimit >= 1 {
		result := make([]string, 0, lineLimit)
		buffer := strings.Builder{}
		for i, line := range lines {
			if i+1 <= lineLimit {
				result = append(result, line)
			} else {
				buffer.WriteString(line)
			}
		}
		lines = append(result, buffer.String())
	}

	return lines
}

func (textImpl) faceKey(size font.Size, weight font.Weight, italic bool) string {
	return fmt.Sprintf("%d-%d-%t", size, weight, italic)
}

func (t *textImpl) face(s ...font.Size) text.Face {
	size := t.fontSize.Value()
	if len(s) != 0 {
		size = s[0]
	}

	if size == 0 {
		size = font.Body
	}

	weight := font.Normal
	if t.fontWeight != nil {
		weight = t.fontWeight.Value()
	}

	italic := false
	if t.fontItalic != nil {
		italic = t.fontItalic.Value()
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
	lines := t.getContent()
	w, h, _ := t.measure(lines)

	if ctxUserSetFrameSize.IsInfWidth() {
		ctxUserSetFrameSize.Width = w
	}

	if ctxUserSetFrameSize.IsInfHeight() {
		ctxUserSetFrameSize.Height = h
	}

	return ctxUserSetFrameSize
}

func (t *textImpl) measure(lines []string) (w, h, lineHeight float64) {
	if len(lines) == 0 {
		lines = []string{" "}
	}

	var (
		kerning          = t.fontKerning.Value()
		lineSpacing      = t.fontLineHeight.Value()
		face             = t.face()
		maxLineRuneCount = 0
		maxW             = 0.0
		totalH           = 0.0
		lineH            = 0.0
	)

	for _, line := range lines {
		w, h := text.Measure(line, face, kerning)
		maxW = max(w, maxW)
		maxLineRuneCount = max(utf8.RuneCountInString(line), maxLineRuneCount)
		if h == 0 {
			lineH = h
		} else {
			lineH = h + lineSpacing
		}
		totalH += lineH
	}

	return maxW + float64(maxLineRuneCount-1)*kerning, totalH, lineH
}

func (t *textImpl) preload(parent *viewCtx, _ ...stackType) (preloadData, layoutFunc) {
	data, layoutFn := t.viewCtx.preload(parent)
	return data, func(start CGPoint, flexBoundsSize CGSize) (CGRect, alignFunc, bool) {
		flexFrameSize := flexBoundsSize.Shrink(data.Padding).Shrink(data.Border)
		if isInf(flexFrameSize.Width) {
			flexFrameSize.Width = data.FrameSize.Width
		}

		if isInf(flexFrameSize.Height) {
			flexFrameSize.Height = data.FrameSize.Height
		}

		result, _, cached := layoutFn(start, flexFrameSize)
		t.cache.SetNextHash(t.Hash())
		t.viewCtx.debugPrintPreload(result, flexFrameSize, data)
		return result, t.viewCtx.align, cached && t.cache.IsNextHashCached()
	}
}

func (t *textImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	t.viewCtx.draw(screen, hook...)

	content := t.content.Value()
	if content == "" {
		return
	}

	op := t.viewCtx.drawOption(t.systemSetBounds(), hook...)
	op.ColorScale.ScaleWithColor(t.foregroundColor.Value())
	if t.cache.IsNextHashCached() {
		screen.DrawImage(t.cache.Load(), op)
		return
	}

	bounds := t.systemSetBounds()
	if !bounds.drawable() {
		return
	}

	var (
		layoutOpt = &text.LayoutOptions{
			LineSpacing: t.fontLineHeight.Value(),
		}
		textBase    = ebiten.NewImage(int(bounds.Dx()), int(bounds.Dy()))
		kerning     = t.fontKerning.Value()
		face        = t.face()
		lines       = t.getContent()
		_, _, lineH = t.measure(lines)
	)

	for i, line := range lines {
		for j, gl := range text.AppendGlyphs(nil, line, face, layoutOpt) {
			if gl.Image == nil {
				continue
			}

			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(j)*kerning, 0)
			opt.GeoM.Translate(gl.X, gl.Y)
			opt.GeoM.Translate(0, float64(i)*lineH)

			textBase.DrawImage(gl.Image, opt)
		}
	}

	t.cache.Update(textBase)
	screen.DrawImage(textBase, op)
}

func (t *textImpl) Hash() string {
	h := xxhash.New()
	h.Write(t.viewCtx.Bytes(true))
	h.Write(helper.BytesString(t.content.Value()))
	return strconv.FormatUint(h.Sum64(), 16)
}

// Preview_Text
func Preview_Text() View {
	return Text(Const("Hello, World2!"))
}
