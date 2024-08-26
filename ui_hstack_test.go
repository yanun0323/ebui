package ebui

import (
	"fmt"
	"testing"
)

func TestHStackBasic(t *testing.T) {
	v3 := HStack()
	v2 := HStack()
	v1 := HStack(v2, v3)

	r := root(v1, frame{200, 100})
	r.calculateStage()

	{
		opt := v1.view()
		msg := fmt.Sprintf("v1: %+v", opt)
		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
		assert(t, opt.size.w, 200, "opt.size.w", msg)
		assert(t, opt.size.h, 100, "opt.size.h", msg)
		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := v2.view()
		msg := fmt.Sprintf("v2: %+v", opt)
		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
		assert(t, opt.size.w, 100, "opt.size.w", msg)
		assert(t, opt.size.h, 100, "opt.size.h", msg)
		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := v3.view()
		msg := fmt.Sprintf("v3: %+v", opt)
		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
		assert(t, opt.size.w, 100, "opt.size.w", msg)
		assert(t, opt.size.h, 100, "opt.size.h", msg)
		assert(t, opt.pos.x, 100, "opt.pos.x", msg)
		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
		assert(t, opt.xx, 100, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}
}
