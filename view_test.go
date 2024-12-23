package ebui

import (
	"strings"
	"testing"
)

func assert[T comparable](t *testing.T, result, expected T, msg string, messages ...string) {
	t.Helper()
	if expected != result {
		switch len(messages) {
		case 0:
			t.Fatalf("%s: expected %v, got %v", msg, expected, result)
		default:
			t.Fatalf("%s: expected %v, got %v\n %s", msg, expected, result, strings.Join(messages, "\n"))
		}
	}
}

// func TestViewCalculateVStack(t *testing.T) {
// 	v3 := VStack().Frame(20, 20) /* 20, 20, 0, 80 */
// 	v2 := VStack(
// 		Spacer(),
// 		v3,
// 	) /* 200, 50, 0, 50 */
// 	v1 := VStack(
// 		Spacer(),
// 		v2,
// 	).Frame(200, 100) /*  200, 100, 0, 0 */

// 	r := root(v1)
// 	r.setWindowSize(200, 100)
// 	r.calculateStage()

// 	{
// 		opt := v1.view()
// 		msg := fmt.Sprintf("v1: %+v", opt)
// 		assert(t, opt.initSize.w, 200, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, 100, "opt.initSize.h", msg)
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
// 		assert(t, opt.pos.y, 50, "opt.pos.y", msg)
// 		assert(t, opt.xx, 0, "opt.xx", msg)
// 		assert(t, opt.yy, 50, "opt.yy", msg)
// 	}

// 	{
// 		opt := v3.view()
// 		msg := fmt.Sprintf("v3: %+v", opt)
// 		assert(t, opt.initSize.w, 20, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, 20, "opt.initSize.h", msg)
// 		assert(t, opt.size.w, 20, "opt.size.w", msg)
// 		assert(t, opt.size.h, 20, "opt.size.h", msg)
// 		assert(t, opt.pos.x, 0, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 80, "opt.pos.y", msg)
// 		assert(t, opt.xx, 0, "opt.xx", msg)
// 		assert(t, opt.yy, 30, "opt.yy", msg)
// 	}
// }

// func TestViewCalculateHStack(t *testing.T) {
// 	v3 := HStack().Frame(20, 20) /* 20, 20, 0, 0 */
// 	v2 := HStack(
// 		Spacer(), /* 80, 100, 100, 0 */
// 		v3,       /* 20, 20, 180, 0 */
// 	)
// 	v1 := HStack(
// 		Spacer(), /*  100, 100, 0, 0 */
// 		v2,       /*  100, 100, 100, 0 */
// 	).Frame(200, 100)

// 	r := root(v1)
// 	r.setWindowSize(200, 100)
// 	r.calculateStage()

// 	{
// 		opt := v1.view()
// 		msg := fmt.Sprintf("v1: %+v", opt)
// 		assert(t, opt.initSize.w, 200, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, 100, "opt.initSize.h", msg)
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
// 		assert(t, opt.pos.x, 100, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
// 		assert(t, opt.xx, 100, "opt.xx", msg)
// 		assert(t, opt.yy, 0, "opt.yy", msg)
// 	}

// 	{
// 		opt := v3.view()
// 		msg := fmt.Sprintf("v3: %+v", opt)
// 		assert(t, opt.initSize.w, 20, "opt.initSize.w", msg)
// 		assert(t, opt.initSize.h, 20, "opt.initSize.h", msg)
// 		assert(t, opt.size.w, 20, "opt.size.w", msg)
// 		assert(t, opt.size.h, 20, "opt.size.h", msg)
// 		assert(t, opt.pos.x, 180, "opt.pos.x", msg)
// 		assert(t, opt.pos.y, 0, "opt.pos.y", msg)
// 		assert(t, opt.xx, 80, "opt.xx", msg)
// 		assert(t, opt.yy, 0, "opt.yy", msg)
// 	}
// }
