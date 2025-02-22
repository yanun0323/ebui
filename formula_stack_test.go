package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTestFormulaStack(t *testing.T) {
	suite.Run(t, new(TestFormulaStackSuite))
}

type TestFormulaStackSuite struct {
	suite.Suite
}

func newVStackForTest(children ...SomeView) *vstackImpl {
	vstack := &vstackImpl{
		children: children,
	}
	vstack.ctx = newViewContext(tagVStack, vstack)
	return vstack
}

func (su *TestFormulaStackSuite) TestVStack() {
	{ // 子視圖全固定大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(Bind(200.0), Bind(200.0)),
		)

		size, inset, layoutFn := vstack.preload()
		su.Equal(200.0, size.Width)
		su.Equal(300.0, size.Height)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500.0, 500.0))
		su.Equal(0.0, result.Start.X)
		su.Equal(0.0, result.Start.Y)
		su.Equal(200.0, result.End.X)
		su.Equal(300.0, result.End.Y)
	}

	{ // 子視圖全固定大小 + VStack Padding
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(Bind(200.0), Bind(200.0)),
		).Padding(Bind(15.0))

		size, inset, layoutFn := vstack.preload()
		su.Equal(200.0, size.Width)
		su.Equal(300.0, size.Height)
		su.Equal(ins(15, 15, 15, 15), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500.0, 500.0))
		su.Equal(0.0, result.Start.X)
		su.Equal(0.0, result.Start.Y)
		su.Equal(200.0, result.End.X)
		su.Equal(300.0, result.End.Y)
	}

	{ // 子視圖有 X 彈性大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(nil, Bind(100.0)),
		)

		size, inset, layoutFn := vstack.preload()
		su.Equal(Inf, size.Width)
		su.Equal(200.0, size.Height)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500.0, 500.0))
		su.Equal(0.0, result.Start.X)
		su.Equal(0.0, result.Start.Y)
		su.Equal(500.0, result.End.X)
		su.Equal(200.0, result.End.Y)
	}

	{ // 子視圖有 Y 彈性大小
		rect1 := newRectangleForTest()
		rect2 := newRectangleForTest()
		vstack := newVStackForTest(
			rect1.Frame(Bind(100.0), Bind(100.0)),
			rect2.Frame(Bind(100.0), nil),
		)

		size, inset, layoutFn := vstack.preload()
		su.Equal(100.0, size.Width)
		su.Equal(Inf, size.Height)
		su.Equal(ins(0, 0, 0, 0), inset)
		su.NotNil(layoutFn)

		result := layoutFn(pt(0, 0), sz(500.0, 500.0))
		su.Equal(0.0, result.Start.X)
		su.Equal(0.0, result.Start.Y)
		su.Equal(100.0, result.End.X)
		su.Equal(500.0, result.End.Y)
	}
}
