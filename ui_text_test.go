package ebui

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/stretchr/testify/suite"
)

func newTextForTest(content string) *textImpl {
	v := &textImpl{
		content: Const(content),
	}
	v.viewCtx = newViewContext(v)
	return v
}

func TestText(t *testing.T) {
	suite.Run(t, new(TextSuite))
}

type TextSuite struct {
	suite.Suite
}

func (su *TextSuite) Test() {
	content := "Hello, World!"
	t := newTextForTest(content)

	height := t.fontLineHeight.Get()
	w, h := text.Measure(content, t.face(), height)

	data, layoutFn := t.preload(nil)
	su.Equal(NewSize(w, h), data.FrameSize)
	su.Equal(NewInset(0, 0, 0, 0), data.Padding)

	bound := layoutFn(NewPoint(0, 0), NewSize(500.0, 500.0))
	su.Equal(NewPoint(0, 0), bound.Start)
	su.Equal(NewPoint(w, h), bound.End)

	su.Equal(NewSize(w, h), t.systemSetFrame().Size())
}
