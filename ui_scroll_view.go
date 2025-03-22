package ebui

import (
	"image/color"
	"strconv"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/animation"
	"github.com/yanun0323/ebui/input"
	"github.com/yanun0323/ebui/internal/helper"
	layout "github.com/yanun0323/ebui/layout"
)

const (
	_scrollIndicatorLength  = 15.0
	_scrollIndicatorPadding = 6.0
)

// ScrollView is a view that can scroll its content.
//
// Use ScrollViewDirection to specify the direction of the scroll view.
//
// Default direction is vertical.
func ScrollView(content View) SomeView {
	sv := &scrollViewImpl{
		contentOffset:    Bind(CGPoint{}),
		content:          content.Body(),
		indicatorOpacity: Bind(0.0),
		isLastHover:      Bind(false),
		indicateCache:    helper.NewHashCache[[2]*ebiten.Image](),
	}
	sv.viewCtx = newViewContext(sv)

	return sv
}

type scrollViewImpl struct {
	*viewCtx

	content          SomeView
	contentOffset    *Binding[CGPoint]
	indicatorOpacity *Binding[float64]
	isLastHover      *Binding[bool]

	indicateCache *helper.HashCache[[2]*ebiten.Image] // [base, main]
}

func (s *scrollViewImpl) maxScrollOffset() CGPoint {
	sFrameSize := s.systemSetFrame().Size()
	cBoundsSize := s.content.systemSetBounds().Size()

	return NewPoint(
		max(cBoundsSize.Width-sFrameSize.Width, 0),
		max(cBoundsSize.Height-sFrameSize.Height, 0),
	)
}

func (s *scrollViewImpl) setScrollOffset(delta input.Vector, d layout.Direction, maxOffset CGPoint) {
	offset := s.contentOffset.Value()
	switch d {
	case layout.DirectionVertical:
		offset.Y += delta.Y
		switch {
		case offset.Y < 0:
			offset.Y = 0
		case offset.Y > maxOffset.Y:
			offset.Y = maxOffset.Y
		}
	case layout.DirectionHorizontal:
		offset.X += delta.X
		switch {
		case offset.X < 0:
			offset.X = 0
		case offset.X > maxOffset.X:
			offset.X = maxOffset.X
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
		childFrameSize := childBoundsSize.Shrink(sData.Padding.Add(sData.Border))
		sBounds, sAlignFunc := sLayout(start, childFrameSize)
		sFrameSize := sBounds.Size().Shrink(sData.Padding.Add(sData.Border))
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
	if !bounds.drawable() {
		return
	}

	offset := s.contentOffset.Value()
	base := ebiten.NewImage(int(bounds.Dx()), int(bounds.Dy()))
	s.content.draw(base, func(opt *ebiten.DrawImageOptions) {
		opt.GeoM.Translate(-offset.X, -offset.Y)
		opt.GeoM.Translate(-bounds.Start.X, -bounds.Start.Y)
	})
	s.drawScrollIndicator(base)

	screen.DrawImage(base, s.drawOption(bounds, hook...))
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
			optOffset = NewPoint(_scrollIndicatorPadding/2, offset.Y*ratio)
			radius = mainSize.Width / 2
		case layout.DirectionHorizontal:
			baseSize = NewSize(sBoundsWidth, _scrollIndicatorLength)
			opt.GeoM.Translate(0, sBoundsHeight-_scrollIndicatorLength)
			opt.ColorScale.ScaleAlpha(float32(s.indicatorOpacity.Value()))

			ratio := sBoundsWidth / cBoundsSize.Width
			mainSize = NewSize(ratio*sBoundsWidth, _scrollIndicatorLength-_scrollIndicatorPadding)
			optOffset = NewPoint(offset.X*ratio, _scrollIndicatorPadding/2)
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

func (s *scrollViewImpl) shiftCursor(cursor input.Vector) input.Vector {
	scrollOffset := s.contentOffset.Value()
	return cursor.Add(scrollOffset.X, scrollOffset.Y)
}

func (s *scrollViewImpl) onHoverEvent(cursor input.Vector) {
	isHover := s.isHover(cursor)
	if isHover != s.isLastHover.Get() {
		opacity := s.indicatorOpacity.Get()
		if isHover {
			if opacity == 0.0 {
				s.indicatorOpacity.Set(1.0, nil)
			}
		} else {
			if opacity != 0.0 {
				s.indicatorOpacity.Set(0.0, _scrollIndicatorDisappearAnimation)
			}
		}
		s.isLastHover.Set(isHover, nil)
	}

	s.viewCtx.onHoverEvent(cursor)
	s.content.onHoverEvent(s.shiftCursor(cursor))
}

var (
	_scrollIndicatorAppearAnimation    = animation.EaseInOut(300 * time.Millisecond)
	_scrollIndicatorDisappearAnimation = _scrollIndicatorAppearAnimation.Delay(300 * time.Millisecond)
)

func (s *scrollViewImpl) onScrollEvent(cursor input.Vector, event input.ScrollEvent) bool {
	opacity := s.indicatorOpacity.Get()
	if !s.isHover(cursor) {
		if opacity != 0.0 {
			s.indicatorOpacity.Set(0.0, _scrollIndicatorDisappearAnimation)
		}

		return false
	}

	defer func() {
		s.content.onScrollEvent(s.shiftCursor(cursor), event)
	}()
	defer s.viewCtx.onScrollEvent(cursor, event)

	d := s.scrollViewDirection.Get()
	maxOffset := s.maxScrollOffset()
	if !event.Delta.IsZero() {
		s.setScrollOffset(event.Delta, d, maxOffset)
	}

	var (
		threshold = 1.0
		scrolling = false
	)

	switch d {
	case layout.DirectionVertical:
		if abs(event.Delta.Y) <= threshold {
			event.Delta.Y = 0
		}
		scrolling = event.Delta.Y != 0
	case layout.DirectionHorizontal:
		if abs(event.Delta.X) <= threshold {
			event.Delta.X = 0
		}
		scrolling = event.Delta.X != 0
	}

	lastScrolling := opacity != 0
	if lastScrolling != scrolling {
		if scrolling {
			s.indicatorOpacity.Set(1.0, nil)
		} else {
			s.indicatorOpacity.Set(0.0, _scrollIndicatorDisappearAnimation)
		}
	}

	return true
}

func (s *scrollViewImpl) onMouseEvent(event input.MouseEvent) {
	s.viewCtx.onMouseEvent(event)

	event.Position = s.shiftCursor(event.Position)
	s.content.onMouseEvent(event)
}

func (s *scrollViewImpl) onKeyEvent(event input.KeyEvent) {
	s.viewCtx.onKeyEvent(event)
	s.content.onKeyEvent(event)
}

func (s *scrollViewImpl) onTypeEvent(event input.TypeEvent) {
	s.viewCtx.onTypeEvent(event)
	s.content.onTypeEvent(event)
}

func (s *scrollViewImpl) onGestureEvent(event input.GestureEvent) {
	s.viewCtx.onGestureEvent(event)

	event.Position = s.shiftCursor(event.Position)
	s.content.onGestureEvent(event)
}

func (s *scrollViewImpl) onTouchEvent(event input.TouchEvent) {
	s.viewCtx.onTouchEvent(event)

	event.Position = s.shiftCursor(event.Position)
	s.content.onTouchEvent(event)
}
