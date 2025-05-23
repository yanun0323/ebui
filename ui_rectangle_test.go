package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRectangle(t *testing.T) {
	suite.Run(t, new(RectangleSuite))
}

type RectangleSuite struct {
	suite.Suite
}

func (su *RectangleSuite) Test() {
	rect := Rectangle().(*rectangleImpl)
	rect.Frame(Bind(NewSize(100, 100)))

	s := rect.viewCtx.userSetFrameSize()
	su.Equal(100.0, s.Width)
	su.Equal(100.0, s.Height)

	r := rect.viewCtx.systemSetFrame()
	su.Equal(CGPoint{}, r.Start)
	su.Equal(CGPoint{}, r.End)
}
