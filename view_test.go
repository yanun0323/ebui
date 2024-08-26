package ebui

import (
	"fmt"
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

func TestViewCalculateVStack(t *testing.T) {
	v3 := VStack().Frame(20, 20) /* 20, 20, 0, 80 */
	v2 := VStack(
		Spacer(),
		v3,
	) /* 200, 50, 0, 50 */
	v1 := VStack(
		Spacer(),
		v2,
	).Frame(200, 100) /*  200, 100, 0, 0 */

	r := root(v1)
	r.setWindowSize(200, 100)
	r.calculateStage()

	{
		opt := v1.params()
		msg := fmt.Sprintf("v1: %+v", opt)
		assert(t, opt.initW, 200, "opt.initW", msg)
		assert(t, opt.initH, 100, "opt.initH", msg)
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
		assert(t, opt.w, 200, "opt.w", msg)
		assert(t, opt.h, 50, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 50, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 50, "opt.yy", msg)
	}

	{
		opt := v3.params()
		msg := fmt.Sprintf("v3: %+v", opt)
		assert(t, opt.initW, 20, "opt.initW", msg)
		assert(t, opt.initH, 20, "opt.initH", msg)
		assert(t, opt.w, 20, "opt.w", msg)
		assert(t, opt.h, 20, "opt.h", msg)
		assert(t, opt.x, 0, "opt.x", msg)
		assert(t, opt.y, 80, "opt.y", msg)
		assert(t, opt.xx, 0, "opt.xx", msg)
		assert(t, opt.yy, 30, "opt.yy", msg)
	}
}

func TestViewCalculateHStack(t *testing.T) {
	v3 := HStack().Frame(20, 20) /* 20, 20, 0, 0 */
	v2 := HStack(
		Spacer(), /* 80, 100, 100, 0 */
		v3,       /* 20, 20, 180, 0 */
	)
	v1 := HStack(
		Spacer(), /*  100, 100, 0, 0 */
		v2,       /*  100, 100, 100, 0 */
	).Frame(200, 100)

	r := root(v1)
	r.setWindowSize(200, 100)
	r.calculateStage()

	{
		opt := v1.params()
		msg := fmt.Sprintf("v1: %+v", opt)
		assert(t, opt.initW, 200, "opt.initW", msg)
		assert(t, opt.initH, 100, "opt.initH", msg)
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
		assert(t, opt.x, 100, "opt.x", msg)
		assert(t, opt.y, 0, "opt.y", msg)
		assert(t, opt.xx, 100, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}

	{
		opt := v3.params()
		msg := fmt.Sprintf("v3: %+v", opt)
		assert(t, opt.initW, 20, "opt.initW", msg)
		assert(t, opt.initH, 20, "opt.initH", msg)
		assert(t, opt.w, 20, "opt.w", msg)
		assert(t, opt.h, 20, "opt.h", msg)
		assert(t, opt.x, 180, "opt.x", msg)
		assert(t, opt.y, 0, "opt.y", msg)
		assert(t, opt.xx, 80, "opt.xx", msg)
		assert(t, opt.yy, 0, "opt.yy", msg)
	}
}
