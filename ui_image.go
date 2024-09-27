package ebui

// import (
// 	"bytes"
// 	"image"
// 	"io"
// 	"io/fs"
// 	"os"
// 	"path/filepath"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/yanun0323/pkg/logs"
// )

// type uiImage struct {
// 	*uiViewBack

// 	img *ebiten.Image
// }

// func Image(path string, embedFile ...fs.FS) SomeView {
// 	ui := &uiImage{}

// 	ui.uiViewBack = newUIView(typesImage, ui)
// 	ui.tryLoadingImage(path, embedFile...)
// 	return ui
// }

// func (ui *uiImage) tryLoadingImage(path string, embedFile ...fs.FS) {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		logs.Errorf("Failed to get working directory: %v", err)
// 		return
// 	}

// 	filename := filepath.Join(wd, path)

// 	var r io.Reader
// 	if len(embedFile) != 0 && embedFile[0] != nil {
// 		f, err := fs.ReadFile(embedFile[0], filename)
// 		if err != nil {
// 			logs.Errorf("Failed to open image with embed file (%s): %v", filename, err)
// 			return
// 		}
// 		r = bytes.NewReader(f)
// 	} else {
// 		r, err = os.Open(filename)
// 		if err != nil {
// 			logs.Errorf("Failed to open image (%s): %v", filename, err)
// 			return
// 		}
// 	}

// 	img, _, err := image.Decode(r)
// 	if err != nil {
// 		logs.Errorf("Failed to decode image (%s): %v", filename, err)
// 		return
// 	}

// 	ui.img = ebiten.NewImageFromImage(img)
// }

// func (ui *uiImage) draw(screen *ebiten.Image) {
// 	if ui.img == nil {
// 		return
// 	}

// 	cache := ui.Copy()
// 	cache.Draw(screen, func(screen *ebiten.Image) {
// 		cache.ApplyViewModifiers(screen)

// 		wRatio := float64(cache.Width()) / float64(ui.img.Bounds().Dx())
// 		hRatio := float64(cache.Height()) / float64(ui.img.Bounds().Dy())
// 		op := &ebiten.DrawImageOptions{}
// 		ratio := min(wRatio, hRatio)
// 		dx, dy := 0, 0

// 		if wRatio != ratio {
// 			dx = abs(int(float64(ui.img.Bounds().Dx())*ratio - float64(cache.Width())))
// 		}

// 		if hRatio != ratio {
// 			dy = abs(int(float64(ui.img.Bounds().Dy())*ratio - float64(cache.Height())))
// 		}

// 		op.GeoM.Scale(ratio, ratio)
// 		op.GeoM.Translate(float64(cache.padding.left+dx/2), float64(cache.padding.top+dy/2))
// 		screen.DrawImage(ui.img, op)
// 	})
// }
