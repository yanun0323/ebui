package ebui

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
)

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
		ctx := Rectangle().(*rectangleImpl)
		s := ctx.userSetFrameSize()
		su.True(s.IsInfWidth(), "%2.f", s.Width)
		su.True(s.IsInfHeight(), "%2.f", s.Height)

		r := ctx.systemSetFrame()
		su.Equal(CGPoint{}, r.Start)
		su.Equal(CGPoint{}, r.End)
	}
	{
		ctx := Rectangle().(*rectangleImpl)
		ctx.Frame(Bind(NewSize(Inf, 100.0)))

		s := ctx.userSetFrameSize()
		su.True(s.IsInfWidth())
		su.Equal(100.0, s.Height)

		r := ctx.systemSetFrame()
		su.Equal(CGPoint{}, r.Start)
		su.Equal(CGPoint{}, r.End)
	}

	{
		ctx := Rectangle().(*rectangleImpl)
		ctx.Frame(Bind(NewSize(100.0, 100.0)))

		s := ctx.userSetFrameSize()
		su.Equal(100.0, s.Width)
		su.Equal(100.0, s.Height)

		r := ctx.systemSetFrame()
		su.Equal(CGPoint{}, r.Start)
		su.Equal(CGPoint{}, r.End)
	}
}

func (su *ViewContextSuite) TestPreload() {
	{ // no size
		ctx := Rectangle().(*rectangleImpl)

		data, layoutFn := ctx.preload(nil)
		su.True(data.FrameSize.IsInfWidth())
		su.True(data.FrameSize.IsInfHeight())
		su.Equal(NewInset(0, 0, 0, 0), data.Padding)
		su.NotNil(layoutFn)

		res, _ := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(CGPoint{}, res.Start)
		su.Equal(NewPoint(500, 500), res.End)

		su.Equal(NewRect(0, 0, 500, 500), ctx.systemSetFrame())
	}

	{ // no size, with padding
		ctx := Rectangle().(*rectangleImpl)
		zs := ctx.Padding(Bind(CGInset{10, 10, 10, 10})).(*zstackImpl)

		data, layoutFn := zs.preload(nil)
		su.Equal(NewSize(0, 0), data.FrameSize)
		su.True(data.IsInfWidth)
		su.True(data.IsInfHeight)
		su.Equal(CGInset{10, 10, 10, 10}, data.Padding)
		su.NotNil(layoutFn)

		res, _ := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(NewPoint(0, 0), res.Start)
		su.Equal(NewPoint(500, 500), res.End)

		su.Equal(NewRect(10, 10, 490, 490), ctx.systemSetFrame())
		su.Equal(NewRect(10, 10, 490, 490), zs.systemSetFrame())
	}

	{ // set size
		ctx := Rectangle().(*rectangleImpl)
		ctx.Frame(Bind(NewSize(100.0, 100.0)))

		data, layoutFn := ctx.preload(nil)
		su.Equal(100.0, data.FrameSize.Width)
		su.Equal(100.0, data.FrameSize.Height)
		su.Equal(NewInset(0, 0, 0, 0), data.Padding)
		su.NotNil(layoutFn)

		res, _ := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(CGPoint{}, res.Start)
		su.Equal(NewPoint(100, 100), res.End)

		su.Equal(NewRect(0, 0, 100, 100), ctx.systemSetFrame())
	}

	{ // set size, with padding
		ctx := Rectangle().(*rectangleImpl)
		ctx.Frame(Bind(NewSize(100.0, 100.0)))
		ctx.Padding(Bind(CGInset{10, 10, 10, 10}))

		data, layoutFn := ctx.preload(nil)
		su.Equal(100.0, data.FrameSize.Width)
		su.Equal(100.0, data.FrameSize.Height)
		su.Equal(CGInset{}, data.Padding)
		su.NotNil(layoutFn)

		res, _ := layoutFn(CGPoint{}, NewSize(500, 500))
		su.Equal(CGPoint{}, res.Start)
		su.Equal(NewPoint(100, 100), res.End)

		su.Equal(NewRect(0, 0, 100, 100), ctx.systemSetFrame())
	}
}
