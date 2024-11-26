package ebui

import "testing"

func TestZStackGetSize(t *testing.T) {
	{
		v := ZStack()
		size := v.getSize()
		assert(t, size.w, -1, "size.w")
		assert(t, size.h, -1, "size.h")
	}

	{
		v := ZStack().Frame(30, 10)
		size := v.getSize()
		assert(t, size.w, 30, "size.w")
		assert(t, size.h, 10, "size.h")
	}

	{
		v := ZStack(
			ZStack().Frame(60, 20),
		)
		size := v.getSize()
		assert(t, size.w, 60, "size.w")
		assert(t, size.h, 20, "size.h")
	}

	{
		v := ZStack(
			ZStack().Frame(60, 20),
			ZStack().Frame(40, 80),
		)
		size := v.getSize()
		assert(t, size.w, 60, "size.w")
		assert(t, size.h, 80, "size.h")
	}

	{
		v := ZStack(
			ZStack().Frame(60, 20),
			ZStack().Frame(40, 80),
		).Frame(30, 10)
		size := v.getSize()
		assert(t, size.w, 30, "size.w")
		assert(t, size.h, 10, "size.h")
	}
}
