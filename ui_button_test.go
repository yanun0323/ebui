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
	btn.viewCtx = newViewContext(btn)
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
			Padding(Bind(NewInset(10, 10, 10, 10))).
			Frame(Bind(NewSize(200, 100)))
		btn := newButtonForTest(func() {}, func() SomeView {
			return rect
		})
		btn.Padding(Bind(NewInset(20, 20, 20, 20)))

		btnData, btnLayoutFn := btn.preload(nil)
		su.Equal(NewSize(220, 120), btnData.FrameSize)
		su.Equal(NewInset(0, 0, 0, 0), btnData.Padding)
		su.NotNil(btnLayoutFn)
	}
}
