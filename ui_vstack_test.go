package ebui

import (
	"testing"
)

// import (
// 	"fmt"
// 	"testing"
// )

// func TestVStackBasic(t *testing.T) {
// 	v3 := VStack()
// 	v2 := VStack()
// 	v1 := VStack(v2, v3)

// 	r := root(v1, size{200, 100})
// 	r.calculateStage()

// 	{
// 		opt := v1.view()
// 		msg := fmt.Sprintf("v1: %+v", opt)
// 		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
// 		assert(t, opt.size.w, 200, "opt.size.w", msg)
// 		assert(t, opt.size.h, 100, "opt.size.h", msg)
// 		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
// 		assert(t, opt.xx, 0, "opt.xx", msg)
// 		assert(t, opt.yy, 0, "opt.yy", msg)
// 	}

// 	{
// 		opt := v2.view()
// 		msg := fmt.Sprintf("v2: %+v", opt)
// 		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
// 		assert(t, opt.size.w, 200, "opt.size.w", msg)
// 		assert(t, opt.size.h, 50, "opt.size.h", msg)
// 		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
// 		assert(t, opt.xx, 0, "opt.xx", msg)
// 		assert(t, opt.yy, 0, "opt.yy", msg)
// 	}

// 	{
// 		opt := v3.view()
// 		msg := fmt.Sprintf("v3: %+v", opt)
// 		assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
// 		assert(t, opt.size.w, 200, "opt.size.w", msg)
// 		assert(t, opt.size.h, 50, "opt.size.h", msg)
// 		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 50, "opt.pos.y", msg)
// 		assert(t, opt.xx, 0, "opt.xx", msg)
// 		assert(t, opt.yy, 50, "opt.yy", msg)
// 	}
// }

func TestVStackGetSize(t *testing.T) {
	{
		v := VStack()
		size := v.getSize()
		assert(t, size.w, -1, "size.w")
		assert(t, size.h, -1, "size.h")
	}

	{
		v := VStack().Frame(30, 10)
		size := v.getSize()
		assert(t, size.w, 30, "size.w")
		assert(t, size.h, 10, "size.h")
	}

	{
		v := VStack(
			VStack().Frame(60, 20),
		)
		size := v.getSize()
		assert(t, size.w, 60, "size.w")
		assert(t, size.h, 20, "size.h")
	}

	{
		v := VStack(
			VStack().Frame(60, 20),
			VStack().Frame(40, 80),
		)
		size := v.getSize()
		assert(t, size.w, 60, "size.w")
		assert(t, size.h, 100, "size.h")
	}

	{
		v := VStack(
			VStack().Frame(60, 20),
			VStack().Frame(40, 80),
		).Frame(30, 10)
		size := v.getSize()
		assert(t, size.w, 30, "size.w")
		assert(t, size.h, 10, "size.h")
	}
}
