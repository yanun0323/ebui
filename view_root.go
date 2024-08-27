package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

func root(contentView SomeView, s ...frame) *rootView {
	v := &rootView{}
	v.uiView = newUIView(typesRoot, v, contentView)
	if len(s) == 0 {
		w, h := ebiten.WindowSize()
		s = append(s, frame{w, h})
	}
	v.setWindowSize(s[0].w, s[0].h)

	return v
}

type rootView struct {
	*uiView
}

func (r *rootView) setWindowSize(w, h int) {
	r.uiView.initSize = frame{w, h}
	r.uiView.size = frame{w, h}
}

func (r *rootView) calculateStage() {
	r.preCacheChildrenSize()
	r.calculationParameters()
}

func (r *rootView) draw(screen *ebiten.Image) {
	r.uiView.Draw(screen, func(screen *ebiten.Image) {
		r.uiView.IterateViewModifiers(func(vm viewModifier) {
			vv := vm(screen, r.uiView)
			if vv != nil {
				vv.draw(screen)
			}
		})
	})
}

func (r *rootView) preCacheChildrenSize() {
	var formula func(SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int)
	formula = func(current SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int) {
		currentView := current.view()
		initW, initH := current.initBounds()
		maxWidth, maxHeight = initW, initH
		for _, sv := range currentView.subviews {
			maxW, maxH, sumW, sumH := formula(sv)

			maxWidth = max(maxWidth, maxW)
			maxHeight = max(maxHeight, maxH)
			sumWidth += rpEq(sumW, -1, 0)
			sumHeight += rpEq(sumH, -1, 0)
		}

		sumWidth = max(sumWidth, initW)
		sumHeight = max(sumHeight, initH)

		switch currentView.types {
		case typesVStack:
			currentView.size.w = maxWidth
			currentView.size.h = rpEq(sumHeight, 0, -1)
		case typesHStack:
			currentView.size.w = rpEq(sumWidth, 0, -1)
			currentView.size.h = maxHeight
		case typesZStack:
			currentView.size.w = maxWidth
			currentView.size.h = maxHeight
		default:
			currentView.size.w = initW
			currentView.size.h = initH
		}

		return currentView.size.w, currentView.size.h, sumWidth, sumHeight
	}

	formula(r.uiView)
	r.uiView.size.w, r.uiView.size.h = r.initBounds()
}

func (r *rootView) calculationParameters() {
	var formula func(v, anchor *uiView)
	formula = func(v, anchor *uiView) {
		// update params from anchor
		v.pos.x = anchor.pos.x
		v.pos.y = anchor.pos.y
		v.xx = anchor.xx
		v.yy = anchor.yy
		v.viewOpacity *= anchor.viewOpacity
		v.fColor = rpZero(v.fColor, anchor.fColor)
		v.fontSizes = rpZero(v.fontSizes, anchor.fontSizes)
		v.isPressing = rpZero(v.isPressing, anchor.isPressing)

		// apply view modifiers
		nextAnchor := v.Copy()
		for _, vm := range nextAnchor.viewModifiers {
			_ = vm(nil, nextAnchor)
		}

		// update cache after modifying
		nextAnchor.pos.x += nextAnchor.padding.left
		nextAnchor.pos.y += nextAnchor.padding.top
		nextAnchor.xx = nextAnchor.padding.left // TODO: check if this is correct
		nextAnchor.yy = nextAnchor.padding.top

		// fmt.Printf("next anchor: %+v\n", nextAnchor)

		// calculate flexible size
		wFlexCount, hFlexCount, wDelta, hDelta, recalculatedSubViews := r.countFlexibleChildren(nextAnchor)

		setSubviews := func() (again bool) {
			if wFlexCount <= -1 || hFlexCount <= -1 {
				l := logs.Default().WithField("current", v.types)
				l.Fatalf("%s: flex count is negative: wFlexCount: %d, hFlexCount: %d", nextAnchor.types, wFlexCount, hFlexCount)
			}

			width := (nextAnchor.Width() - wDelta) / rpZero(wFlexCount, 1)
			height := (nextAnchor.Height() - hDelta) / rpZero(hFlexCount, 1)

			for _, sv := range nextAnchor.subviews {
				svp := sv.view()
				ll := logs.Default().WithField("parent", v.types).WithField("current", svp.types)
				ll.Debugf("x,y(%d, %d), xx,yy(%d, %d), w,h(%d, %d)", svp.pos.x, svp.pos.y, svp.xx, svp.yy, svp.size.w, svp.size.h)

				if !recalculatedSubViews[svp] {
					recalculatedSubViews[svp] = true
					// calculate that does width/height need to recalculate
					switch v.types {
					case typesVStack, typesHStack, typesZStack:
						if svp.size.w > width {
							ll.Debugf("width out of range! svp.size.w: %d > width: %d", svp.size.w, width)
							wFlexCount--
							wDelta += svp.size.w
							again = true
						}

						if svp.size.h > height {
							ll.Debugf("height out of range! svp.size.h: %d > height: %d", svp.size.h, height)
							hFlexCount--
							hDelta += svp.size.h
							again = true
						}
					}

					if again {
						return true
					}
				}

				// set size to subviews
				switch v.types {
				case typesVStack:
					svp.size.w = rpEq(svp.initSize.w, -1, max(svp.size.w, nextAnchor.Width()))
					svp.size.h = rpEq(svp.initSize.h, -1, max(svp.size.h, height))
				case typesHStack:
					svp.size.w = rpEq(svp.initSize.w, -1, max(svp.size.w, width))
					svp.size.h = rpEq(svp.initSize.h, -1, max(svp.size.h, nextAnchor.Height()))
				default:
					svp.size.w = rpEq(svp.initSize.w, -1, max(svp.size.w, width))
					svp.size.h = rpEq(svp.initSize.h, -1, max(svp.size.h, height))
				}

				// deep calculate subviews
				formula(svp, nextAnchor)

				// update cache after calculating subviews
				switch v.types {
				case typesVStack:
					nextAnchor.pos.y, nextAnchor.yy = nextAnchor.pos.y+svp.size.h, nextAnchor.yy+svp.size.h
				case typesHStack:
					nextAnchor.pos.x, nextAnchor.xx = nextAnchor.pos.x+svp.size.w, nextAnchor.xx+svp.size.w
				}
			}

			return false
		}

		again := true
		for again {
			again = setSubviews()
		}
	}

	formula(r.uiView, r.uiView)
}

func (r *rootView) countFlexibleChildren(v *uiView) (widthCount, heightCount, widthDelta, heightDelta int, recalculatedSubViews map[*uiView]bool) {
	table := make(map[*uiView]bool, len(v.subviews))
	wCount, hCount, wDelta, hDelta := 0, 0, 0, 0
	for _, sv := range v.subviews {
		svp := sv.view()
		vvw, vvh := sv.initBounds()

		if vvw == -1 {
			wCount++
		} else {
			table[svp] = true
			wDelta += vvw
		}

		if vvh == -1 {
			hCount++
		} else {
			table[svp] = true
			hDelta += vvh
		}
	}

	return wCount, hCount, wDelta, hDelta, table
}
