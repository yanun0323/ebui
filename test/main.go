package main

import (
	"embed"
	"image/color"
	"math/rand"
	"time"

	. "github.com/yanun0323/ebui"
)

//go:embed *
var resource embed.FS

func main() {
	clr := New(ColorGray)
	size := New(Size{W: 30, H: 30})
	offset := New(Point{X: 0, Y: 100})

	go func() {
		for {
			time.Sleep(time.Second)
			if clr.Get() == ColorGray {
				clr.Set(color.Gray{50})
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

	gen := func(count int) []SomeView {
		views := make([]SomeView, 0, count*5)

		ofs := New(Point{X: rand.Intn(500), Y: rand.Intn(500)})
		views = append(views, Rectangle().
			Frame(size).
			Offset(offset).
			Offset(ofs).
			ForegroundColor(clr),
		)

		for i := 0; i < count; i++ {
			ofs := New(Point{X: rand.Intn(500), Y: rand.Intn(500)})
			views = append(views, Rectangle().
				Frame(size).
				Offset(offset).
				Offset(ofs).
				ForegroundColor(clr),
			)
		}

		for i := 0; i < count; i++ {
			ofs := New(Point{X: rand.Intn(500), Y: rand.Intn(500)})
			views = append(views, Circle().
				Frame(size).
				Offset(offset).
				Offset(ofs).
				ForegroundColor(clr),
			)
		}

		for i := 0; i < count; i++ {
			ofs := New(Point{X: rand.Intn(500), Y: rand.Intn(500)})
			views = append(views, Text(New("Hello")).
				Frame(size).
				Offset(offset).
				Offset(ofs).
				ForegroundColor(clr),
			)
		}

		return views
	}

	contentView := ZStack(gen(100)...)

	SetWindowSize(640, 480)
	Run("Windows Title", contentView,
		WithDebug(),
		WithMemMonitor(),
	)
}
