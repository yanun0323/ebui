package ebui

import (
	"fmt"
	"testing"
)

func TestSpacerVStackBasic(t *testing.T) {
	sp1 := Spacer()
	sp2 := Spacer()
	sp3 := Spacer()
	r := root(
		VStack(
			sp1,
			VStack( /* w:200, h:25 */
				HStack().Frame(20, 20),
				sp3,
			),
			sp2,
			Spacer(),
		),
		size{200, 100},
	)
	r.calculateStage()

	{
		opt := sp1.params()
		msg := fmt.Sprintf("sp1: %+v", opt)
		assert(t, opt.initW, -1, "opt.initW", msg)
		assert(t, opt.initH, -1, "opt.initH", msg)
		assert(t, opt.w, 200, "opt.w", msg)
		assert(t, opt.h, 25, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 0, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := sp2.params()
		msg := fmt.Sprintf("sp2: %+v", opt)
		assert(t, opt.initW, -1, "opt.initW", msg)
		assert(t, opt.initH, -1, "opt.initH", msg)
		assert(t, opt.w, 200, "opt.w", msg)
		assert(t, opt.h, 25, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 50, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 50, "opt.yy", msg)
	}

	{
		opt := sp3.params()
		msg := fmt.Sprintf("sp3: %+v", opt)
		assert(t, opt.initW, -1, "opt.initW", msg)
		assert(t, opt.initH, -1, "opt.initH", msg)
		assert(t, opt.w, 200, "opt.w", msg)
		assert(t, opt.h, 5, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 45, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 20, "opt.yy", msg)
	}
}
