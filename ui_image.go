package ebui

import (
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type imageImpl struct {
	*ctx

	image *Binding[*ebiten.Image]
}

func Image[T string | *ebiten.Image](img *Binding[T]) SomeView {
	switch content := any(img).(type) {
	case *Binding[*ebiten.Image]:
		v := &imageImpl{
			image: content,
		}
		v.ctx = newViewContext(v)
		return v
	case *Binding[string]:
		var img *ebiten.Image
		content.addListener(func() {
			img = getImage(content.Get())
		})
		constraint := BindFunc(func() *ebiten.Image {
			if img == nil {
				img = getImage(content.Get())
				if img == nil {
					img = ebiten.NewImage(1, 1)
				}
			}

			return img
		}, func(i *ebiten.Image) {})

		return Image(constraint)
	}

	return nil
}

func (v *imageImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions {
	op := v.ctx.draw(screen, hook...)
	img := v.image.Get()
	if img == nil {
		return op
	}

	frame := v.ctx.systemSetFrame()
	if !frame.drawable() {
		return op
	}

	frameSize := frame.Size()
	imgSize := sz(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	scale := v.getScale(frameSize, imgSize)

	opt := &ebiten.DrawImageOptions{}
	opt.ColorScale.ScaleWithColorScale(op.ColorScale)
	opt.GeoM.Scale(scale.X, scale.Y)
	opt.GeoM.Concat(op.GeoM)

	screen.DrawImage(img, opt)

	return opt
}

func (v *imageImpl) getScale(frameSize, imgSize CGSize) CGPoint {
	scale := pt(1, 1)
	if !v.ctx.scaleToFit.Get() {
		return scale
	}

	keepAspectRatio := v.ctx.keepAspectRatio.Get()
	if !keepAspectRatio {
		return pt(frameSize.Width/imgSize.Width, frameSize.Height/imgSize.Height)
	}

	scaleX := frameSize.Width / imgSize.Width
	scaleY := frameSize.Height / imgSize.Height
	s := scaleX
	if scaleY < scaleX {
		s = scaleY
	}

	return pt(s, s)
}

func getImage(filename string) *ebiten.Image {
	path := filename
	if dir, ok := resourceDir.Load().(string); ok && len(dir) != 0 {
		println("resourceDir:", dir)
		if !filepath.IsAbs(dir) {
			println("not abs:", dir)
			dir, _ = filepath.Abs(dir)
		}
		path = filepath.Join(dir, filename)
	}

	println("image:", path)

	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}

	return ebiten.NewImageFromImage(img)
}
