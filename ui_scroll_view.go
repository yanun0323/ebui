package ebui

import (
	"image/color"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/internal/helper"
	layout "github.com/yanun0323/ebui/layout"
)

func ScrollView(content View) SomeView {
	sv := &scrollViewImpl{
		contentOffset:    Bind(CGPoint{}),
		content:          content.Body(),
		indicatorOpacity: Bind(0.0),
		indicateCache:    helper.NewHashCache[[2]*ebiten.Image](),
	}
	sv.viewCtx = newViewContext(sv)
	globalEventManager.RegisterHandler(sv)
	return sv
}

type scrollViewImpl struct {
	*viewCtx

	content          SomeView
	contentOffset    *Binding[CGPoint]
	indicatorOpacity *Binding[float64]
	hovered          atomic.Bool

	indicateCache *helper.HashCache[[2]*ebiten.Image] // [base, main]
}

func (s *scrollViewImpl) isScrollable(sFrameSize CGSize, cBoundsSize CGSize, d layout.Direction) bool {
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
	if !s.hovered.Load() {
		return
	}

	d := s.scrollViewDirection.Get()
	sFrameSize := s.systemSetFrame().Size()
	cBoundsSize := s.content.systemSetBounds().Size()
	if !s.isScrollable(sFrameSize, cBoundsSize, d) {
		return
	}

	offset := s.contentOffset.Value()
	switch d {
	case layout.DirectionVertical:
		offset.Y += delta.Y
		floor := -(cBoundsSize.Height - sFrameSize.Height)
		switch {
		case offset.Y > 0:
			offset.Y = 0
		case offset.Y < floor:
			offset.Y = floor
		}
	case layout.DirectionHorizontal:
		offset.X += delta.X
		floor := -(cBoundsSize.Width - sFrameSize.Width)
		switch {
		case offset.X > 0:
			offset.X = 0
		case offset.X < floor:
			offset.X = floor
		}
	}

	s.contentOffset.Set(offset, nil)
}

func (s *scrollViewImpl) Hash() string {
	h := xxhash.New()
	h.Write(s.content.bytes())
	h.Write(s.bytes())

	return strconv.FormatUint(h.Sum64(), 16)
}

func (s *scrollViewImpl) preload(parent *viewCtxEnv, types ...stackType) (preloadData, layoutFunc) {
	stackFormula := &stackPreloader{
		types:                 newStackTypeFromDirection(s.scrollViewDirection.Get()),
		stackCtx:              s.viewCtx,
		children:              []SomeView{s.content},
		preloadStackOnlyFrame: true,
	}

	s.indicateCache.SetNextHash(s.Hash())

	return stackFormula.preload(parent, types...)
}

func (s *scrollViewImpl) draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	s.viewCtx.draw(screen, hook...)
	bounds := s.systemSetBounds()

	offset := s.contentOffset.Value()
	contentHook := make([]func(*ebiten.DrawImageOptions), 0, len(hook)+1)
	contentHook = append(contentHook, hook...)
	contentHook = append(contentHook, func(opt *ebiten.DrawImageOptions) {
		opt.GeoM.Translate(offset.X, offset.Y)
		opt.GeoM.Translate(-bounds.Start.X, -bounds.Start.Y)
	})

	base := ebiten.NewImage(int(bounds.Dx()), int(bounds.Dy()))

	s.content.draw(base, contentHook...)
	s.drawScrollIndicator(base, hook...)

	screen.DrawImage(base, s.drawOption(bounds))
}

func (s *scrollViewImpl) drawScrollIndicator(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions)) {
	opacity := s.indicatorOpacity.Value()
	if opacity == 0.0 {
		return
	}

	var (
		radius   float64
		baseSize CGSize
		mainSize CGSize

		length           = 15.0
		baseColor        = NewColor(16, 16, 16, 64)
		mainColor        = NewColor(64, 64, 64, 64)
		indicatorPadding = 6.0

		d                           = s.scrollViewDirection.Get()
		offset                      = s.contentOffset.Value()
		sBoundsSize                 = s.systemSetBounds().Size()
		sBoundsHeight, sBoundsWidth = sBoundsSize.Height, sBoundsSize.Width
		cBoundsSize                 = s.content.systemSetBounds().Size()

		opt       = s.drawOption(CGRect{}, hook...)
		optOffset CGPoint
	)

	{
		switch d {
		case layout.DirectionVertical:
			baseSize = NewSize(length, sBoundsHeight)
			opt.GeoM.Translate(sBoundsWidth-length, 0)
			opt.ColorScale.ScaleAlpha(float32(s.indicatorOpacity.Value()))

			ratio := sBoundsHeight / cBoundsSize.Height
			mainSize = NewSize(length-indicatorPadding, ratio*sBoundsHeight)
			optOffset = NewPoint(indicatorPadding/2, -offset.Y*ratio)
			radius = mainSize.Width / 2
		case layout.DirectionHorizontal:
			baseSize = NewSize(sBoundsWidth, length)
			opt.GeoM.Translate(0, sBoundsHeight-length)
			opt.ColorScale.ScaleAlpha(float32(s.indicatorOpacity.Value()))

			ratio := sBoundsWidth / cBoundsSize.Width
			mainSize = NewSize(ratio*sBoundsWidth, length-indicatorPadding)
			optOffset = NewPoint(-offset.X*ratio, indicatorPadding/2)
			radius = mainSize.Height / 2
		}
	}

	var img [2]*ebiten.Image
	if s.indicateCache.IsNextHashCached() {
		img = s.indicateCache.Load()
	} else {
		img[0] = ebiten.NewImage(int(baseSize.Width), int(baseSize.Height))
		img[0].Fill(baseColor)

		w := int(mainSize.Width * _roundedScale)
		h := int(mainSize.Height * _roundedScale)
		r := float64(int(radius * _roundedScale))
		img[1] = ebiten.NewImage(w, h)
		img[1].Fill(mainColor)
		cornerHandler := newCornerHandler(w, h, r)
		cornerHandler.Execute(func(isOutside, isBorder bool, x, y int) {
			if isOutside {
				img[1].Set(x, y, color.Transparent)
				return
			}
		})

		s.indicateCache.Update(img)
	}

	screen.DrawImage(img[0], opt)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(_roundedScaleInverse, _roundedScaleInverse)
	op.Filter = ebiten.FilterLinear
	op.GeoM.Concat(opt.GeoM)
	op.GeoM.Translate(optOffset.X, optOffset.Y)
	op.ColorScale.ScaleWithColorScale(opt.ColorScale)
	screen.DrawImage(img[1], op)
}

func (s *scrollViewImpl) HandleWheelEvent(event wheelEvent) {
	s.setScrollOffset(event.Delta)
	hovered := onHover(s.systemSetBounds())
	if s.hovered.Swap(hovered) != hovered {
		opacity := 0.0
		if hovered {
			opacity = 1.0
		}
		s.indicatorOpacity.Set(opacity, animation.EaseInOut(300*time.Millisecond))
	}
}

func (s *scrollViewImpl) HandleTouchEvent(event touchEvent) {

}

func (s *scrollViewImpl) HandleKeyEvent(event keyEvent) {
}

func (s *scrollViewImpl) HandleInputEvent(event inputEvent) {
}
