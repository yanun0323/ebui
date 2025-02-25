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
		su.Equal(CGPoint{}, r.Start)
		su.Equal(CGPoint{}, r.End)
	}
	{
		ctx := newTestViewContextForTest()
		ctx.Frame(nil, Bind(100.0))

		s := ctx.userSetFrameSize()
		su.True(s.IsInfX)
		su.Equal(100.0, s.Frame.Height)

		r := ctx.systemSetFrame()
		su.Equal(CGPoint{}, r.Start)
		su.Equal(CGPoint{}, r.End)
	}

	{
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))

		s := ctx.userSetFrameSize()
		su.Equal(100.0, s.Frame.Width)
		su.Equal(100.0, s.Frame.Height)

		r := ctx.systemSetFrame()
		su.Equal(CGPoint{}, r.Start)
		su.Equal(CGPoint{}, r.End)
	}
}

func (su *ViewContextSuite) TestPreload() {
	{ // 沒有設定大小
		ctx := newTestViewContextForTest()

		s, inset, layoutFn := ctx.preload()
		su.True(s.IsInfX)
		su.True(s.IsInfY)
		su.Equal(NewInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		res := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(CGPoint{}, res.Start)
		su.Equal(NewPoint(500, 500), res.End)

		su.Equal(NewRect(0, 0, 500, 500), ctx.systemSetFrame())
	}

	{ // 沒有設定大小，有設定 padding
		ctx := newTestViewContextForTest()
		ctx.Padding(Bind(CGInset{10, 10, 10, 10}))

		s, inset, layoutFn := ctx.preload()
		su.Equal(NewSize(0, 0), s.Frame)
		su.True(s.IsInfX)
		su.True(s.IsInfY)
		su.Equal(NewInset(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		res := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(NewPoint(0, 0), res.Start)
		su.Equal(NewPoint(520, 520), res.End)

		su.Equal(NewRect(10, 10, 510, 510), ctx.systemSetFrame())
	}

	{ // 設定大小
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))

		s, inset, layoutFn := ctx.preload()
		su.Equal(100.0, s.Frame.Width)
		su.Equal(100.0, s.Frame.Height)
		su.Equal(NewInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		res := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(CGPoint{}, res.Start)
		su.Equal(NewPoint(100, 100), res.End)

		su.Equal(NewRect(0, 0, 100, 100), ctx.systemSetFrame())
	}

	{ // 設定大小，有設定 padding
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))
		ctx.Padding(Bind(CGInset{10, 10, 10, 10}))

		s, inset, layoutFn := ctx.preload()
		su.Equal(100.0, s.Frame.Width)
		su.Equal(100.0, s.Frame.Height)
		su.Equal(NewInset(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		res := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(CGPoint{}, res.Start)
		su.Equal(NewPoint(120, 120), res.End)

		su.Equal(NewRect(10, 10, 110, 110), ctx.systemSetFrame())
	}
}
