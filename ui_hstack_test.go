package ebui

import "testing"

// import (
// 	"fmt"
// 	"testing"
// )

// func TestHStackBasic(t *testing.T) {
// 	v3 := HStack()
// 	v2 := HStack()
// 	v1 := HStack(v2, v3)

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
// 		assert(t, opt.size.w, 100, "opt.size.w", msg)
// 		assert(t, opt.size.h, 100, "opt.size.h", msg)
// 		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
// 		assert(t, opt.xx, 0, "opt.xx", msg)
// 		assert(t, opt.yy, 0, "opt.yy", msg)
// 	}

//		{
//			opt := v3.view()
//			msg := fmt.Sprintf("v3: %+v", opt)
//			assert(t, opt.initSize.w, -1, "opt.initSize.w", msg)
//			assert(t, opt.initSize.h, -1, "opt.initSize.h", msg)
//			assert(t, opt.size.w, 100, "opt.size.w", msg)
//			assert(t, opt.size.h, 100, "opt.size.h", msg)
//			assert(t, opt.pos.x, 100, "opt.pos.x", msg)
//			assert(t, opt.pos.y, 0, "opt.pos.y", msg)
//			assert(t, opt.xx, 100, "opt.xx", msg)
//			assert(t, opt.yy, 0, "opt.yy", msg)
//		}
//	}

func TestHStackGetSize(t *testing.T) {
	{
		v := HStack()
		size := v.getSize()
		assert(t, size.w, -1, "size.w")
		assert(t, size.h, -1, "size.h")
	}

	{
		v := HStack().Frame(30, 10)
		size := v.getSize()
		assert(t, size.w, 30, "size.w")
		assert(t, size.h, 10, "size.h")
	}

	{
		v := HStack(
			HStack().Frame(60, 20),
		)
		size := v.getSize()
		assert(t, size.w, 60, "size.w")
		assert(t, size.h, 20, "size.h")
	}

	{
		v := HStack(
			HStack().Frame(60, 20),
			HStack().Frame(40, 80),
		)
		size := v.getSize()
		assert(t, size.w, 100, "size.w")
		assert(t, size.h, 80, "size.h")
	}

	{
		v := HStack(
			HStack().Frame(60, 20),
			HStack().Frame(40, 80),
		).Frame(30, 10)
		size := v.getSize()
		assert(t, size.w, 30, "size.w")
		assert(t, size.h, 10, "size.h")
	}
}
