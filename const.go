package ebui

import (
	"embed"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	_roundedScale        float64 = 3.0
	_roundedScaleInverse float64 = 1.0 / _roundedScale
)

/*
	COLOR
*/

var (
	white       = NewColor(255)
	black       = NewColor(0)
	transparent = CGColor{}
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
