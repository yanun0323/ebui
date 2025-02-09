package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageView struct {
	*viewContext

	image    *ebiten.Image
	frame    image.Rectangle
	scale    float64
	rotation float64
}

func Image(img *ebiten.Image) SomeView {
	iv := &ImageView{
		image: img,
		scale: 1.0,
	}
	iv.viewContext = NewViewContext(iv)
	return iv
}

func (iv *ImageView) layout(bounds image.Rectangle) image.Rectangle {
	size := iv.image.Bounds().Size()
	iv.frame = image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X+int(float64(size.X)*iv.scale),
		bounds.Min.Y+int(float64(size.Y)*iv.scale),
	)
	return iv.frame
}

func (iv *ImageView) draw(screen *ebiten.Image) {
	iv.viewContext.draw(screen)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(iv.image.Bounds().Dx())/2, -float64(iv.image.Bounds().Dy())/2)
	op.GeoM.Rotate(iv.rotation)
	op.GeoM.Scale(iv.scale, iv.scale)
	op.GeoM.Translate(
		float64(iv.frame.Min.X)+float64(iv.frame.Dx())/2,
		float64(iv.frame.Min.Y)+float64(iv.frame.Dy())/2,
	)
	screen.DrawImage(iv.image, op)
}

// 新增的方法，用於設置縮放比例
func (iv *ImageView) WithScale(scale float64) *ImageView {
	iv.scale = scale
	return iv
}

// 新增的方法，用於設置旋轉角度
func (iv *ImageView) WithRotation(rotation float64) *ImageView {
	iv.rotation = rotation
	return iv
}
