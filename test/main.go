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
	size := New(Size{W: 640, H: 480})

	go func() {
		for {
			time.Sleep(time.Second)
			s := size.Get()
			s.W -= 1
			size.Set(s)
			if clr.Get() == ColorGray {
				clr.Set(ColorWhite)
			} else {
				clr.Set(ColorGray)
			}
		}
	}()

	contentView := Text(New("Hello")).Frame(size).ForegroundColor(clr).Opacity(New(0.5))

	SetWindowSize(640, 480)
	Run("Windows Title", contentView,
		WithDebug(),
		WithMemMonitor(),
	)
}
