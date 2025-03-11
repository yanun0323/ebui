package ebui

import "github.com/hajimehoshi/ebiten/v2"

type textFieldImpl struct {
	SomeView

	isFocused bool
	content   *Binding[string]
}

func TextField(content *Binding[string]) SomeView {
	return &textFieldImpl{
		SomeView: ZStack(Text(content)),
		content:  content,
	}
}

// Button 的事件處理
func (t *textFieldImpl) HandleTouchEvent(event touchEvent) bool {
	switch event.Phase {
	case touchPhaseBegan:
		if event.Position.In(t.systemSetBounds()) {
			t.isFocused = true
			return true
		}
	case touchPhaseMoved:
		if t.isFocused {
			return true
		}
	case touchPhaseEnded, touchPhaseCancelled:
	}
	return false
}

func (t *textFieldImpl) HandleKeyEvent(event keyEvent) bool {
	if !t.isFocused || !event.Pressed {
		return false
	}

	ebiten.AppendInputChars(nil)

	return false
}
