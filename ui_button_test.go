package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestButton(t *testing.T) {
	suite.Run(t, new(ButtonSuite))
}

type ButtonSuite struct {
	suite.Suite
}

func (su *ButtonSuite) TestButton() {

	{ // button with no size, button label with size
		rect := Rectangle().
			Padding(Bind(NewInset(10, 10, 10, 10))).
			Frame(Bind(NewSize(200, 100)))
		btn := Button("", func() {}, func() SomeView {
			return rect
		})
		btn.Padding(Bind(NewInset(20, 20, 20, 20)))

		btnData, btnLayoutFn := btn.preload(nil)
		su.Equal(NewSize(220, 120), btnData.FrameSize)
		su.Equal(NewInset(0, 0, 0, 0), btnData.Padding)
		su.NotNil(btnLayoutFn)
	}
}
