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
	su.Equal(sz(w, h), frameSize.Frame)
	su.Equal(ins(0, 0, 0, 0), padding)

	bound := layoutFn(pt(0, 0), sz(500.0, 500.0))
	su.Equal(pt(0, 0), bound.Start)
	su.Equal(pt(w, h), bound.End)

	su.Equal(sz(w, h), t.systemSetFrame().Size())
}
