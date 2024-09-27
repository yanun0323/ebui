package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/pkg/logs"
)

func root(contentView SomeView) *rootView {
	r := &rootView{}
	r.uiView = newView(typesRoot, r, contentView)

	w, h := ebiten.WindowSize()
	r.uiView.flexible = size{w, h}
	r.preCacheChildrenSizeToFlexible()
	r.calculateParameters()

	if len(r.uiView.contents) != 0 {
		invokeSomeView(r.uiView.contents[0], func(svp *uiView) {
			_rootStart.Store(point{
				x: (w - svp.flexible.w) / 2,
				y: (h - svp.flexible.h) / 2,
			})
			logs.Warnf("!!! w: %d, h %d, svp.flexible.w: %d, svp.flexible.h: %d", w, h, svp.flexible.w, svp.flexible.h)
		})
	}

	r.uiView.ActionUpdate()

	return r
}

type rootView struct {
	*uiView
}

func (r *rootView) Draw(screen *ebiten.Image) {
	r.uiView.Draw(screen)
}

func (r *rootView) preCacheChildrenSizeToFlexible() {
	windowWidth, windowHeight := ebiten.WindowSize()
	var deepGetSize func(SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int)
	deepGetSize = func(sv SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int) {
		maxWidth, maxHeight, sumWidth, sumHeight = -1, -1, 0, 0
		invokeSomeView(sv, func(ui *uiView) {
			initSize := size{windowWidth, windowHeight}
			if ui.types != typesRoot {
				initSize = ui.GetFrameWithMargin()
			}
			maxWidth, maxHeight = 0, 0

			spacerW, spacerH := false, false
			for _, content := range ui.contents {
				maxW, maxH, sumW, sumH := deepGetSize(content)
				invokeSomeView(content, func(svp *uiView) {
					spacerW = spacerW || svp.frame.w == -1
					spacerH = spacerH || svp.frame.h == -1

					if spacerW {
						maxWidth = -1
					} else {
						maxWidth = max(maxWidth, maxW)
					}

					if spacerH {
						maxHeight = -1
					} else {
						maxHeight = max(maxHeight, maxH)
					}

					sumWidth += rpEq(sumW, -1, 0)
					sumHeight += rpEq(sumH, -1, 0)
				})
			}

			sumWidth = max(sumWidth, initSize.w)
			sumHeight = max(sumHeight, initSize.h)

			switch ui.types {
			case typesVStack:
				ui.flexible.w = maxWidth
				ui.flexible.h = rpEq(sumHeight, 0, -1)
			case typesHStack:
				ui.flexible.w = rpEq(sumWidth, 0, -1)
				ui.flexible.h = maxHeight
			case typesZStack:
				ui.flexible.w = rpEq(sumWidth, 0, maxWidth)
				ui.flexible.h = rpEq(sumHeight, 0, maxHeight)
			default:
				ui.flexible.w = initSize.w
				ui.flexible.h = initSize.h
			}

			maxWidth = ui.flexible.w
			maxHeight = ui.flexible.h
		})

		return maxWidth, maxHeight, rpEq(sumWidth, 0, -1), rpEq(sumHeight, 0, -1)
	}

	deepGetSize(r)
}

func (r *rootView) calculateParameters() {
	var deepCalculate func(ui, anchor *uiView)
	deepCalculate = func(ui, anchor *uiView) {
		// update params from anchor
		ui.uiViewEnvironment.set(anchor.uiViewEnvironment)
		if ui != anchor {
			ui.opacity *= anchor.opacity
			ui.start = anchor.start
		}

		// apply view modifiers
		nextAnchor := ui.Copy()
		nextAnchor.start = anchor.start

		// calculate flexible size
		wFlexCount, hFlexCount, wSizedDelta, hSizedDelta, wSizedMax, hSizedMax, sizedSubViews := r.countFlexibleChildren(nextAnchor)

		setSubviews := func() (again bool) {
			l := logs.Default().WithField("current", ui.types)
			if wFlexCount < 0 || hFlexCount < 0 {
				l.Fatalf("%s: flex count is negative: wFlexCount: %d, hFlexCount: %d", nextAnchor.types, wFlexCount, hFlexCount)
			}

			width, height := 0, 0
			switch ui.types {
			case typesVStack:
				width = (nextAnchor.flexible.w - wSizedMax) / rpZero(wFlexCount, 1)
				height = (nextAnchor.flexible.h - hSizedDelta) / rpZero(hFlexCount, 1)
			case typesHStack:
				width = (nextAnchor.flexible.w - wSizedDelta) / rpZero(wFlexCount, 1)
				height = (nextAnchor.flexible.h - hSizedMax) / rpZero(hFlexCount, 1)
			default:
				width = (nextAnchor.flexible.w - wSizedDelta) / rpZero(wFlexCount, 1)
				height = (nextAnchor.flexible.h - hSizedDelta) / rpZero(hFlexCount, 1)
			}
			l.Infof("%s, width: %d, height: %d", ui.types, width, height)
			l.Infof("%s, flexible.w: %d, flexible.h: %d, delta.w: %d, delta.h: %d, max.w: %d, max.h: %d",
				ui.types,
				nextAnchor.flexible.w, nextAnchor.flexible.h,
				wSizedDelta, hSizedDelta,
				wSizedMax, hSizedMax,
			)

			for _, sv := range nextAnchor.contents {
				again := false
				invokeSomeView(sv, func(svp *uiView) {
					ll := logs.Default().WithField("parent", ui.types).WithField("current", svp.types)
					ll.Debugf("start: %v, offset: %v, frame: %v", svp.start, svp.offset, svp.frame)

					if !sizedSubViews[svp] {
						sizedSubViews[svp] = true
						// calculate that does width/height need to recalculate
						switch ui.types {
						case typesVStack, typesHStack, typesZStack:
							if svp.flexible.w > width {
								ll.Warnf("width out of range! svp.flexible.w: %d > width: %d", svp.flexible.w, width)
								wFlexCount--
								wSizedDelta += svp.flexible.w
								again = true
							}

							if svp.flexible.h > height {
								ll.Warnf("height out of range! svp.flexible.h: %d > height: %d", svp.flexible.h, height)
								hFlexCount--
								hSizedDelta += svp.flexible.h
								again = true
							}
						}

						if again {
							ll.Info("again")
							return
						}
					}

					// set position to subviews
					svp.start = nextAnchor.start

					// set size to subviews
					switch ui.types {
					case typesVStack:
						svp.flexible.w = rpEq(svp.flexible.w, -1, nextAnchor.flexible.w)
						svp.flexible.h = rpEq(svp.flexible.h, -1, height)
					case typesHStack:
						svp.flexible.w = rpEq(svp.flexible.w, -1, width)
						svp.flexible.h = rpEq(svp.flexible.h, -1, nextAnchor.flexible.h)
					default:
						svp.flexible.w = rpEq(svp.flexible.w, -1, width)
						svp.flexible.h = rpEq(svp.flexible.h, -1, height)
					}

					ll.Infof("svp.flexible: %v, svp.start %v", svp.flexible, svp.start)

					// deep calculate subviews
					deepCalculate(svp, nextAnchor)

					// update cache after calculating subviews
					switch ui.types {
					case typesVStack:
						nextAnchor.start.y += svp.flexible.h
					case typesHStack:
						nextAnchor.start.x += svp.flexible.w
					}
				})

				if again {
					return true
				}
			}

			return false
		}

		again := true
		for again {
			again = setSubviews()
		}
	}

	deepCalculate(r.uiView, r.uiView)
	logs.Warnf("root element: %s (flexible: %v)", r.types, r.flexible)
	invokeSomeView(r.uiView.contents[0], func(ui *uiView) {
		logs.Warnf("first element: %s (flexible: %v)", ui.types, ui.flexible)
	})
}

func (*rootView) countFlexibleChildren(ui *uiView) (widthCount, heightCount, widthDelta, heightDelta, widthMax, heightMax int, recalculatedSubViews map[*uiView]bool) {
	sizedTable := make(map[*uiView]bool, len(ui.contents))
	wCount, hCount, wDelta, hDelta, wMax, hMax := 0, 0, 0, 0, 0, 0
	for _, sv := range ui.contents {
		invokeSomeView(sv, func(ui *uiView) {
			size := ui.GetFrameWithMargin()

			if size.w == -1 {
				wCount++
			} else {
				sizedTable[ui] = true
				wDelta += size.w
				wMax = max(wMax, size.w)
			}

			if size.h == -1 {
				hCount++
			} else {
				sizedTable[ui] = true
				hDelta += size.h
				hMax = max(hMax, size.h)
			}
		})
	}

	return wCount, hCount, wDelta, hDelta, wMax, hMax, sizedTable
}

// func (r *rootView) draw(screen *ebiten.Image) {
// 	r.uiViewBack.Draw(screen, func(screen *ebiten.Image) {
// 		r.uiViewBack.ApplyViewModifiers(screen)
// 	})
// }

// func (r *rootView) preCacheChildrenSize() {
// 	var formula func(SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int)
// 	formula = func(current SomeView) (maxWidth, maxHeight, sumWidth, sumHeight int) {
// 		currentView := current.view()
// 		initW, initH := current.initBounds()
// 		maxWidth, maxHeight = initW, initH
// 		for _, sv := range currentView.subviews {
// 			maxW, maxH, sumW, sumH := formula(sv)

// 			maxWidth = max(maxWidth, maxW)
// 			maxHeight = max(maxHeight, maxH)
// 			sumWidth += rpEq(sumW, -1, 0)
// 			sumHeight += rpEq(sumH, -1, 0)
// 		}

// 		sumWidth = max(sumWidth, initW)
// 		sumHeight = max(sumHeight, initH)

// 		switch currentView.types {
// 		case typesVStack:
// 			currentView.size.w = maxWidth
// 			currentView.size.h = rpEq(sumHeight, 0, -1)
// 		case typesHStack:
// 			currentView.size.w = rpEq(sumWidth, 0, -1)
// 			currentView.size.h = maxHeight
// 		case typesZStack:
// 			currentView.size.w = maxWidth
// 			currentView.size.h = maxHeight
// 		default:
// 			currentView.size.w = initW
// 			currentView.size.h = initH
// 		}

// 		return currentView.size.w, currentView.size.h, sumWidth, sumHeight
// 	}

// 	formula(r.uiViewBack)
// 	r.uiViewBack.size.w, r.uiViewBack.size.h = r.initBounds()
// }

// func (r *rootView) calculationParameters() {
// 	var formula func(v, anchor *uiViewBack)
// 	formula = func(v, anchor *uiViewBack) {
// 		// update params from anchor
// 		v.pos.x = anchor.pos.x
// 		v.pos.y = anchor.pos.y
// 		v.xx = anchor.xx
// 		v.yy = anchor.yy
// 		v.viewOpacity *= anchor.viewOpacity
// 		v.fColor = rpZero(v.fColor, anchor.fColor)
// 		v.fontSizes = rpZero(v.fontSizes, anchor.fontSizes)
// 		v.isPressing = rpZero(v.isPressing, anchor.isPressing)

// 		// apply view modifiers
// 		nextAnchor := v.Copy()
// 		nextAnchor.ApplyViewModifiers(nil)

// 		// update cache after modifying
// 		nextAnchor.pos.x += nextAnchor.padding.left
// 		nextAnchor.pos.y += nextAnchor.padding.top
// 		nextAnchor.xx = nextAnchor.padding.left // TODO: check if this is correct
// 		nextAnchor.yy = nextAnchor.padding.top

// 		// fmt.Printf("next anchor: %+v\n", nextAnchor)

// 		// calculate flexible size
// 		wFlexCount, hFlexCount, wDelta, hDelta, recalculatedSubViews := r.countFlexibleChildren(nextAnchor)

// 		setSubviews := func() (again bool) {
// 			if wFlexCount <= -1 || hFlexCount <= -1 {
// 				l := logs.Default().WithField("current", v.types)
// 				l.Fatalf("%s: flex count is negative: wFlexCount: %d, hFlexCount: %d", nextAnchor.types, wFlexCount, hFlexCount)
// 			}

// 			width := (nextAnchor.Width() - wDelta) / rpZero(wFlexCount, 1)
// 			height := (nextAnchor.Height() - hDelta) / rpZero(hFlexCount, 1)

// 			for _, sv := range nextAnchor.subviews {
// 				svp := sv.view()
// 				ll := logs.Default().WithField("parent", v.types).WithField("current", svp.types)
// 				ll.Debugf("x,y(%d, %d), xx,yy(%d, %d), w,h(%d, %d)", svp.pos.x, svp.pos.y, svp.xx, svp.yy, svp.size.w, svp.size.h)

// 				if !recalculatedSubViews[svp] {
// 					recalculatedSubViews[svp] = true
// 					// calculate that does width/height need to recalculate
// 					switch v.types {
// 					case typesVStack, typesHStack, typesZStack:
// 						if svp.size.w > width {
// 							ll.Debugf("width out of range! svp.size.w: %d > width: %d", svp.size.w, width)
// 							wFlexCount--
// 							wDelta += svp.size.w
// 							again = true
// 						}

// 						if svp.size.h > height {
// 							ll.Debugf("height out of range! svp.size.h: %d > height: %d", svp.size.h, height)
// 							hFlexCount--
// 							hDelta += svp.size.h
// 							again = true
// 						}
// 					}

// 					if again {
// 						return true
// 					}
// 				}

// 				// set size to subviews
// 				switch v.types {
// 				case typesVStack:
// 					svp.size.w = rpEq(svp.initSize.w, -1, max(svp.size.w, nextAnchor.Width()))
// 					svp.size.h = rpEq(svp.initSize.h, -1, max(svp.size.h, height))
// 				case typesHStack:
// 					svp.size.w = rpEq(svp.initSize.w, -1, max(svp.size.w, width))
// 					svp.size.h = rpEq(svp.initSize.h, -1, max(svp.size.h, nextAnchor.Height()))
// 				default:
// 					svp.size.w = rpEq(svp.initSize.w, -1, max(svp.size.w, width))
// 					svp.size.h = rpEq(svp.initSize.h, -1, max(svp.size.h, height))
// 				}

// 				// deep calculate subviews
// 				formula(svp, nextAnchor)

// 				// update cache after calculating subviews
// 				switch v.types {
// 				case typesVStack:
// 					nextAnchor.pos.y += svp.size.h
// 					nextAnchor.yy += svp.size.h
// 				case typesHStack:
// 					nextAnchor.pos.x += svp.size.w
// 					nextAnchor.xx += svp.size.w
// 				}
// 			}

// 			return false
// 		}

// 		again := true
// 		for again {
// 			again = setSubviews()
// 		}
// 	}

// 	formula(r.uiViewBack, r.uiViewBack)
// }

// func (r *rootView) countFlexibleChildren(v *uiViewBack) (widthCount, heightCount, widthDelta, heightDelta int, recalculatedSubViews map[*uiViewBack]bool) {
// 	table := make(map[*uiViewBack]bool, len(v.subviews))
// 	wCount, hCount, wDelta, hDelta := 0, 0, 0, 0
// 	for _, sv := range v.subviews {
// 		svp := sv.view()
// 		vvw, vvh := sv.initBounds()

// 		if vvw == -1 {
// 			wCount++
// 		} else {
// 			table[svp] = true
// 			wDelta += vvw
// 		}

// 		if vvh == -1 {
// 			hCount++
// 		} else {
// 			table[svp] = true
// 			hDelta += vvh
// 		}
// 	}

// 	return wCount, hCount, wDelta, hDelta, table
// }
