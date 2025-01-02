package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ViewCache struct {
	lastFrame   image.Rectangle
	cachedImage *ebiten.Image
	isDirty     bool
	lastHash    uint64
}

func NewViewCache() *ViewCache {
	return &ViewCache{
		isDirty: true,
	}
}

func (vc *ViewCache) Draw(screen *ebiten.Image, view View, currentHash uint64) {
	if currentHash != vc.lastHash {
		vc.isDirty = true
		vc.lastHash = currentHash
	}

	if !vc.lastFrame.Eq(view.Layout(screen.Bounds())) {
		vc.isDirty = true
		vc.lastFrame = view.Layout(screen.Bounds())
	}

	if vc.isDirty || vc.cachedImage == nil {
		if vc.cachedImage == nil || !vc.lastFrame.Size().Eq(vc.cachedImage.Bounds().Size()) {
			vc.cachedImage = ebiten.NewImage(vc.lastFrame.Dx(), vc.lastFrame.Dy())
		}
		vc.cachedImage.Clear()
		view.Draw(vc.cachedImage)
		vc.isDirty = false
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(vc.lastFrame.Min.X), float64(vc.lastFrame.Min.Y))
	screen.DrawImage(vc.cachedImage, op)
}
