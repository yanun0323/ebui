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

func newRectangleForTest() *rectangleImpl {
	rect := &rectangleImpl{}
	rect.ctx = newViewContext(rect)
	return rect
}

func (su *RectangleSuite) Test() {
	rect := newRectangleForTest()
	rect.Frame(Bind(100.0), Bind(100.0))

	s := rect.ctx.userSetFrameSize()
	su.Equal(100.0, s.Frame.Width)
	su.Equal(100.0, s.Frame.Height)

	r := rect.ctx.systemSetFrame()
	su.Equal(CGPoint{}, r.Start)
	su.Equal(CGPoint{}, r.End)
}
