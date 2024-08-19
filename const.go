package ebui

import (
	"bytes"
	"image/color"
	"log"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	_mplusFaceSource        = defaultFont()
	_defaultForegroundColor = color.Black
	_globalTicker           = atomic.Int64{}
)

func defaultFont() *text.GoTextFaceSource {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func currentTicker() int64 {
	return _globalTicker.Load()
}

func tickTock() {
	_ = _globalTicker.Add(1)
}
