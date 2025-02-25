package font

import (
    "log"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "golang.org/x/image/font/gofont/goregular"
)

var (
    defaultFace font.Face
)

func init() {
    tt, err := opentype.Parse(goregular.TTF)
    if err != nil {
        log.Fatal(err)
    }

    defaultFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
        Size:    24,
        DPI:     72,
        Hinting: font.HintingFull,
    })
    if err != nil {
        log.Fatal(err)
    }
}

func GetFontFace(size Size, weight Weight) font.Face {
    // 這裡先返回預設字體，之後可以根據 size 和 weight 來返回不同的字體
    return defaultFace
} 