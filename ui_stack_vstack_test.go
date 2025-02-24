package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestVStack(t *testing.T) {
	suite.Run(t, new(TestVStackSuite))
}

type TestVStackSuite struct {
	suite.Suite
}

func newVStackForTest(children ...SomeView) *vstackImpl {
	vstack := &vstackImpl{
		children: children,
	}
	vstack.ctx = newViewContext(vstack)
	return vstack
}

func (su *TestVStackSuite) TestVStack() {
	{ // 子視圖全固定大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(Bind(200.0), Bind(200.0)),
		)

		size, inset, layoutFn := vstack.preload()
		su.Equal(sz(200, 300), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(false, size.IsInfY)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500, 500))
		su.Equal(ptZero, result.Start)
		su.Equal(pt(200, 300), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(pt(0, 0), frameVstack.Start)
		su.Equal(pt(200, 300), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(pt(0, 0), frame1.Start)
		su.Equal(pt(100, 100), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(pt(0, 100), frame2.Start)
		su.Equal(pt(200, 300), frame2.End)
	}

	{ // 子視圖全固定大小 + VStack Padding
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(Bind(200.0), Bind(200.0)),
		).Padding(Bind(15.0))

		size, inset, layoutFn := vstack.preload()
		su.Equal(sz(200, 300), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(false, size.IsInfY)
		su.Equal(ins(15, 15, 15, 15), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500, 500))
		su.Equal(ptZero, result.Start)
		su.Equal(pt(230, 330), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(pt(15, 15), frameVstack.Start)
		su.Equal(pt(215, 315), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(pt(15, 15), frame1.Start)
		su.Equal(pt(115, 115), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(pt(15, 115), frame2.Start)
		su.Equal(pt(215, 315), frame2.End)
	}

	{ // 子視圖有 X 彈性大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(nil, Bind(200.0)),
		)

		size, inset, layoutFn := vstack.preload()
		su.Equal(sz(100, 300), size.Frame)
		su.Equal(true, size.IsInfX)
		su.Equal(false, size.IsInfY)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500, 500))
		su.Equal(ptZero, result.Start)
		su.Equal(pt(500, 300), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(pt(0, 0), frameVstack.Start)
		su.Equal(pt(500, 300), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(pt(0, 0), frame1.Start)
		su.Equal(pt(100, 100), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(pt(0, 100), frame2.Start)
		su.Equal(pt(500, 300), frame2.End)
	}

	{ // 子視圖有 Y 彈性大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(Bind(200.0), nil),
		)

		size, inset, layoutFn := vstack.preload()
		su.Equal(sz(200, 100), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(true, size.IsInfY)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500, 500))
		su.Equal(ptZero, result.Start)
		su.Equal(pt(200, 500), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(pt(0, 0), frameVstack.Start)
		su.Equal(pt(200, 500), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(pt(0, 0), frame1.Start)
		su.Equal(pt(100, 100), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(pt(0, 100), frame2.Start)
		su.Equal(pt(200, 500), frame2.End)
	}

	{ // 子視圖有多個 Y 彈性大小 + VStack Padding
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		rect3 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(300.0), Bind(100.0)),
			rect2.Frame(Bind(100.0), nil),
			rect3.Frame(Bind(100.0), nil).Padding(Bind(15.0)),
		).Padding(Bind(10.0))

		// Y: (500-(10*2)-(100)-(15*2))/2 = (500-20-100-30)/2 = 350/2 = 175
		// FlexY: 175

		size, inset, layoutFn := vstack.preload()
		su.Equal(sz(300, 130), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(true, size.IsInfY)
		su.Equal(ins(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500.0, 500.0))
		su.Equal(ptZero, result.Start)
		su.Equal(pt(320, 500), result.End)

		vstackFrame := vstack.systemSetFrame()
		su.Equal(pt(10, 10), vstackFrame.Start)
		su.Equal(pt(310, 490), vstackFrame.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(pt(10, 10), frame1.Start)
		su.Equal(pt(310, 110), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(pt(10, 110), frame2.Start)
		su.Equal(pt(110, 285), frame2.End)

		frame3 := rect3.systemSetFrame()
		su.Equal(pt(25, 300.0), frame3.Start)
		su.Equal(pt(125.0, 475.0), frame3.End)
	}

	{
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		rect3 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1,
			rect2.Frame(Bind(100.0), Bind(100.0)),
			rect3,
		)

		_, _, layoutFn := vstack.preload()
		layoutFn(pt(0, 0), sz(500, 500))

		frameVstack := vstack.systemSetFrame()
		su.Equal(pt(0, 0), frameVstack.Start)
		su.Equal(pt(500, 500), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(pt(0, 0), frame1.Start)
		su.Equal(pt(500, 200), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(pt(0, 200), frame2.Start)
		su.Equal(pt(100, 300), frame2.End)

		frame3 := rect3.systemSetFrame()
		su.Equal(pt(0, 300), frame3.Start)
		su.Equal(pt(500, 500), frame3.End)
	}
}
