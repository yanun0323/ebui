package ebui

import (
	"embed"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

/*
	FONT
*/

var (
	//go:embed resource/NotoSansTC.ttf
	_defaultTTF          embed.FS
	_defaultFontResource = defaultFont()
	_fontTagWeight       = parseTag("wght") /* 100-900 */
	_fontTagItalic       = parseTag("ital") /* 0-1 */
)

func defaultFont() *text.GoTextFaceSource {
	f, err := _defaultTTF.Open("resource/NotoSansTC.ttf")
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
