package ebui

import (
	"image/color"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/input"
)

// TODO: Implement the TextField

var (
	textFieldPlaceholderColor     = NewColor(128, 128, 128, 128)
	textFieldBackgroundColor      = NewColor(128, 128, 128, 128)
	textFieldUnfocusedBorderColor = NewColor(255, 64, 64, 64)
)

type textFieldImpl struct {
	*stackImpl
	text *textImpl

	isFocused    bool
	focusedColor *Binding[CGColor]
	content      *Binding[string]
	placeholder  *Binding[string]
	cursorPos    int
}

func TextField[T string | *Binding[string]](content *Binding[string], placeholder ...T) SomeView {
	ph := Const("")
	if len(placeholder) != 0 {
		switch phT := any(placeholder[0]).(type) {
		case string:
			return TextField(content, Const(phT))
		case *Binding[string]:
			ph = phT
		}
	}
	text := Text(content).(*textImpl)
	focusedColor := Bind(textFieldUnfocusedBorderColor)
	zs := ZStack(text).BackgroundColor(Bind(textFieldBackgroundColor)).Border(Const(NewInset(1, 1, 1, 1)), focusedColor).(*stackImpl)
	tf := &textFieldImpl{
		stackImpl:    zs,
		text:         text,
		isFocused:    false,
		focusedColor: focusedColor,
		content:      content,
		placeholder:  ph,
	}
	zs.viewCtx._owner = tf

	return tf
}

func (t *textFieldImpl) setFocused(focused bool) {
	t.isFocused = focused
	if focused {
		t.focusedColor.Set(AccentColor.Value())
	} else {
		t.focusedColor.Set(textFieldUnfocusedBorderColor)
	}
}

func (t *textFieldImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	t.stackImpl.draw(screen, hook...)

	if t.isFocused && time.Now().Unix()%4%2 == 0 {
		lines := t.text.getContent()
		_, h, _ := t.text.measure(lines)
		rect := t.systemSetBounds()
		opt := t.stackImpl.drawOption(rect, hook...)
		opt.GeoM.Translate(3, 3)

		stroke := ebiten.NewImage(1, int(h)-6)
		stroke.Fill(color.White)
		screen.DrawImage(stroke, opt)
	}
}

func (t *textFieldImpl) onMouseEvent(event input.MouseEvent) {
	defer t.viewCtx.onMouseEvent(event)

	switch event.Phase {
	case input.MousePhaseBegan:
		if t.isHover(event.Position) {
			t.setFocused(true)
		} else {
			t.setFocused(false)
		}
	case input.MousePhaseMoved:
	case input.MousePhaseEnded, input.MousePhaseCancelled:
	}
}

func (t *textFieldImpl) onKeyEvent(event input.KeyEvent) {
	if !t.isFocused {
		return
	}
	defer t.viewCtx.onKeyEvent(event)

	if event.Phase != input.KeyPhaseJustReleased {
		switch event.Key {
		case input.KeyBackspace:
			t.content.Set(removeLastChar(t.content.Value()))
			t.cursorPos--
		}
	}
}

func (t *textFieldImpl) onTypeEvent(event input.TypeEvent) {
	if !t.isFocused {
		return
	}
	defer t.viewCtx.onTypeEvent(event)

	if t.isFocused {
		content := t.content.Value()
		if t.cursorPos == utf8.RuneCountInString(content) {
			t.cursorPos += 1
		}
		t.content.Set(content + string(event.Char))
	}
}
