package ebui

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/stretchr/testify/suite"
)

func newTextForTest(content string) *textImpl {
	v := &textImpl{
		content: Bind(content),
		face:    createDefaultFace(),
	}
	v.ctx = newViewContext(v)
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
	w, h := text.Measure(content, t.face, height)

	frameSize, padding, layoutFn := t.preload()
	su.Equal(Size(w, h), frameSize.Frame)
	su.Equal(Inset(0, 0, 0, 0), padding)

	bound := layoutFn(Point(0, 0), Size(500.0, 500.0))
	su.Equal(Point(0, 0), bound.Start)
	su.Equal(Point(w, h), bound.End)

	su.Equal(Size(w, h), t.systemSetFrame().Size())
}
