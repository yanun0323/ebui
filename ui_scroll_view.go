package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	layout "github.com/yanun0323/ebui/layout"
)

func ScrollView(direction layout.Direction, content View) SomeView {
	sv := &scrollViewImpl{
		content:   content.Body(),
		direction: direction,
	}
	sv.viewCtx = newViewContext(sv)
	globalEventManager.RegisterHandler(sv)
	return sv
}

type scrollViewImpl struct {
	*viewCtx

	content   SomeView
	offset    CGPoint
	direction layout.Direction
}

func (s *scrollViewImpl) isScrollable(d layout.Direction) bool {
	sFrameSize := s.systemSetFrame().Size()
	cBoundsSize := s.content.systemSetBounds().Size()
	logf("sFrameSize: %v, cBoundsSize: %v", sFrameSize, cBoundsSize)
	switch d {
	case layout.DirectionVertical:
		if sFrameSize.Height >= cBoundsSize.Height {
			return false
		}
	case layout.DirectionHorizontal:
		if sFrameSize.Width >= cBoundsSize.Width {
			return false
		}
	}

	return true
}

func (s *scrollViewImpl) setScrollOffset(delta CGPoint) {
	if !s.isScrollable(s.direction) {
		return
	}

	switch s.direction {
	case layout.DirectionVertical:
		s.offset.Y += delta.Y
	case layout.DirectionHorizontal:
		s.offset.X += delta.X
	}
}

func (s *scrollViewImpl) preload(parent *viewCtxEnv, types ...formulaType) (preloadData, layoutFunc) {
	stackFormula := &formulaStack{
		types:                           formulaZStack,
		stackCtx:                        s.viewCtx,
		children:                        []SomeView{s.content},
		ignorePreloadingChildSummedSize: true,
	}

	return stackFormula.preload(parent, types...)
}

func (s *scrollViewImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	hooks := make([]func(*ebiten.DrawImageOptions), 0, len(hook)+1)
	hooks = append(hooks, hook...)
	hooks = append(hooks, func(opt *ebiten.DrawImageOptions) {
		opt.GeoM.Translate(s.offset.X, s.offset.Y)
	})

	s.viewCtx.draw(screen, hooks...)
	s.content.draw(screen, hooks...)
}

func (s *scrollViewImpl) HandleWheelEvent(event wheelEvent) {
	s.setScrollOffset(event.Delta)
}

func (s *scrollViewImpl) HandleTouchEvent(event touchEvent) {

}

func (s *scrollViewImpl) HandleKeyEvent(event keyEvent) {
}

func (s *scrollViewImpl) HandleInputEvent(event inputEvent) {
}
