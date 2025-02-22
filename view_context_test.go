package ebui

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
)

func newTestViewContextForTest() *ctx {
	return newViewContext(tagNone, newViewContext(tagNone, nil))
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
		r := ctx.userSetFrame()
		su.Inf(r.End.X)
		su.Inf(r.End.Y)

		r = ctx.systemSetFrame()
		su.Equal(ptZero, r.Start)
		su.Equal(ptZero, r.End)
	}
	{
		ctx := newTestViewContextForTest()
		ctx.Frame(nil, Bind(100.0))

		r := ctx.userSetFrame()
		su.Inf(r.End.X)
		su.Equal(100.0, r.End.Y)

		r = ctx.systemSetFrame()
		su.Equal(ptZero, r.Start)
		su.Equal(ptZero, r.End)
	}

	{
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))

		r := ctx.userSetFrame()
		su.Equal(ptZero, r.Start)
		su.Equal(pt(100, 100), r.End)

		r = ctx.systemSetFrame()
		su.Equal(ptZero, r.Start)
		su.Equal(ptZero, r.End)
	}
}

func (su *ViewContextSuite) TestPreload() {
	{ // 沒有設定大小
		ctx := newTestViewContextForTest()

		s, inset, layoutFn := ctx.preload()
		su.Inf(s.Width)
		su.Inf(s.Height)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		res := layoutFn(ptZero, sz(500, 500))
		su.Equal(ptZero, res.Start)
		su.Equal(pt(500, 500), res.End)

		su.Equal(rect(0, 0, 500, 500), ctx.systemSetFrame())
	}

	{ // 沒有設定大小，有設定 padding
		ctx := newTestViewContextForTest()
		ctx.Padding(Bind(10.0))

		s, inset, layoutFn := ctx.preload()
		su.Inf(s.Width)
		su.Inf(s.Height)
		su.Equal(ins(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		res := layoutFn(ptZero, sz(500, 500))
		su.Equal(ptZero, res.Start)
		su.Equal(pt(500, 500), res.End)

		su.Equal(rect(10, 10, 490, 490), ctx.systemSetFrame())
	}

	{ // 設定大小
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))

		s, inset, layoutFn := ctx.preload()
		su.Equal(100.0, s.Width)
		su.Equal(100.0, s.Height)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		res := layoutFn(ptZero, sz(500, 500))
		su.Equal(ptZero, res.Start)
		su.Equal(pt(100, 100), res.End)

		su.Equal(rect(0, 0, 100, 100), ctx.systemSetFrame())
	}

	{ // 設定大小，有設定 padding
		ctx := newTestViewContextForTest()
		ctx.Frame(Bind(100.0), Bind(100.0))
		ctx.Padding(Bind(10.0))

		s, inset, layoutFn := ctx.preload()
		su.Equal(100.0, s.Width)
		su.Equal(100.0, s.Height)
		su.Equal(ins(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		res := layoutFn(ptZero, sz(500, 500))
		su.Equal(ptZero, res.Start)
		su.Equal(pt(100, 100), res.End)

		su.Equal(rect(10, 10, 90, 90), ctx.systemSetFrame())
	}
}
