package ebui

import "unicode/utf8"

type selectableText struct {
	content *Binding[string]
	cursor  int
}

func newSelectableText(content *Binding[string]) *selectableText {
	s := &selectableText{
		content: content,
		cursor:  utf8.RuneCountInString(content.Get()),
	}

	content.AddListener(func(oldVal, newVal string) {
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
	content := []rune(s.content.Get())
	content = append(content[:s.cursor], content[s.cursor+1:]...)
	s.content.Set(string(content))
}

func (s *selectableText) insertCharAtCursor(char rune) {
	content := []rune(s.content.Get())
	content = append(content[:s.cursor], append([]rune{char}, content[s.cursor:]...)...)
	s.cursor++
	s.content.Set(string(content))
}
