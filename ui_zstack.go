package ebui

// import (
// 	"github.com/hajimehoshi/ebiten/v2"
// )

// /* Check Interface Implementation */
// var _ SomeView = (*zstackView)(nil)

// func ZStack(views ...View) *zstackView {
// 	v := &zstackView{}
// 	v.uiViewBack = newUIView(typesZStack, v, views...)
// 	return v
// }

// type zstackView struct {
// 	*uiViewBack
// }

// func (v *zstackView) draw(screen *ebiten.Image) {
// 	cache := v.uiViewBack.Copy()
// 	cache.Draw(screen, func(screen *ebiten.Image) {
// 		cache.ApplyViewModifiers(screen)
// 	})
// }
