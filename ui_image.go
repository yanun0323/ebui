package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageView struct {
	*ctx

	image    *ebiten.Image
	frame    CGRect
	scale    float64
	rotation float64
}

func Image(img *ebiten.Image) SomeView {
	iv := &ImageView{
		image: img,
		scale: 1.0,
	}
	iv.ctx = newViewContext(tagImage, iv)
	return iv
}

func (iv *ImageView) layout(bounds CGRect) (CGRect, func(CGRect)) {
	size := iv.image.Bounds().Size()
	iv.frame = rect(
		bounds.Start.X,
		bounds.Start.Y,
		bounds.Start.X+float64(size.X)*iv.scale,
		bounds.Start.Y+float64(size.Y)*iv.scale,
	)
	return iv.frame, nil
}

func (iv *ImageView) draw(screen *ebiten.Image, bounds ...CGRect) {
	iv.ctx.draw(screen)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(iv.image.Bounds().Dx())/2, -float64(iv.image.Bounds().Dy())/2)
	op.GeoM.Rotate(iv.rotation)
	op.GeoM.Scale(iv.scale, iv.scale)
	op.GeoM.Translate(
		float64(iv.frame.Start.X)+float64(iv.frame.Dx())/2,
		float64(iv.frame.Start.Y)+float64(iv.frame.Dy())/2,
	)
	screen.DrawImage(iv.image, op)
}
