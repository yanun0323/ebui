package ebui

import (
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type imageImpl struct {
	*viewCtx

	image *Binding[*ebiten.Image]
}

func Image[T string | *ebiten.Image](img *Binding[T]) SomeView {
	switch content := any(img).(type) {
	case *Binding[*ebiten.Image]:
		v := &imageImpl{
			image: content,
		}
		v.viewCtx = newViewContext(v)
		return v
	case *Binding[string]:
		path := bindCombineOneWay(resourceDir, content, func(dir, filename string) string {
			return getImageFilename(dir, filename)
		})

		return Image(BindOneWay(path, getImage))
	}

	return nil
}

func (v *imageImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	v.viewCtx.draw(screen, hook...)
	img := v.image.Value()
	if img == nil {
		return
	}

	frame := v.viewCtx.systemSetFrame()
	if !frame.drawable() {
		return
	}

	op := v.viewCtx.drawOption(frame, hook...)

	frameSize := frame.Size()
	imgSize := NewSize(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	scale := v.getScale(frameSize, imgSize)

	opt := &ebiten.DrawImageOptions{}
	opt.ColorScale.ScaleWithColorScale(op.ColorScale)
	opt.GeoM.Scale(scale.X, scale.Y)
	opt.GeoM.Concat(op.GeoM)

	screen.DrawImage(img, opt)
}

func (v *imageImpl) getScale(frameSize, imgSize CGSize) CGPoint {
	scale := NewPoint(1, 1)
	if !v.viewCtx.scaleToFit.Value() {
		return scale
	}

	keepAspectRatio := v.viewCtx.keepAspectRatio.Value()
	if !keepAspectRatio {
		return NewPoint(frameSize.Width/imgSize.Width, frameSize.Height/imgSize.Height)
	}

	scaleX := frameSize.Width / imgSize.Width
	scaleY := frameSize.Height / imgSize.Height
	s := min(scaleY, scaleX)

	return NewPoint(s, s)
}

func getImageFilename(dir, filename string) string {
	path := filename
	if len(dir) != 0 {
		if !filepath.IsAbs(dir) {
			dir, _ = filepath.Abs(dir)
		}
		path = filepath.Join(dir, filename)
	}

	return path
}

func getImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		logf("error: %+v", err)
		return nil
	}
	defer f.Close()

	logf("get image: %s", path)

	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}

	return ebiten.NewImageFromImage(img)
}
