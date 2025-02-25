package ebui

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
)

func newTestViewContextForTest() *ctx {
	return newRectangleForTest().ctx
}

func TestViewContext(t *testing.T) {
	suite.Run(t, new(ViewContextSuite))
}

type ViewContextSuite struct {
	suite.Suite
}

func (su *ViewContextSuite) Inf(f float64) {
	su.T().Helper()
	su.True(math.IsInf(f, 1))
}

func (su *ViewContextSuite) TestSetFrame() {
	{
		ctx := newTestViewContextForTest()
		s := ctx.userSetFrameSize()
		su.True(s.IsInfX)
		su.True(s.IsInfY)

		r := ctx.systemSetFrame()
		su.Equal(Point{}, r.Start)
		su.Equal(Point{}, r.End)
	}
	{
		ctx := newTestViewContextForTest()
		ctx.Frame(nil, Bind(100.0))

		s := ctx.userSetFrameSize()
		su.True(s.IsInfX)
		su.Equal(100.0, s.Frame.Height)

		r := ctx.systemSetFrame()
		su.Equal(Point{}, r.Start)
		su.Equal(Point{}, r.End)
	}

	{
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))

		s := ctx.userSetFrameSize()
		su.Equal(100.0, s.Frame.Width)
		su.Equal(100.0, s.Frame.Height)

		r := ctx.systemSetFrame()
		su.Equal(Point{}, r.Start)
		su.Equal(Point{}, r.End)
	}
}

func (su *ViewContextSuite) TestPreload() {
	{ // 沒有設定大小
		ctx := newTestViewContextForTest()

		s, inset, layoutFn := ctx.preload()
		su.True(s.IsInfX)
		su.True(s.IsInfY)
		su.Equal(CGInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		res := layoutFn(Point{}, CGSize(500, 500))
		su.Equal(Point{}, res.Start)
		su.Equal(CGPoint(500, 500), res.End)

		su.Equal(CGRect(0, 0, 500, 500), ctx.systemSetFrame())
	}

	{ // 沒有設定大小，有設定 padding
		ctx := newTestViewContextForTest()
		ctx.Padding(Bind(Inset{10, 10, 10, 10}))

		s, inset, layoutFn := ctx.preload()
		su.Equal(CGSize(0, 0), s.Frame)
		su.True(s.IsInfX)
		su.True(s.IsInfY)
		su.Equal(CGInset(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		res := layoutFn(Point{}, CGSize(500, 500))
		su.Equal(CGPoint(0, 0), res.Start)
		su.Equal(CGPoint(520, 520), res.End)

		su.Equal(CGRect(10, 10, 510, 510), ctx.systemSetFrame())
	}

	{ // 設定大小
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))

		s, inset, layoutFn := ctx.preload()
		su.Equal(100.0, s.Frame.Width)
		su.Equal(100.0, s.Frame.Height)
		su.Equal(CGInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		res := layoutFn(Point{}, CGSize(500, 500))
		su.Equal(Point{}, res.Start)
		su.Equal(CGPoint(100, 100), res.End)

		su.Equal(CGRect(0, 0, 100, 100), ctx.systemSetFrame())
	}

	{ // 設定大小，有設定 padding
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))
		ctx.Padding(Bind(Inset{10, 10, 10, 10}))

		s, inset, layoutFn := ctx.preload()
		su.Equal(100.0, s.Frame.Width)
		su.Equal(100.0, s.Frame.Height)
		su.Equal(CGInset(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		res := layoutFn(Point{}, CGSize(500, 500))
		su.Equal(Point{}, res.Start)
		su.Equal(CGPoint(120, 120), res.End)

		su.Equal(CGRect(10, 10, 110, 110), ctx.systemSetFrame())
	}
}
