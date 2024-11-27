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
	"github.com/yanun0323/pkg/sys"
)

type uiImage struct {
	*uiView

	img *ebiten.Image
}

func Image(path string, embedFile ...fs.FS) SomeView {
	ui := &uiImage{}

	ui.uiView = newView(typesImage, ui)
	ui.tryLoadingImage(path, embedFile...)
	return ui
}

func (ui *uiImage) tryLoadingImage(path string, embedFile ...fs.FS) {
	wd, err := os.Getwd()
	if err != nil {
		logs.Errorf("Failed to get working directory: %v", err)
		return
	}

	exeDir, err := os.Executable()
	if err != nil {
		logs.Errorf("Failed to get executable directory: %v", err)
		return
	}

	filenames := []string{filepath.Join(wd, path), filepath.Join(exeDir, path)}

	var r io.Reader
	for _, filename := range filenames {
		if len(embedFile) != 0 && embedFile[0] != nil {
			f, err := fs.ReadFile(embedFile[0], filename)
			if err != nil {
				logs.Errorf("Failed to open image with embed file (%s): %v", filename, err)
				continue
			}
			r = bytes.NewReader(f)
		} else {
			r, err = os.Open(filename)
			if err != nil {
				logs.Errorf("Failed to open image (%s): %v", filename, err)
				continue
			}
		}

		img, _, err := image.Decode(r)
		if err != nil {
			logs.Errorf("Failed to decode image (%s): %v", filename, err)
			continue
		}

		ui.img = ebiten.NewImageFromImage(img)
	}
}

func (v *uiImage) getSize() size {
	if v.isCached {
		return v.cachedSize
	}

	s := v.uiView.getSize()
	if s.IsZero() && !v.resizable {
		r := v.img.Bounds()
		s = size{r.Dx(), r.Dy()}
	}

	v.isCached = true
	v.cachedSize = s
	return v.cachedSize
}

func (ui *uiImage) draw(screen *ebiten.Image) {
	if ui.img == nil {
		return
	}

	drawSize := ui.getDrawSize(ui.cachedSize)
	wRatio := float64(drawSize.w) / float64(ui.img.Bounds().Dx())
	hRatio := float64(drawSize.h) / float64(ui.img.Bounds().Dy())
	dx := sys.If(wRatio <= 0, 0.0, float64(ui.start.x))
	dy := sys.If(hRatio <= 0, 0.0, float64(ui.start.y))
	wRatio = sys.If(wRatio < _minimumFloat64, _minimumFloat64, wRatio)
	hRatio = sys.If(hRatio < _minimumFloat64, _minimumFloat64, hRatio)

	op := &ebiten.DrawImageOptions{}
	if ui.aspectRatio != 0 {
		if wRatio <= hRatio {
			hRatio = wRatio * ui.aspectRatio
			dy += (float64(drawSize.h) - float64(ui.img.Bounds().Dy())*hRatio) / 2
		} else {
			wRatio = hRatio / ui.aspectRatio
			dx += (float64(drawSize.w) - float64(ui.img.Bounds().Dx())*wRatio) / 2
		}
	}

	op.GeoM.Scale(wRatio, hRatio)
	op.GeoM.Translate(dx, dy)
	screen.DrawImage(ui.img, op)
}
