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
	vstack.viewCtx = newViewContext(vstack)
	return vstack
}

func (su *TestVStackSuite) TestVStack() {
	{ // 子視圖全固定大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(NewSize(100, 100))),
			rect2.Frame(Bind(NewSize(200, 200))),
		)

		size, inset, layoutFn := vstack.preload(nil)
		su.Equal(NewSize(200, 300), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(false, size.IsInfY)
		su.Equal(NewInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(NewPoint(0, 0), NewSize(500, 500))
		su.Equal(CGPoint{}, result.Start)
		su.Equal(NewPoint(200, 300), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(NewPoint(0, 0), frameVstack.Start)
		su.Equal(NewPoint(200, 300), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(0, 0), frame1.Start)
		su.Equal(NewPoint(100, 100), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(NewPoint(0, 100), frame2.Start)
		su.Equal(NewPoint(200, 300), frame2.End)
	}

	{ // 子視圖全固定大小 + VStack Padding
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(NewSize(100, 100))),
			rect2.Frame(Bind(NewSize(200, 200))),
		).Padding(Bind(CGInset{15, 15, 15, 15}))

		size, inset, layoutFn := vstack.preload(nil)
		su.Equal(NewSize(200, 300), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(false, size.IsInfY)
		su.Equal(NewInset(15, 15, 15, 15), inset)
		su.NotNil(layoutFn)

		result := layoutFn(NewPoint(0, 0), NewSize(500, 500))
		su.Equal(CGPoint{}, result.Start)
		su.Equal(NewPoint(230, 330), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(NewPoint(15, 15), frameVstack.Start)
		su.Equal(NewPoint(215, 315), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(15, 15), frame1.Start)
		su.Equal(NewPoint(115, 115), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(NewPoint(15, 115), frame2.Start)
		su.Equal(NewPoint(215, 315), frame2.End)
	}

	{ // 子視圖有 X 彈性大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(NewSize(100, 100))),
			rect2.Frame(Bind(NewSize(Inf, 200))),
		)

		size, inset, layoutFn := vstack.preload(nil)
		su.Equal(NewSize(100, 300), size.Frame)
		su.Equal(true, size.IsInfX)
		su.Equal(false, size.IsInfY)
		su.Equal(NewInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(NewPoint(0, 0), NewSize(500, 500))
		su.Equal(CGPoint{}, result.Start)
		su.Equal(NewPoint(500, 300), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(NewPoint(0, 0), frameVstack.Start)
		su.Equal(NewPoint(500, 300), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(0, 0), frame1.Start)
		su.Equal(NewPoint(100, 100), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(NewPoint(0, 100), frame2.Start)
		su.Equal(NewPoint(500, 300), frame2.End)
	}

	{ // 子視圖有 Y 彈性大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(NewSize(100, 100))),
			rect2.Frame(Bind(NewSize(200.0, Inf))),
		)

		size, inset, layoutFn := vstack.preload(nil)
		su.Equal(NewSize(200, 100), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(true, size.IsInfY)
		su.Equal(NewInset(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(NewPoint(0, 0), NewSize(500, 500))
		su.Equal(CGPoint{}, result.Start)
		su.Equal(NewPoint(200, 500), result.End)

		frameVstack := vstack.systemSetFrame()
		su.Equal(NewPoint(0, 0), frameVstack.Start)
		su.Equal(NewPoint(200, 500), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(0, 0), frame1.Start)
		su.Equal(NewPoint(100, 100), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(NewPoint(0, 100), frame2.Start)
		su.Equal(NewPoint(200, 500), frame2.End)
	}

	{ // 子視圖有多個 Y 彈性大小 + VStack Padding
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		rect3 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(NewSize(300, 100))),
			rect2.Frame(Bind(NewSize(100, Inf))),
			rect3.Frame(Bind(NewSize(100, Inf))).Padding(Bind(CGInset{15, 15, 15, 15})),
		).Padding(Bind(CGInset{10, 10, 10, 10}))

		// Y: (500-(10*2)-(100)-(15*2))/2 = (500-20-100-30)/2 = 350/2 = 175
		// FlexY: 175

		size, inset, layoutFn := vstack.preload(nil)
		su.Equal(NewSize(300, 0), size.Frame)
		su.Equal(false, size.IsInfX)
		su.Equal(true, size.IsInfY)
		su.Equal(NewInset(10, 10, 10, 10), inset)
		su.NotNil(layoutFn)

		result := layoutFn(NewPoint(0, 0), NewSize(500.0, 500.0))
		su.Equal(CGPoint{}, result.Start)
		su.Equal(NewPoint(320, 500), result.End)

		vstackFrame := vstack.systemSetFrame()
		su.Equal(NewPoint(10, 10), vstackFrame.Start)
		su.Equal(NewPoint(310, 490), vstackFrame.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(10, 10), frame1.Start)
		su.Equal(NewPoint(310, 110), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(NewPoint(10, 110), frame2.Start)
		su.Equal(NewPoint(110, 285), frame2.End)

		frame3 := rect3.systemSetFrame()
		su.Equal(NewPoint(25, 300.0), frame3.Start)
		su.Equal(NewPoint(125.0, 445.0), frame3.End)
	}

	{
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		rect3 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1,
			rect2.Frame(Bind(NewSize(100, 100))),
			rect3,
		)

		_, _, layoutFn := vstack.preload(nil)
		layoutFn(NewPoint(0, 0), NewSize(500, 500))

		frameVstack := vstack.systemSetFrame()
		su.Equal(NewPoint(0, 0), frameVstack.Start)
		su.Equal(NewPoint(500, 500), frameVstack.End)

		frame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(0, 0), frame1.Start)
		su.Equal(NewPoint(500, 200), frame1.End)

		frame2 := rect2.systemSetFrame()
		su.Equal(NewPoint(0, 200), frame2.Start)
		su.Equal(NewPoint(100, 300), frame2.End)

		frame3 := rect3.systemSetFrame()
		su.Equal(NewPoint(0, 300), frame3.Start)
		su.Equal(NewPoint(500, 500), frame3.End)
	}
}
