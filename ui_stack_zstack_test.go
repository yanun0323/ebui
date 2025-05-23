package ebui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestZStack(t *testing.T) {
	suite.Run(t, new(ZSTestZStackSuite))
}

type ZSTestZStackSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *ZSTestZStackSuite) SetupSuite() {
	su.ctx = context.Background()
}

func (su *ZSTestZStackSuite) TestZStack() {
	{ // child view with no size
		rect1 := Rectangle().(*rectangleImpl)
		z := ZStack(rect1).(*stackImpl)

		data, layoutFn := z.preload(nil)
		su.Equal(NewSize(0, 0), data.FrameSize)
		su.Equal(true, data.IsInfWidth)
		su.Equal(true, data.IsInfHeight)
		su.Equal(NewInset(0, 0, 0, 0), data.Padding)
		su.Equal(NewInset(0, 0, 0, 0), data.Border)

		result, _, _ := layoutFn(NewPoint(0, 0), NewSize(100, 100))
		su.Equal(NewPoint(0, 0), result.Start)
		su.Equal(NewPoint(100, 100), result.End)
	}

	{ // child view with fixed size
		rect1 := Rectangle().Frame(Bind(NewSize(100, 100)))
		z := ZStack(rect1).(*stackImpl)

		data, layoutFn := z.preload(nil)
		su.Equal(NewSize(100, 100), data.FrameSize)
		su.Equal(false, data.IsInfWidth)
		su.Equal(false, data.IsInfHeight)
		su.Equal(NewInset(0, 0, 0, 0), data.Padding)
		su.Equal(NewInset(0, 0, 0, 0), data.Border)

		result, _, _ := layoutFn(NewPoint(0, 0), NewSize(200, 200))
		su.Equal(NewPoint(0, 0), result.Start)
		su.Equal(NewPoint(100, 100), result.End)
	}

	{ // child view with no size + Padding
		rect1 := Rectangle().(*rectangleImpl)
		z := ZStack(rect1).(*stackImpl)
		zz := (z.Padding(Bind(CGInset{10, 10, 10, 10}))).(*stackImpl)

		data, layoutFn := zz.preload(nil)
		su.Equal(NewSize(0, 0), data.FrameSize)
		su.Equal(true, data.IsInfWidth)
		su.Equal(true, data.IsInfHeight)
		su.Equal(NewInset(10, 10, 10, 10), data.Padding)
		su.Equal(NewInset(0, 0, 0, 0), data.Border)

		final, _, _ := layoutFn(NewPoint(0, 0), NewSize(200, 200))
		su.Equal(NewPoint(0, 0), final.Start)
		su.Equal(NewPoint(200, 200), final.End)

		rect1Frame := rect1.viewCtx.systemSetFrame()
		su.Equal(NewPoint(10, 10), rect1Frame.Start)
		su.Equal(NewPoint(190, 190), rect1Frame.End)

		zFrame := z.viewCtx.systemSetFrame()
		su.Equal(NewPoint(10, 10), zFrame.Start)
		su.Equal(NewPoint(190, 190), zFrame.End)

		zzFrame := zz.viewCtx.systemSetFrame()
		su.Equal(NewPoint(10, 10), zzFrame.Start)
		su.Equal(NewPoint(190, 190), zzFrame.End)
	}

	{ // child view with fixed size + Padding
		rect1 := Rectangle().Frame(Bind(NewSize(100, 100))).(*rectangleImpl)
		z := ZStack(rect1).(*stackImpl)
		zz := (z.Padding(Bind(CGInset{10, 10, 10, 10}))).(*stackImpl)

		data, layoutFn := zz.preload(nil)
		su.Equal(NewSize(100, 100), data.FrameSize)
		su.Equal(false, data.IsInfWidth)
		su.Equal(false, data.IsInfHeight)
		su.Equal(NewInset(10, 10, 10, 10), data.Padding)
		su.Equal(NewInset(0, 0, 0, 0), data.Border)

		final, _, _ := layoutFn(NewPoint(0, 0), NewSize(200, 200))
		su.Equal(NewPoint(0, 0), final.Start)
		su.Equal(NewPoint(120, 120), final.End)

		rect1Frame := rect1.viewCtx.systemSetFrame()
		su.Equal(NewPoint(10, 10), rect1Frame.Start)
		su.Equal(NewPoint(110, 110), rect1Frame.End)

		zFrame := z.viewCtx.systemSetFrame()
		su.Equal(NewPoint(10, 10), zFrame.Start)
		su.Equal(NewPoint(110, 110), zFrame.End)

		zzFrame := zz.viewCtx.systemSetFrame()
		su.Equal(NewPoint(10, 10), zzFrame.Start)
		su.Equal(NewPoint(110, 110), zzFrame.End)
	}

	{ // child view with no size + Padding + Border
		rect1 := Rectangle().(*rectangleImpl)
		z := ZStack(rect1).(*stackImpl)
		zz := (z.Padding(Bind(CGInset{10, 10, 10, 10})).
			Border(Bind(CGInset{10, 10, 10, 10}), Bind(white))).(*stackImpl)

		data, layoutFn := zz.preload(nil)
		su.Equal(NewSize(0, 0), data.FrameSize)
		su.Equal(true, data.IsInfWidth)
		su.Equal(true, data.IsInfHeight)
		su.Equal(NewInset(10, 10, 10, 10), data.Padding)
		su.Equal(NewInset(10, 10, 10, 10), data.Border)

		final, _, _ := layoutFn(NewPoint(0, 0), NewSize(200, 200))
		su.Equal(NewPoint(0, 0), final.Start)
		su.Equal(NewPoint(200, 200), final.End)

		rect1Frame := rect1.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), rect1Frame.Start)
		su.Equal(NewPoint(180, 180), rect1Frame.End)

		zFrame := z.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), zFrame.Start)
		su.Equal(NewPoint(180, 180), zFrame.End)

		zzFrame := zz.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), zzFrame.Start)
		su.Equal(NewPoint(180, 180), zzFrame.End)
	}

	{ // child view with fixed size + Padding + Border
		rect1 := Rectangle().Frame(Bind(NewSize(100, 100))).(*rectangleImpl)
		z := ZStack(rect1).(*stackImpl)
		zz := (z.Padding(Bind(CGInset{10, 10, 10, 10})).
			Border(Bind(CGInset{10, 10, 10, 10}), Bind(white))).(*stackImpl)

		data, layoutFn := zz.preload(nil)
		su.Equal(NewSize(100, 100), data.FrameSize)
		su.Equal(false, data.IsInfWidth)
		su.Equal(false, data.IsInfHeight)
		su.Equal(NewInset(10, 10, 10, 10), data.Padding)
		su.Equal(NewInset(10, 10, 10, 10), data.Border)

		final, _, _ := layoutFn(NewPoint(0, 0), NewSize(200, 200))
		su.Equal(NewPoint(0, 0), final.Start)
		su.Equal(NewPoint(140, 140), final.End)

		rect1Frame := rect1.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), rect1Frame.Start)
		su.Equal(NewPoint(120, 120), rect1Frame.End)

		zFrame := z.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), zFrame.Start)
		su.Equal(NewPoint(120, 120), zFrame.End)

		zzFrame := zz.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), zzFrame.Start)
		su.Equal(NewPoint(120, 120), zzFrame.End)
	}

	{
		rect1 := Rectangle().Frame(Bind(NewSize(100, 500))).
			Padding(Bind(CGInset{10, 10, 10, 10})).
			Border(Bind(CGInset{10, 10, 10, 10}), Bind(white)).(*stackImpl)
		rect2 := Rectangle().Frame(Bind(NewSize(200, 100))).
			Padding(Bind(CGInset{50, 50, 50, 50})).
			Border(Bind(CGInset{50, 50, 50, 50}), Bind(white)).(*stackImpl)
		z := ZStack(rect1, rect2).(*stackImpl)

		data, layoutFn := z.preload(nil)
		su.Equal(NewSize(400, 540), data.FrameSize)
		su.Equal(false, data.IsInfWidth)
		su.Equal(false, data.IsInfHeight)
		su.Equal(NewInset(0, 0, 0, 0), data.Padding)
		su.Equal(NewInset(0, 0, 0, 0), data.Border)

		final, _, _ := layoutFn(NewPoint(0, 0), NewSize(1000, 1000))
		su.Equal(NewPoint(0, 0), final.Start)
		su.Equal(NewPoint(400, 540), final.End)

		rect1Frame := rect1.viewCtx.systemSetFrame()
		su.Equal(NewPoint(20, 20), rect1Frame.Start)
		su.Equal(NewPoint(120, 520), rect1Frame.End)

		rect2Frame := rect2.viewCtx.systemSetFrame()
		su.Equal(NewPoint(100, 100), rect2Frame.Start)
		su.Equal(NewPoint(300, 200), rect2Frame.End)

	}

	{
		var (
			white  = NewColor(255, 255, 255, 255)
			red    = NewColor(255, 0, 0, 255)
			green  = NewColor(0, 255, 0, 255)
			yellow = NewColor(255, 255, 0, 255)
		)

		rect := Rectangle().
			Frame(Bind(NewSize(100, 100))).
			BackgroundColor(Bind(red)).
			Border(Bind(CGInset{10, 10, 10, 10}), Bind(white)).(*rectangleImpl)

		z1 := rect.Padding(Bind(CGInset{10, 10, 10, 10})).
			BackgroundColor(Bind(green)).
			Border(Bind(CGInset{5, 5, 5, 5}), Bind(yellow)).(*stackImpl)

		rData, _ := rect.preload(nil)
		su.Equal(NewSize(100, 100), rData.FrameSize)
		su.Equal(false, rData.IsInfWidth)
		su.Equal(false, rData.IsInfHeight)
		su.Equal(NewInset(0, 0, 0, 0), rData.Padding)
		su.Equal(NewInset(10, 10, 10, 10), rData.Border)

		z1Data, _ := z1.preload(nil)
		su.Equal(NewSize(120, 120), z1Data.FrameSize)
		su.Equal(false, z1Data.IsInfWidth)
		su.Equal(false, z1Data.IsInfHeight)
		su.Equal(NewInset(10, 10, 10, 10), z1Data.Padding)
		su.Equal(NewInset(5, 5, 5, 5), z1Data.Border)

		root := ZStack(rect, z1).(*stackImpl)

		data, layoutFn := root.preload(nil)
		su.Equal(NewSize(150, 150), data.FrameSize)
		su.Equal(false, data.IsInfWidth)
		su.Equal(false, data.IsInfHeight)
		su.Equal(NewInset(0, 0, 0, 0), data.Padding)
		su.Equal(NewInset(0, 0, 0, 0), data.Border)

		_ = layoutFn
	}
}
