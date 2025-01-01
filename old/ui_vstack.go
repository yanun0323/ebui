package ebui

// type vstack struct {
// 	*view

// 	views []SomeView
// }

// func VStack(views ...SomeView) SomeView {
// 	v := &vstack{
// 		views: views,
// 	}
// 	v.view = newView(idVStack, v)
// 	return v
// }

// func (v *vstack) bounds() (min, current, max Size) {
// 	return v.view.bounds()
// }

// func (v *vstack) update(container Size) {
// 	v.view.updateRenderCache()

// 	size := v.param.frameSize
// 	flexWidthCount := 0
// 	flexHeightCount := 0

// 	for _, child := range v.views {
// 		_, childSize, _ := child.bounds()
// 		size.W = max(size.W, childSize.W)
// 		size.H += childSize.H
// 		flexWidthCount += sys.If(childSize.W <= 0, 0, 1)
// 		flexHeightCount += sys.If(childSize.H <= 0, 0, 1)
// 	}

// 	size.W = sys.If(flexWidthCount != 0, 0, size.W)
// 	size.H = sys.If(flexHeightCount != 0, 0, size.H)

// }

// func (v *vstack) draw(screen *ebiten.Image) {
// 	v.updateRenderCache()

// 	for _, child := range v.views {
// 		child.draw(screen)
// 	}
// }
