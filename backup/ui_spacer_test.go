package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSpacer(t *testing.T) {
	suite.Run(t, new(SpacerSuite))
}

type SpacerSuite struct {
	suite.Suite
}

func (su *SpacerSuite) Test() {
	{
		sp1 := Spacer()
		sp2 := Spacer()
		sp3 := Spacer()
		sp4 := Spacer()
		rect1 := Rectangle().(*rectangleImpl)
		rect1.Frame(Bind(NewSize(100, 100)))

		view := ZStack(
			VStack(
				sp1,
				HStack(
					sp2,
					rect1,
					sp3,
				),
				sp4,
			),
		)
		_, layoutFn := view.preload(nil)
		bound, _, _ := layoutFn(NewPoint(0, 0), NewSize(500, 500))
		su.Equal(NewPoint(0, 0), bound.Start)
		su.Equal(NewPoint(500, 500), bound.End)

		rectFrame1 := rect1.systemSetFrame()
		su.Equal(NewPoint(200, 200), rectFrame1.Start)
		su.Equal(NewPoint(300, 300), rectFrame1.End)
	}
}
