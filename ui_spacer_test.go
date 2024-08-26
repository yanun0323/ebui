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
		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
		assert(t, opt.size.w, 200, "opt.size.w", msg)
		assert(t, opt.size.h, 25, "opt.size.h", msg)
		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := sp2.params()
		msg := fmt.Sprintf("sp2: %+v", opt)
		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
		assert(t, opt.size.w, 200, "opt.size.w", msg)
		assert(t, opt.size.h, 25, "opt.size.h", msg)
		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
		assert(t, opt.pos.y, 50, "opt.pos.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 50, "opt.yy", msg)
	}

	{
		opt := sp3.params()
		msg := fmt.Sprintf("sp3: %+v", opt)
		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
		assert(t, opt.size.w, 200, "opt.size.w", msg)
		assert(t, opt.size.h, 5, "opt.size.h", msg)
		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
		assert(t, opt.pos.y, 45, "opt.pos.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 20, "opt.yy", msg)
	}
}
