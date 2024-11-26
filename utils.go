package ebui

// type ordered interface {
// 	~int | ~int8 | ~int16 | ~int32 | ~int64 |
// 		~uintptr |
// 		~float32 | ~float64
// }

// func abs[T ordered](num T) T {
// 	if num < 0 {
// 		return -num
// 	}

// 	return num
// }

// // rpZero returns result if v is zero value, otherwise returns v.
// func rpZero[T comparable](v, result T) T {
// 	var zero T
// 	return rpEq(v, zero, result)
// }

// // rpEq returns result if v is eq, otherwise returns v.
// func rpEq[T comparable](v, eq, result T) T {
// 	if v == eq {
// 		return result
// 	}

// 	return v
// }

// first returns the zero if slice is empty, otherwise returns the first element of slice.
func first[T any](slice []T, zero T) T {
	if len(slice) != 0 {
		return slice[0]
	}
	return zero
}

// makeImageRounded makes image rounded.
// func makeImageRounded(img *ebiten.Image, current *uiViewBack, round int) {
// 	if img == nil || round <= 0 {
// 		return
// 	}

// 	sW, sH := img.Bounds().Dx(), img.Bounds().Dy()
// 	if sW <= 0 || sH <= 0 {
// 		return
// 	}

// 	set := func(x, y int) {
// 		// x += current.padding.left
// 		// y += current.padding.top
// 		if x >= 0 && y >= 0 && x < sW && y < sH {
// 			img.Set(x, y, color.Transparent)
// 		}
// 	}

// 	w, h := current.size.w, current.size.h

// 	for x := 0; x < round; x++ {
// 		// left-top
// 		ltX, ltY := round, round
// 		for y := 0; y < round; y++ {
// 			dx := x - ltX
// 			dy := y - ltY
// 			if dx*dx+dy*dy > round*round {
// 				set(x, y)
// 			}
// 		}

// 		// left-bottom
// 		lbX, lbY := round, h-round
// 		for y := h - round; y < h; y++ {
// 			dx := x - lbX
// 			dy := y - lbY
// 			if dx*dx+dy*dy > round*round {
// 				set(x, y)
// 			}
// 		}
// 	}

// 	for x := w - round; x < w; x++ {
// 		// right-top
// 		rtX, rtY := w-round, round
// 		for y := 0; y < round; y++ {
// 			dx := x - rtX
// 			dy := y - rtY
// 			if dx*dx+dy*dy > round*round {
// 				set(x, y)
// 			}
// 		}

// 		// right-bottom
// 		rbX, rbY := w-round, h-round
// 		for y := h - round; y < h; y++ {
// 			dx := x - rbX
// 			dy := y - rbY
// 			if dx*dx+dy*dy > round*round {
// 				set(x, y)
// 			}
// 		}
// 	}
// }
