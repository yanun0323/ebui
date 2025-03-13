package ebui

import "github.com/yanun0323/ebui/layout"

func animateValue[T bindable](startValue, endValue T, delta float64) T {
	switch any(startValue).(type) {
	case layout.Align:
		return endValue
	case int:
		sv := any(startValue).(int)
		ev := any(endValue).(int)
		res := float64(sv) + float64(ev-sv)*delta
		return any(int(res)).(T)
	case int8:
		sv := any(startValue).(int8)
		ev := any(endValue).(int8)
		res := float64(sv) + float64(ev-sv)*delta
		return any(int8(res)).(T)
	case int16:
		sv := any(startValue).(int16)
		ev := any(endValue).(int16)
		res := float64(sv) + float64(ev-sv)*delta
		return any(int16(res)).(T)
	case int32:
		sv := any(startValue).(int32)
		ev := any(endValue).(int32)
		res := float64(sv) + float64(ev-sv)*delta
		return any(int32(res)).(T)
	case int64:
		sv := any(startValue).(int64)
		ev := any(endValue).(int64)
		res := float64(sv) + float64(ev-sv)*delta
		return any(int64(res)).(T)
	case uint:
		sv := any(startValue).(uint)
		ev := any(endValue).(uint)
		res := float64(sv) + float64(ev-sv)*delta
		return any(uint(res)).(T)
	case uint8:
		sv := any(startValue).(uint8)
		ev := any(endValue).(uint8)
		res := float64(sv) + float64(ev-sv)*delta
		return any(uint8(res)).(T)
	case uint16:
		sv := any(startValue).(uint16)
		ev := any(endValue).(uint16)
		res := float64(sv) + float64(ev-sv)*delta
		return any(uint16(res)).(T)
	case uint32:
		sv := any(startValue).(uint32)
		ev := any(endValue).(uint32)
		res := float64(sv) + float64(ev-sv)*delta
		return any(uint32(res)).(T)
	case uint64:
		sv := any(startValue).(uint64)
		ev := any(endValue).(uint64)
		res := float64(sv) + float64(ev-sv)*delta
		return any(uint64(res)).(T)
	case uintptr:
		sv := any(startValue).(uintptr)
		ev := any(endValue).(uintptr)
		res := float64(sv) + float64(ev-sv)*delta
		return any(uintptr(res)).(T)
	case float32:
		sv := any(startValue).(float32)
		ev := any(endValue).(float32)
		res := float64(sv) + float64(ev-sv)*delta
		return any(float32(res)).(T)
	case float64:
		sv := any(startValue).(float64)
		ev := any(endValue).(float64)
		res := sv + (ev-sv)*delta
		return any(res).(T)
	case CGPoint:
		sv := any(startValue).(CGPoint)
		ev := any(endValue).(CGPoint)
		x := sv.X + (ev.X-sv.X)*delta
		y := sv.Y + (ev.Y-sv.Y)*delta
		return any(CGPoint{X: x, Y: y}).(T)
	case CGSize:
		sv := any(startValue).(CGSize)
		ev := any(endValue).(CGSize)
		w := sv.Width + (ev.Width-sv.Width)*delta
		h := sv.Height + (ev.Height-sv.Height)*delta
		return any(CGSize{Width: w, Height: h}).(T)
	case CGRect:
		sv := any(startValue).(CGRect)
		ev := any(endValue).(CGRect)
		sx := sv.Start.X + (ev.Start.X-sv.Start.X)*delta
		sy := sv.Start.Y + (ev.Start.Y-sv.Start.Y)*delta
		ex := sv.End.X + (ev.End.X-sv.End.X)*delta
		ey := sv.End.Y + (ev.End.Y-sv.End.Y)*delta
		return any(CGRect{Start: CGPoint{X: sx, Y: sy}, End: CGPoint{X: ex, Y: ey}}).(T)
	case CGInset:
		sv := any(startValue).(CGInset)
		ev := any(endValue).(CGInset)
		t := sv.Top + (ev.Top-sv.Top)*delta
		r := sv.Right + (ev.Right-sv.Right)*delta
		b := sv.Bottom + (ev.Bottom-sv.Bottom)*delta
		l := sv.Left + (ev.Left-sv.Left)*delta
		return any(CGInset{Top: t, Right: r, Bottom: b, Left: l}).(T)
	case CGColor:
		sv := any(startValue).(CGColor)
		ev := any(endValue).(CGColor)
		r := uint8(float64(sv.R) + (float64(ev.R)-float64(sv.R))*delta)
		g := uint8(float64(sv.G) + (float64(ev.G)-float64(sv.G))*delta)
		b := uint8(float64(sv.B) + (float64(ev.B)-float64(sv.B))*delta)
		a := uint8(float64(sv.A) + (float64(ev.A)-float64(sv.A))*delta)
		return any(CGColor{R: r, G: g, B: b, A: a}).(T)
	}

	return endValue
}
