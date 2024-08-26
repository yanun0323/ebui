package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

type size struct {
	w, h int
}

func root(contentView SomeView, s ...size) *rootView {
	v := &rootView{}
	v.view = newView(typeRoot, v, contentView)
	if len(s) == 0 {
		w, h := ebiten.WindowSize()
		s = append(s, size{w, h})
	}
	v.setWindowSize(s[0].w, s[0].h)

	return v
}

type rootView struct {
	*view
}

func (r *rootView) setWindowSize(w, h int) {
	r.view.initW, r.view.initH = w, h
	r.view.w, r.view.h = w, h
}

func (r *rootView) calculateStage() {
	r.preCacheChildrenSize()
	r.calculationParameters()
}

func (r *rootView) draw(screen *ebiten.Image) {
	r.view.Draw(screen, func(screen *ebiten.Image) {
		r.view.IterateViewModifiers(func(vm viewModifier) {
			vv := vm(screen, r.view)
			if vv != nil {
				vv.draw(screen)
			}
		})
	})
}

func (r *rootView) preCacheChildrenSize() {
	var formula func(SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int)
	formula = func(current SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int) {
		currentView := current.params()
		initW, initH := current.initBounds()
		maxWidth, maxHeight = initW, initH
		for _, sv := range currentView.subviews {
			maxW, maxH, sumW, sumH := formula(sv)

			maxWidth = max(maxWidth, maxW)
			maxHeight = max(maxHeight, maxH)
			sumWidth += replaceIf(sumW, -1, 0)
			sumHeight += replaceIf(sumH, -1, 0)
		}

		sumWidth = max(sumWidth, initW)
		sumHeight = max(sumHeight, initH)

		switch currentView.types {
		case typeVStack:
			currentView.w = maxWidth
			currentView.h = replaceIf(sumHeight, 0, -1)
		case typeHStack:
			currentView.w = replaceIf(sumWidth, 0, -1)
			currentView.h = maxHeight
		case typeZStack:
			currentView.w = maxWidth
			currentView.h = maxHeight
		default:
			currentView.w = -1
			currentView.h = -1
		}

		return currentView.w, currentView.h, sumWidth, sumHeight
	}

	formula(r.view)
	r.view.w, r.view.h = r.initBounds()
}

func (r *rootView) calculationParameters() {
	var formula func(v, anchor *view)
	formula = func(v, anchor *view) {
		// update params from anchor
		v.x = anchor.x
		v.y = anchor.y
		v.xx = anchor.xx
		v.yy = anchor.yy
		v.viewOpacity *= anchor.viewOpacity
		v.fColor = replaceIfZero(v.fColor, anchor.fColor)
		v.fontSizes = replaceIfZero(v.fontSizes, anchor.fontSizes)
		v.isPressing = replaceIfZero(v.isPressing, anchor.isPressing)

		// apply view modifiers
		nextAnchor := v.Copy()
		for _, vm := range nextAnchor.viewModifiers {
			_ = vm(nil, nextAnchor)
		}

		// update cache after modifying
		nextAnchor.x += nextAnchor.paddingLeft
		nextAnchor.y += nextAnchor.paddingTop
		nextAnchor.xx = nextAnchor.paddingLeft // TODO: check if this is correct
		nextAnchor.yy = nextAnchor.paddingTop

		// fmt.Printf("next anchor: %+v\n", nextAnchor)

		// calculate flexible size
		wFlexCount, hFlexCount, wDelta, hDelta, recalculatedSubViews := r.countFlexibleChildren(nextAnchor)

		setSubviews := func() (again bool) {
			if wFlexCount <= -1 || hFlexCount <= -1 {
				l := logs.New(logs.LevelDebug).WithField("current", v.types)
				l.Fatalf("%s: flex count is negative: wFlexCount: %d, hFlexCount: %d", nextAnchor.types, wFlexCount, hFlexCount)
			}

			width := (nextAnchor.Width() - wDelta) / replaceIfZero(wFlexCount, 1)
			height := (nextAnchor.Height() - hDelta) / replaceIfZero(hFlexCount, 1)

			for _, sv := range nextAnchor.subviews {
				svp := sv.params()
				ll := logs.New(logs.LevelDebug).WithField("parent", v.types).WithField("current", svp.types)
				ll.Debugf("x,y(%d, %d), xx,yy(%d, %d), w,h(%d, %d)", svp.x, svp.y, svp.xx, svp.yy, svp.w, svp.h)

				if !recalculatedSubViews[svp] {
					recalculatedSubViews[svp] = true
					// calculate that does width/height need to recalculate
					switch v.types {
					case typeVStack, typeHStack, typeZStack:
						if svp.w > width {
							ll.Warnf("width out of range! svp.w: %d > width: %d", svp.w, width)
							wFlexCount--
							wDelta += svp.w
							again = true
						}

						if svp.h > height {
							ll.Warnf("height out of range! svp.h: %d > height: %d", svp.h, height)
							hFlexCount--
							hDelta += svp.h
							again = true
						}
					}

					if again {
						return true
					}
				}

				// set size to subviews
				switch v.types {
				case typeVStack:
					svp.w = replaceIf(svp.initW, -1, max(svp.w, nextAnchor.Width()))
					svp.h = replaceIf(svp.initH, -1, max(svp.h, height))
				case typeHStack:
					svp.w = replaceIf(svp.initW, -1, max(svp.w, width))
					svp.h = replaceIf(svp.initH, -1, max(svp.h, nextAnchor.Height()))
				default:
					svp.w = replaceIf(svp.initW, -1, max(svp.w, width))
					svp.h = replaceIf(svp.initH, -1, max(svp.h, height))
				}

				// deep calculate subviews
				formula(svp, nextAnchor)

				// update cache after calculating subviews
				switch v.types {
				case typeVStack:
					nextAnchor.y, nextAnchor.yy = nextAnchor.y+svp.h, nextAnchor.yy+svp.h
				case typeHStack:
					nextAnchor.x, nextAnchor.xx = nextAnchor.x+svp.w, nextAnchor.xx+svp.w
				}
			}

			return false
		}

		again := true
		for again {
			again = setSubviews()
		}
	}

	formula(r.view, r.view)
}

func (r *rootView) countFlexibleChildren(v *view) (widthCount, heightCount, widthDelta, heightDelta int, recalculatedSubViews map[*view]bool) {
	table := make(map[*view]bool, len(v.subviews))
	wCount, hCount, wDelta, hDelta := 0, 0, 0, 0
	for _, sv := range v.subviews {
		svp := sv.params()
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
