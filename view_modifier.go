package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/sys"
)

type uiViewModifier struct {
	uiViewLayout

	borderWidth int
	borderColor color.Color
}

func newViewModifier() uiViewModifier {
	return uiViewModifier{
		uiViewLayout: _zeroUIViewLayout,
	}
}

func (vm *uiViewModifier) drawModifiers(screen *ebiten.Image, ui *uiView) {
	pos := vm.getDrawSize(ui.cachedSize)
	if pos.w <= 0 || pos.h <= 0 {
		return
	}

	img := ebiten.NewImage(pos.w, pos.h)
	img.Fill(sys.If(vm.background == nil, color.Color(color.Transparent), vm.background))

	if vm.borderColor != nil && vm.borderWidth >= 1 {
		for x := 0; x < pos.w; x++ {
			for y := 0; y < pos.h; y++ {
				if x == 0 || x == pos.w-1 || y == 0 || y == pos.h-1 {
					img.Set(x, y, vm.borderColor)
				}
			}
		}
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(ui.start.x), float64(ui.start.y))
	opt.GeoM.Translate(float64(vm.offset.x), float64(vm.offset.y))
	opt.GeoM.Translate(float64(vm.margin.left), float64(vm.margin.top))
	ui.handlePreference(opt)
	screen.DrawImage(img, opt)
}
