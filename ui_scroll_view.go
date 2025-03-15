package ebui

import (
	"image/color"
	"strconv"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/internal/helper"
	layout "github.com/yanun0323/ebui/layout"
)

const (
	_scrollIndicatorLength  = 15.0
	_scrollIndicatorPadding = 6.0
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

	indicateCache *helper.HashCache[[2]*ebiten.Image] // [base, main]
}

func (s *scrollViewImpl) clampScrollDelta(sFrameSize CGSize, cBoundsSize CGSize, delta CGPoint) CGPoint {
	if sFrameSize.Height >= cBoundsSize.Height {
		delta.Y = 0
	}

	if sFrameSize.Width >= cBoundsSize.Width {
		delta.X = 0
	}

	return delta
}

func (s *scrollViewImpl) setScrollOffset(delta CGPoint, d layout.Direction) {
	sFrameSize := s.systemSetFrame().Size()
	cBoundsSize := s.content.systemSetBounds().Size()
	delta = s.clampScrollDelta(sFrameSize, cBoundsSize, delta)
	if delta.IsZero() {
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
	_, cLayout := s.content.preload(parent, types...)
	sData, sLayout := s.viewCtx.preload(parent, types...)

	return sData, func(start CGPoint, childBoundsSize CGSize) (CGRect, alignFunc) {
		childFrameSize := childBoundsSize.Shrink(sData.Padding).Shrink(sData.Border)
		sBounds, sAlignFunc := sLayout(start, childFrameSize)
		sFrameSize := sBounds.Size().Shrink(sData.Padding).Shrink(sData.Border)
		cBounds, cAlignFunc := cLayout(start, sFrameSize)

		sBoundsSize := sBounds.Size()
		cBoundsSize := cBounds.Size()

		vec := NewPoint(s.alignment.Get().Vector())
		offset := NewPoint(
			max(sBoundsSize.Width-cBoundsSize.Width, 0)*vec.X,
			max(sBoundsSize.Height-cBoundsSize.Height, 0)*vec.Y,
		)
		cAlignFunc(offset)

		return sBounds, func(offset CGPoint) {
			sAlignFunc(offset)
			cAlignFunc(offset)
		}
	}
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

	if !bounds.drawable() {
		return
	}

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

		baseColor = NewColor(16, 16, 16, 64)
		mainColor = NewColor(64, 64, 64, 64)

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
			baseSize = NewSize(_scrollIndicatorLength, sBoundsHeight)
			opt.GeoM.Translate(sBoundsWidth-_scrollIndicatorLength, 0)
			opt.ColorScale.ScaleAlpha(float32(s.indicatorOpacity.Value()))

			ratio := sBoundsHeight / cBoundsSize.Height
			mainSize = NewSize(_scrollIndicatorLength-_scrollIndicatorPadding, ratio*sBoundsHeight)
			optOffset = NewPoint(_scrollIndicatorPadding/2, -offset.Y*ratio)
			radius = mainSize.Width / 2
		case layout.DirectionHorizontal:
			baseSize = NewSize(sBoundsWidth, _scrollIndicatorLength)
			opt.GeoM.Translate(0, sBoundsHeight-_scrollIndicatorLength)
			opt.ColorScale.ScaleAlpha(float32(s.indicatorOpacity.Value()))

			ratio := sBoundsWidth / cBoundsSize.Width
			mainSize = NewSize(ratio*sBoundsWidth, _scrollIndicatorLength-_scrollIndicatorPadding)
			optOffset = NewPoint(-offset.X*ratio, _scrollIndicatorPadding/2)
			radius = mainSize.Height / 2
		}
	}

	img := s.indicateCache.Load()
	if s.indicateCache.IsNextCacheOutdated() || img[0] == nil || img[1] == nil {
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
	anim := animation.EaseInOut(300 * time.Millisecond)
	if !onHover(s.systemSetBounds()) {
		if opacity := s.indicatorOpacity.Get(); opacity != 0.0 {
			s.indicatorOpacity.Set(0.0, anim)
		}

		return
	}

	d := s.scrollViewDirection.Get()
	if !event.Delta.IsZero() {
		s.setScrollOffset(event.Delta, d)
	}

	scrolling := false
	switch d {
	case layout.DirectionVertical:
		scrolling = event.Delta.Y != 0
	case layout.DirectionHorizontal:
		scrolling = event.Delta.X != 0
	}

	if scrolling {
		opacity := s.indicatorOpacity.Get()
		if opacity != 1.0 {
			s.indicatorOpacity.Set(1.0, anim)
		}
	} else {
		opacity := s.indicatorOpacity.Get()
		if opacity != 0.0 {
			s.indicatorOpacity.Set(0.0, anim)
		}
	}
}

func (s *scrollViewImpl) HandleTouchEvent(event touchEvent) {

}

func (s *scrollViewImpl) HandleKeyEvent(event keyEvent) {
}

func (s *scrollViewImpl) HandleInputEvent(event inputEvent) {
}
