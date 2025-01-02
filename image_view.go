package view

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageView struct {
	image     *ebiten.Image
	frame     image.Rectangle
	scale     float64
	rotation  float64
}

func Image(img *ebiten.Image) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			return &ImageView{
				image: img,
				scale: 1.0,
			}
		},
	}
}

func (iv *ImageView) Layout(bounds image.Rectangle) image.Rectangle {
	size := iv.image.Bounds().Size()
	iv.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X + int(float64(size.X)*iv.scale),
		bounds.Min.Y + int(float64(size.Y)*iv.scale),
	)
	return iv.frame
}

func (iv *ImageView) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(iv.image.Bounds().Dx())/2, -float64(iv.image.Bounds().Dy())/2)
	op.GeoM.Rotate(iv.rotation)
	op.GeoM.Scale(iv.scale, iv.scale)
	op.GeoM.Translate(
		float64(iv.frame.Min.X) + float64(iv.frame.Dx())/2,
		float64(iv.frame.Min.Y) + float64(iv.frame.Dy())/2,
	)
	screen.DrawImage(iv.image, op)
} 