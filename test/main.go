package main

import (
	"embed"
	"time"

	. "github.com/yanun0323/ebui"
)

//go:embed *
var resource embed.FS

func main() {
	clr := New(ColorGray)
	size := New(Size{W: 150, H: 150})
	offset := New(Point{X: 0, Y: 0})

	go func() {
		for {
			time.Sleep(time.Second)
			if clr.Get() == ColorGray {
				clr.Set(ColorWhite)
			} else {
				clr.Set(ColorGray)
			}
		}
	}()

	go func() {
		right := true
		for {
			time.Sleep(time.Millisecond)
			o := offset.Get()
			if right {
				o.X += 1
				offset.Set(o)
			} else {

				o.X -= 1
				offset.Set(o)
			}

			if o.X > 500 {
				right = false
			}

			if o.X < 0 {
				right = true
			}
		}
	}()

	contentView := ZStack(
		Circle().
			Frame(size).
			Offset(offset).
			ForegroundColor(clr).
			Opacity(New(0.5)),
		Text(New("Hello")).
			Offset(offset).
			Frame(New(Size{W: 100, H: 100})).
			ForegroundColor(New(ColorWhite)),
	)

	SetWindowSize(640, 480)
	Run("Windows Title", contentView,
		WithDebug(),
		WithMemMonitor(),
	)
}
