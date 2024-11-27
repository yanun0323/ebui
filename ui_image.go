package ebui

import (
	"embed"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
	"github.com/yanun0323/pkg/sys"
)

var (
	supportedImageFormat = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".avif", ".heic", ".heif"}
)

type uiImage struct {
	*uiView

	img *ebiten.Image
}

func Image(path string, embedFile ...embed.FS) SomeView {
	ui := &uiImage{}

	ui.uiView = newView(typesImage, ui)
	ui.tryLoadingImage(path, embedFile...)
	return ui
}

func (ui *uiImage) tryLoadingImage(path string, embedFile ...embed.FS) {
	filenames := make([]string, 0, len(supportedImageFormat))
	hasExtension := false
	lowercasePath := strings.ToLower(path)
	for _, format := range supportedImageFormat {
		hasExtension = hasExtension || strings.HasSuffix(lowercasePath, format)
		filenames = append(filenames, path+format)
	}

	if hasExtension {
		filenames = []string{path}
	}

	paths := make([]string, 0, len(filenames)*2)
	if len(embedFile) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			logs.Errorf("get working directory, err: %v", err)
			return
		}

		exeDir, err := os.Executable()
		if err != nil {
			logs.Errorf("get executable directory, err: %v", err)
			return
		}

		for _, filename := range filenames {
			paths = append(paths, filepath.Join(wd, filename), filepath.Join(exeDir, filename))
		}
	} else {
		paths = append(paths, filenames...)
	}

	var (
		r   io.ReadCloser
		err error
	)

	for _, p := range paths {
		if len(embedFile) != 0 {
			if dir, err := embedFile[0].ReadDir("."); err == nil {
				for _, f := range dir {
					i, _ := f.Info()
					logs.Infof("name: %s, size: %d", f.Name(), i.Size())
				}
			}

			r, err = embedFile[0].Open(p)
			if err != nil {
				logs.Warnf("open image with embed file (%s), err: %v", p, err)
				continue
			}
			defer r.Close()
		} else {
			r, err = os.Open(p)
			if err != nil {
				logs.Warnf("open image (%s), err: %v", p, err)
				continue
			}
			defer r.Close()
		}

		img, _, err := image.Decode(r)
		if err != nil {
			logs.Errorf("decode image (%s), err: %v", p, err)
			continue
		}

		ui.img = ebiten.NewImageFromImage(img)
		break
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
