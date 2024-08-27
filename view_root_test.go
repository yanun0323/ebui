package ebui

import (
	"fmt"
	"testing"
)

func TestPreCacheChildrenSizeSmallerChildren(t *testing.T) {
	v2 := VStack().Frame(20, 20)
	v1 := VStack(v2).Frame(50, 50)
	r := &rootView{}
	r.uiView = newUIView(typesRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v2.view()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initSize.w, 20, "p.initSize.w", msg)
		assert(t, p.initSize.h, 20, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	{
		p := v1.view()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initSize.w, 50, "p.initSize.w", msg)
		assert(t, p.initSize.h, 50, "p.initSize.h", msg)
		assert(t, p.size.w, 50, "p.size.w", msg)
		assert(t, p.size.h, 50, "p.size.h", msg)
	}
}

func TestPreCacheChildrenSizeBiggerChildren(t *testing.T) {
	v2 := VStack().Frame(70, 70)
	v1 := VStack(v2).Frame(50, 50)
	r := &rootView{}
	r.uiView = newUIView(typesRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v2.view()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initSize.w, 70, "p.initSize.w", msg)
		assert(t, p.initSize.h, 70, "p.initSize.h", msg)
		assert(t, p.size.w, 70, "p.size.w", msg)
		assert(t, p.size.h, 70, "p.size.h", msg)
	}

	{
		p := v1.view()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initSize.w, 50, "p.initSize.w", msg)
		assert(t, p.initSize.h, 50, "p.initSize.h", msg)
		assert(t, p.size.w, 70, "p.size.w", msg)
		assert(t, p.size.h, 70, "p.size.h", msg)
	}

	{
		p := r.view()
		msg := fmt.Sprintf("r: %+v", p)
		assert(t, p.initSize.w, 100, "p.initSize.w", msg)
		assert(t, p.initSize.h, 100, "p.initSize.h", msg)
		assert(t, p.size.w, 100, "p.size.w", msg)
		assert(t, p.size.h, 100, "p.size.h", msg)
	}
}

func TestCalculateParametersSingle(t *testing.T) {
	v2 := VStack().Frame(20, 20)
	v1 := VStack(v2)
	r := &rootView{}
	r.uiView = newUIView(typesRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v2.view()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initSize.w, 20, "p.initSize.w", msg)
		assert(t, p.initSize.h, 20, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	{
		p := v1.view()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initSize.w, -1, "p.initSize.w", msg)
		assert(t, p.initSize.h, -1, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	r.calculationParameters()

	{
		p := v2.view()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initSize.w, 20, "p.initSize.w", msg)
		assert(t, p.initSize.h, 20, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	{
		p := v1.view()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initSize.w, -1, "p.initSize.w", msg)
		assert(t, p.initSize.h, -1, "p.initSize.h", msg)
		assert(t, p.size.w, 100, "p.size.w", msg)
		assert(t, p.size.h, 100, "p.size.h", msg)
	}
}

func TestCalculateParametersMultiple(t *testing.T) {
	v3 := HStack()
	v2 := VStack().Frame(20, 20)
	v1 := VStack(v2, v3)
	r := &rootView{}
	r.uiView = newUIView(typesRoot, r, v1)
	r.setWindowSize(100, 100)

	r.preCacheChildrenSize()

	{
		p := v3.view()
		msg := fmt.Sprintf("v3: %+v", p)
		assert(t, p.initSize.w, -1, "p.initSize.w", msg)
		assert(t, p.initSize.h, -1, "p.initSize.h", msg)
		assert(t, p.size.w, -1, "p.size.w", msg)
		assert(t, p.size.h, -1, "p.size.h", msg)
	}

	{
		p := v2.view()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initSize.w, 20, "p.initSize.w", msg)
		assert(t, p.initSize.h, 20, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	{
		p := v1.view()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initSize.w, -1, "p.initSize.w", msg)
		assert(t, p.initSize.h, -1, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	r.calculationParameters()

	{
		p := v3.view()
		msg := fmt.Sprintf("v3: %+v", p)
		assert(t, p.initSize.w, -1, "p.initSize.w", msg)
		assert(t, p.initSize.h, -1, "p.initSize.h", msg)
		assert(t, p.size.w, 100, "p.size.w", msg)
		assert(t, p.size.h, 80, "p.size.h", msg)
	}

	{
		p := v2.view()
		msg := fmt.Sprintf("v2: %+v", p)
		assert(t, p.initSize.w, 20, "p.initSize.w", msg)
		assert(t, p.initSize.h, 20, "p.initSize.h", msg)
		assert(t, p.size.w, 20, "p.size.w", msg)
		assert(t, p.size.h, 20, "p.size.h", msg)
	}

	{
		p := v1.view()
		msg := fmt.Sprintf("v1: %+v", p)
		assert(t, p.initSize.w, -1, "p.initSize.w", msg)
		assert(t, p.initSize.h, -1, "p.initSize.h", msg)
		assert(t, p.size.w, 100, "p.size.w", msg)
		assert(t, p.size.h, 100, "p.size.h", msg)
	}
}
