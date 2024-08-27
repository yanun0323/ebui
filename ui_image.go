package ebui

import (
	"bytes"
	"image"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

type uiImage struct {
	*uiView

	img *ebiten.Image
}

func Image(path string, embedFile ...fs.FS) SomeView {
	ui := &uiImage{}

	ui.uiView = newUIView(typesImage, ui)
	ui.tryLoadingImage(path, embedFile...)
	return ui
}

func (ui *uiImage) tryLoadingImage(path string, embedFile ...fs.FS) {
	wd, err := os.Getwd()
	if err != nil {
		logs.Errorf("Failed to get working directory: %v", err)
		return
	}

	filename := filepath.Join(wd, path)

	var r io.Reader
	if len(embedFile) != 0 && embedFile[0] != nil {
		f, err := fs.ReadFile(embedFile[0], filename)
		if err != nil {
			logs.Errorf("Failed to open image with embed file (%s): %v", filename, err)
			return
		}
		r = bytes.NewReader(f) 
	} else {
		r, err = os.Open(filename)
		if err != nil {
			logs.Errorf("Failed to open image (%s): %v", filename, err)
			return
		}
	}

	img, _, err := image.Decode(r)
	if err != nil {
		logs.Errorf("Failed to decode image (%s): %v", filename, err)
		return
	}

	ui.img = ebiten.NewImageFromImage(img)
}

func (ui *uiImage) draw(screen *ebiten.Image) {
	if ui.img == nil {
		return
	}

	cache := ui.Copy()
	cache.Draw(screen, func(img *ebiten.Image) {
		cache.IterateViewModifiers(func(vm viewModifier) {
			_ = vm(img, cache)
		})

		wRatio := float64(cache.Width()) / float64(ui.img.Bounds().Dx())
		hRatio := float64(cache.Height()) / float64(ui.img.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		ratio := max(wRatio, hRatio)
		op.GeoM.Scale(ratio, ratio)
		op.GeoM.Translate(float64(cache.padding.left), float64(cache.padding.top))
		screen.DrawImage(ui.img, op)
	})
}
