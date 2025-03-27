package ebui

import (
	"unicode/utf8"

	"github.com/yanun0323/ebui/animation"
)

// TODO: Implement selectable text

type selectableText struct {
	content *Binding[string]
	cursor  int
}

func newSelectableText(content *Binding[string]) *selectableText {
	s := &selectableText{
		content: content,
		cursor:  utf8.RuneCountInString(content.Value()),
	}

	content.AddListener(func(oldVal, newVal string, animStyle ...animation.Style) {
		if s.cursor >= utf8.RuneCountInString(oldVal) {
			s.cursor = utf8.RuneCountInString(newVal)
		}
	})

	return s
}

func (s *selectableText) forwardCursor() {
	s.cursor++
}

func (s *selectableText) backwardCursor() {
	s.cursor--
}

func (s *selectableText) removeCharAtCursor() {
	content := []rune(s.content.Value())
	content = append(content[:s.cursor], content[s.cursor+1:]...)
	s.content.Set(string(content))
}

func (s *selectableText) insertCharAtCursor(char rune) {
	content := []rune(s.content.Value())
	content = append(content[:s.cursor], append([]rune{char}, content[s.cursor:]...)...)
	s.cursor++
	s.content.Set(string(content))
}
