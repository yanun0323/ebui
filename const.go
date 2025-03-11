package ebui

import (
	"embed"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	_roundedScale        = 3.0
	_roundedScaleInverse = 1.0 / _roundedScale
)

/*
	FONT
*/

var (
	//go:embed resource/NotoSansTC.ttf
	defaultTTF          embed.FS
	defaultFontResource = defaultFont()
	fontTagWeight       = parseTag("wght") /* 100-900 */
	fontTagItalic       = parseTag("ital") /* 0-1 */
)

var (
	AccentColor = NewColor(0, 0, 255, 255)
	transparent = AnyColor(color.Transparent)
)

func defaultFont() *text.GoTextFaceSource {
	f, err := defaultTTF.Open("resource/NotoSansTC.ttf")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open font: %w", err))
	}
	defer f.Close()

	s, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create font: %w", err))
	}

	return s
}

func parseTag(tag string) text.Tag {
	result, err := text.ParseTag(tag)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse tag: %s, err: %w", tag, err))
	}

	return result
}
