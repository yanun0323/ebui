package ebui

import (
	"fmt"
	"testing"
)

func TestHStackBasic(t *testing.T) {
	v3 := HStack()
	v2 := HStack()
	v1 := HStack(v2, v3)

	r := root(v1, size{200, 100})
	r.calculateStage()

	{
		opt := v1.params()
		msg := fmt.Sprintf("v1: %+v", opt)
		assert(t, opt.initW, -1, "opt.initW", msg)
		assert(t, opt.initH, -1, "opt.initH", msg)
		assert(t, opt.w, 200, "opt.w", msg)
		assert(t, opt.h, 100, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 0, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := v2.params()
		msg := fmt.Sprintf("v2: %+v", opt)
		assert(t, opt.initW, -1, "opt.initW", msg)
		assert(t, opt.initH, -1, "opt.initH", msg)
		assert(t, opt.w, 100, "opt.w", msg)
		assert(t, opt.h, 100, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 0, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := v3.params()
		msg := fmt.Sprintf("v3: %+v", opt)
		assert(t, opt.initW, -1, "opt.initW", msg)
		assert(t, opt.initH, -1, "opt.initH", msg)
		assert(t, opt.w, 100, "opt.w", msg)
		assert(t, opt.h, 100, "opt.h", msg)
		assert(t, opt.x, 100, "opt.x", msg)
		assert(t, opt.y, 0, "opt.y", msg)
		assert(t, opt.xx, 100, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}
}
