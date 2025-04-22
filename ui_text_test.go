package ebui

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/stretchr/testify/suite"
)

func TestText(t *testing.T) {
	suite.Run(t, new(TextSuite))
}

type TextSuite struct {
	suite.Suite
}

func (su *TextSuite) Test() {
	content := "Hello, World!"
	t := Text(content).(*textImpl)

	height := t.fontLineHeight.Value()
	w, h := text.Measure(content, t.face(), height)

	data, layoutFn := t.preload(nil)
	su.Equal(NewSize(w, h), data.FrameSize)
	su.Equal(NewInset(0, 0, 0, 0), data.Padding)

	bound, _, _ := layoutFn(NewPoint(0, 0), NewSize(500.0, 500.0))
	su.Equal(NewPoint(0, 0), bound.Start)
	su.Equal(NewPoint(w, h), bound.End)

	su.Equal(NewSize(w, h), t.systemSetFrame().Size())
}
