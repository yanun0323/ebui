package ebui

import (
	"fmt"
	"testing"
)

func TestPreCacheChildrenSizeSmallerChildren(t *testing.T) {
	v2 := VStack().Frame(20, 20)
	v1 := VStack(v2).Frame(50, 50)
	r := &rootView{}
	r.view = newView(typeRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v2.params()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initW, 20, "p.initW", msg)
		assert(t, p.initH, 20, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	{
		p := v1.params()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initW, 50, "p.initW", msg)
		assert(t, p.initH, 50, "p.initH", msg)
		assert(t, p.w, 50, "p.w", msg)
		assert(t, p.h, 50, "p.h", msg)
	}
}

func TestPreCacheChildrenSizeBiggerChildren(t *testing.T) {
	v2 := VStack().Frame(70, 70)
	v1 := VStack(v2).Frame(50, 50)
	r := &rootView{}
	r.view = newView(typeRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v2.params()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initW, 70, "p.initW", msg)
		assert(t, p.initH, 70, "p.initH", msg)
		assert(t, p.w, 70, "p.w", msg)
		assert(t, p.h, 70, "p.h", msg)
	}

	{
		p := v1.params()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initW, 50, "p.initW", msg)
		assert(t, p.initH, 50, "p.initH", msg)
		assert(t, p.w, 70, "p.w", msg)
		assert(t, p.h, 70, "p.h", msg)
	}

	{
		p := r.params()
		msg := fmt.Sprintf("r: %+v", p)
		assert(t, p.initW, 100, "p.initW", msg)
		assert(t, p.initH, 100, "p.initH", msg)
		assert(t, p.w, 100, "p.w", msg)
		assert(t, p.h, 100, "p.h", msg)
	}
}

func TestCalculateParametersSingle(t *testing.T) {
	v2 := VStack().Frame(20, 20)
	v1 := VStack(v2)
	r := &rootView{}
	r.view = newView(typeRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v2.params()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initW, 20, "p.initW", msg)
		assert(t, p.initH, 20, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	{
		p := v1.params()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initW, -1, "p.initW", msg)
		assert(t, p.initH, -1, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	r.calculationParameters()

	{
		p := v2.params()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initW, 20, "p.initW", msg)
		assert(t, p.initH, 20, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	{
		p := v1.params()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initW, -1, "p.initW", msg)
		assert(t, p.initH, -1, "p.initH", msg)
		assert(t, p.w, 100, "p.w", msg)
		assert(t, p.h, 100, "p.h", msg)
	}
}

func TestCalculateParametersMultiple(t *testing.T) {
	v3 := HStack()
	v2 := VStack().Frame(20, 20)
	v1 := VStack(v2, v3)
	r := &rootView{}
	r.view = newView(typeRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v3.params()
		msg := fmt.Sprintf("v3: %+v", p)
		assert(t, p.initW, -1, "p.initW", msg)
		assert(t, p.initH, -1, "p.initH", msg)
		assert(t, p.w, -1, "p.w", msg)
		assert(t, p.h, -1, "p.h", msg)
	}

	{
		p := v2.params()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initW, 20, "p.initW", msg)
		assert(t, p.initH, 20, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	{
		p := v1.params()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initW, -1, "p.initW", msg)
		assert(t, p.initH, -1, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	r.calculationParameters()

	{
		p := v3.params()
		msg := fmt.Sprintf("v3: %+v", p)
		assert(t, p.initW, -1, "p.initW", msg)
		assert(t, p.initH, -1, "p.initH", msg)
		assert(t, p.w, 100, "p.w", msg)
		assert(t, p.h, 80, "p.h", msg)
	}

	{
		p := v2.params()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initW, 20, "p.initW", msg)
		assert(t, p.initH, 20, "p.initH", msg)
		assert(t, p.w, 20, "p.w", msg)
		assert(t, p.h, 20, "p.h", msg)
	}

	{
		p := v1.params()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initW, -1, "p.initW", msg)
		assert(t, p.initH, -1, "p.initH", msg)
		assert(t, p.w, 100, "p.w", msg)
		assert(t, p.h, 100, "p.h", msg)
	}
}
