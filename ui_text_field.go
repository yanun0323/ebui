package ebui

import (
	"image/color"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	textFieldPlaceholderColor     = NewColor(128, 128, 128, 128)
	textFieldBackgroundColor      = NewColor(128, 128, 128, 128)
	textFieldUnfocusedBorderColor = NewColor(255, 64, 64, 64)
)

type textFieldImpl struct {
	*zstackImpl
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
	zs := ZStack(text).BackgroundColor(Bind(textFieldBackgroundColor)).Border(Const(NewInset(1, 1, 1, 1)), focusedColor).(*zstackImpl)
	tf := &textFieldImpl{
		zstackImpl:   zs,
		text:         text,
		isFocused:    false,
		focusedColor: focusedColor,
		content:      content,
		placeholder:  ph,
	}
	zs.viewCtx._owner = tf

	globalEventManager.RegisterHandler(tf)
	return tf
}

func (t *textFieldImpl) setFocused(focused bool) {
	logf("setFocused: %v", focused)
	t.isFocused = focused
	if focused {
		t.focusedColor.Set(AccentColor)
	} else {
		t.focusedColor.Set(textFieldUnfocusedBorderColor)
	}
}

func (t *textFieldImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	t.zstackImpl.draw(screen, hook...)

	if t.isFocused && time.Now().Unix()%4%2 == 0 {
		_, h := t.text.measure(t.content.Get())
		rect := t.systemSetBounds()
		opt := t.zstackImpl.drawOption(rect, hook...)
		opt.GeoM.Translate(3, 3)

		stroke := ebiten.NewImage(1, int(h)-6)
		stroke.Fill(color.White)
		screen.DrawImage(stroke, opt)
	}
}

func (t *textFieldImpl) HandleTouchEvent(event touchEvent) {
	switch event.Phase {
	case touchPhaseBegan:
		if event.Position.In(t.systemSetBounds()) {
			t.setFocused(true)
		} else {
			t.setFocused(false)
		}
	case touchPhaseMoved:
	case touchPhaseEnded, touchPhaseCancelled:
	}
}

func (t *textFieldImpl) HandleKeyEvent(event keyEvent) {
	if !t.isFocused {
		return
	}

	if event.Phase != keyPhaseJustReleased {
		switch event.Key {
		case ebiten.KeyBackspace:
			t.content.Update(func(val string) string {
				return removeLastChar(val)
			})
			t.cursorPos--
		}
	}
}

func (t *textFieldImpl) HandleInputEvent(event inputEvent) {
	if t.isFocused {
		content := t.content.Get()
		if t.cursorPos == utf8.RuneCountInString(content) {
			t.cursorPos += 1
		}
		t.content.Update(func(val string) string {
			return val + string(event.Char)
		})
	}
}
