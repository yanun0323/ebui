package ebui

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func newButtonForTest(action func(), label func() SomeView) *buttonImpl {
	btn := &buttonImpl{
		action: action,
		label:  label,
	}
	btn.ctx = newViewContext(btn)
	return btn
}

func TestButton(t *testing.T) {
	suite.Run(t, new(ButtonSuite))
}

type ButtonSuite struct {
	suite.Suite
}

func (su *ButtonSuite) TestButton() {

	{ // 按鈕無大小，按鈕標籤有大小
		rect := newRectangleForTest().
			Frame(Bind(200.0), Bind(100.0)).
			Padding(Bind(CGInset(10, 10, 10, 10)))
		btn := newButtonForTest(func() {}, func() SomeView {
			return rect
		})
		btn.Padding(Bind(CGInset(20, 20, 20, 20)))

		btnFrameSize, btnInset, btnLayoutFn := btn.preload()
		su.Equal(CGSize(220, 120), btnFrameSize.Frame)
		su.Equal(CGInset(20, 20, 20, 20), btnInset)
		su.NotNil(btnLayoutFn)
	}
}
