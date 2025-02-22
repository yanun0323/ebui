package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func RectangleForTest() (*rectangleImpl, *ctx) {
	rect := &rectangleImpl{}
	ctx := newViewContext(tagRectangle, rect)
	rect.ctx = ctx
	return rect, ctx
}

func TestRectangle(t *testing.T) {
	suite.Run(t, new(RectangleSuite))
}

type RectangleSuite struct {
	suite.Suite
}

func newRectangleForTest() *rectangleImpl {
	rect := &rectangleImpl{}
	rect.ctx = newViewContext(tagRectangle, rect)
	return rect
}

func (su *RectangleSuite) Test() {
	rect := newRectangleForTest()
	rect.Frame(Bind(100.0), Bind(100.0))

	r := rect.ctx.userSetFrame()
	su.Equal(ptZero, r.Start)
	su.Equal(pt(100, 100), r.End)

	r = rect.ctx.systemSetFrame()
	su.Equal(ptZero, r.Start)
	su.Equal(ptZero, r.End)
}
