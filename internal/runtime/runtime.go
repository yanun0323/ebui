package runtime

import (
	"math"

	ui "github.com/yanun0323/ebui"
)

type Runtime struct {
	measurer Measurer
	rootView ui.View
	root     *node
	size     ui.Size
	frame    Frame

	focused     *node
	captured    *node
	lastPointer ui.Point
	dirtyLayout bool
	dirtyFrame  bool
}

func New(root ui.View, measurer Measurer) *Runtime {
	rt := &Runtime{
		measurer:    measurer,
		rootView:    root,
		dirtyLayout: true,
		dirtyFrame:  true,
	}
	if root != nil {
		rt.root = newNode(rt, nil, ui.Render(root))
		rt.root.markDirtyView()
		rt.refresh()
	}
	return rt
}

func (rt *Runtime) Resize(size ui.Size) {
	if rt.size == size {
		return
	}
	rt.size = size
	rt.dirtyLayout = true
	rt.dirtyFrame = true
}

func (rt *Runtime) Update(events []Event) {
	if rt.root == nil {
		return
	}
	for _, event := range events {
		rt.dispatch(event)
	}
	rt.refresh()
}

func (rt *Runtime) Frame() Frame {
	return rt.frame
}

func (rt *Runtime) refresh() {
	if rt.root == nil {
		return
	}
	rt.rebuildDirty(rt.root)
	if rt.dirtyLayout {
		rt.relayout()
	}
	if rt.dirtyFrame {
		rt.rebuildFrame()
	}
}

func (rt *Runtime) rebuildDirty(current *node) {
	if current == nil {
		return
	}
	if current.kind == ui.ElementKindView && current.dirtyView {
		current.cleanupDependencies()
		ctx := &buildContext{
			owner: current,
			seen:  make(map[uintptr]struct{}),
		}
		built := current.view.Body(ctx)
		current.buildCount++
		current.children = rt.reconcileChildren(current, current.children, []ui.Element{built})
		current.dirtyView = false
		current.markLayoutDirty()
	}

	for _, child := range current.children {
		rt.rebuildDirty(child)
	}
}

func (rt *Runtime) reconcileChildren(parent *node, existing []*node, next []ui.Element) []*node {
	result := make([]*node, 0, len(next))
	for i, element := range next {
		if i < len(existing) && rt.compatible(existing[i], element) {
			existing[i].updateFromElement(element)
			existing[i].children = rt.reconcileChildren(existing[i], existing[i].children, element.Children())
			result = append(result, existing[i])
			continue
		}

		mounted := newNode(rt, parent, element)
		mounted.children = rt.reconcileChildren(mounted, nil, element.Children())
		result = append(result, mounted)
	}

	for i := len(next); i < len(existing); i++ {
		existing[i].destroy()
	}

	return result
}

func (rt *Runtime) compatible(existing *node, next ui.Element) bool {
	if existing.kind != next.Kind() {
		return false
	}
	if existing.kind == ui.ElementKindView {
		return existing.viewCompatible(next.View())
	}
	return true
}

func (n *node) destroy() {
	n.cleanupDependencies()
	for _, child := range n.children {
		child.destroy()
	}
}

func (rt *Runtime) relayout() {
	if rt.root == nil {
		return
	}
	constraints := ui.Flexible(rt.size)
	if rt.size.IsZero() {
		constraints = ui.Flexible(ui.S(ui.Inf, ui.Inf))
	}
	size := rt.measureNode(rt.root, constraints, defaultResolvedStyle())
	if !rt.size.IsZero() {
		size = rt.size
	}
	rt.arrangeNode(rt.root, ui.R(0, 0, size.Width, size.Height), defaultResolvedStyle())
	rt.dirtyLayout = false
}

func (rt *Runtime) measureNode(n *node, constraints ui.Constraints, parentStyle resolvedStyle) ui.Size {
	if n == nil {
		return ui.Size{}
	}

	style := n.element.Style()
	resolved := n.resolve(parentStyle)
	chrome := n.chromeInsets()
	innerConstraints := constraints.Shrink(chrome)
	var content ui.Size

	if n.kind == ui.ElementKindView {
		if len(n.children) > 0 {
			content = rt.measureNode(n.children[0], innerConstraints, resolved)
		}
	} else {
		switch n.kind {
		case ui.ElementKindEmpty:
			content = ui.Size{}
		case ui.ElementKindText:
			props := n.props.(ui.TextProps)
			content = rt.measurer.MeasureText(props.Content, resolved.Font, resolved.LineLimit)
		case ui.ElementKindImage:
			props := n.props.(ui.ImageProps)
			content = rt.measurer.ImageSize(props.Source)
			if style.HasImageFit && innerConstraints.Max.Width != ui.Inf && innerConstraints.Max.Height != ui.Inf {
				content = innerConstraints.Max
			}
		case ui.ElementKindRectangle, ui.ElementKindCircle:
			content = innerConstraints.Min
		case ui.ElementKindDivider:
			content = rt.measureDivider(n, innerConstraints)
		case ui.ElementKindSpacer:
			content = ui.Size{}
		case ui.ElementKindVStack, ui.ElementKindHStack:
			content = rt.measureStack(n, innerConstraints, resolved)
		case ui.ElementKindZStack:
			content = rt.measureZStack(n, innerConstraints, resolved)
		case ui.ElementKindScrollView:
			content = rt.measureScrollView(n, innerConstraints, resolved)
		case ui.ElementKindButton:
			content = rt.measureButton(n, innerConstraints, resolved)
		case ui.ElementKindToggle:
			content = rt.measureToggle(n, innerConstraints, resolved)
		case ui.ElementKindSlider:
			content = rt.measureSlider(n, innerConstraints)
		default:
			content = ui.Size{}
		}
	}

	total := content.AddInsets(chrome)
	if style.HasFrame {
		total = style.Frame
	}
	n.measured = constraints.Constrain(total)
	return n.measured
}

func (rt *Runtime) measureDivider(n *node, constraints ui.Constraints) ui.Size {
	switch n.layoutKind() {
	case ui.ElementKindHStack:
		height := constraints.Max.Height
		if height == ui.Inf {
			height = 1
		}
		return ui.S(1, height)
	case ui.ElementKindVStack:
		width := constraints.Max.Width
		if width == ui.Inf {
			width = 1
		}
		return ui.S(width, 1)
	default:
		return ui.S(1, 1)
	}
}

func (n *node) layoutKind() ui.ElementKind {
	current := n.parent
	for current != nil {
		if current.kind != ui.ElementKindView {
			return current.kind
		}
		current = current.parent
	}
	return ui.ElementKindEmpty
}

func (rt *Runtime) measureStack(n *node, constraints ui.Constraints, parent resolvedStyle) ui.Size {
	style := n.element.Style()
	spacing := 0.0
	if style.HasSpacing {
		spacing = style.Spacing
	}

	totalMain := 0.0
	maxCross := 0.0
	visible := 0

	for _, child := range n.children {
		childConstraints := constraints
		if n.kind == ui.ElementKindVStack {
			childConstraints.Max.Height = ui.Inf
		} else {
			childConstraints.Max.Width = ui.Inf
		}
		size := rt.measureNode(child, childConstraints, parent)
		if child.kind != ui.ElementKindSpacer {
			if visible > 0 {
				totalMain += spacing
			}
			visible++
		}
		if n.kind == ui.ElementKindVStack {
			totalMain += size.Height
			maxCross = math.Max(maxCross, size.Width)
		} else {
			totalMain += size.Width
			maxCross = math.Max(maxCross, size.Height)
		}
	}

	if n.kind == ui.ElementKindVStack {
		return constraints.Constrain(ui.S(maxCross, totalMain))
	}
	return constraints.Constrain(ui.S(totalMain, maxCross))
}

func (rt *Runtime) measureZStack(n *node, constraints ui.Constraints, parent resolvedStyle) ui.Size {
	var max ui.Size
	for _, child := range n.children {
		size := rt.measureNode(child, constraints, parent)
		max.Width = math.Max(max.Width, size.Width)
		max.Height = math.Max(max.Height, size.Height)
	}
	return constraints.Constrain(max)
}

func (rt *Runtime) measureScrollView(n *node, constraints ui.Constraints, parent resolvedStyle) ui.Size {
	if len(n.children) == 0 {
		return ui.Size{}
	}
	contentConstraints := constraints
	props := n.props.(ui.ScrollViewProps)
	if props.Direction == ui.DirectionVertical {
		contentConstraints.Max.Height = ui.Inf
	} else {
		contentConstraints.Max.Width = ui.Inf
	}
	content := rt.measureNode(n.children[0], contentConstraints, parent)
	if props.Direction == ui.DirectionVertical {
		return constraints.Constrain(ui.S(content.Width, constraints.Max.Height))
	}
	return constraints.Constrain(ui.S(constraints.Max.Width, content.Height))
}

func (rt *Runtime) measureButton(n *node, constraints ui.Constraints, parent resolvedStyle) ui.Size {
	if len(n.children) == 0 {
		return constraints.Constrain(ui.S(96, 36))
	}
	return rt.measureNode(n.children[0], constraints, parent)
}

func (rt *Runtime) measureToggle(n *node, constraints ui.Constraints, parent resolvedStyle) ui.Size {
	base := ui.S(52, 32)
	if len(n.children) == 0 {
		return constraints.Constrain(base)
	}
	label := rt.measureNode(n.children[0], constraints, parent)
	return constraints.Constrain(ui.S(label.Width+68, math.Max(label.Height, base.Height)))
}

func (rt *Runtime) measureSlider(n *node, constraints ui.Constraints) ui.Size {
	return constraints.Constrain(ui.S(240, 32))
}

func (rt *Runtime) arrangeNode(n *node, rect ui.Rect, parentStyle resolvedStyle) {
	if n == nil {
		return
	}

	style := n.element.Style()
	resolved := n.resolve(parentStyle)
	rect = rect.Offset(n.localOffset())
	n.frame = rect
	n.content = rect.Inset(n.chromeInsets())

	switch n.kind {
	case ui.ElementKindView:
		if len(n.children) > 0 {
			rt.arrangeNode(n.children[0], n.content, resolved)
		}
	case ui.ElementKindVStack, ui.ElementKindHStack:
		rt.arrangeStack(n, resolved)
	case ui.ElementKindZStack:
		for _, child := range n.children {
			rt.arrangeAlignedChild(child, n.content, resolved, n.childAlignment(ui.AlignCenter))
		}
	case ui.ElementKindScrollView:
		rt.arrangeScrollView(n, resolved)
	case ui.ElementKindButton:
		if len(n.children) > 0 {
			rt.arrangeAlignedChild(n.children[0], n.content, resolved, n.childAlignment(ui.AlignCenter))
		}
	case ui.ElementKindToggle:
		if len(n.children) > 0 {
			labelRect := ui.R(n.content.Min.X, n.content.Min.Y, math.Max(0, n.content.Width()-68), n.content.Height())
			rt.arrangeNode(n.children[0], labelRect, resolved)
		}
	default:
		for _, child := range n.children {
			rt.arrangeNode(child, n.content, resolved)
		}
	}

	if style.HasFrame {
		n.frame = ui.R(rect.Min.X, rect.Min.Y, style.Frame.Width, style.Frame.Height)
		n.content = n.frame.Inset(n.chromeInsets())
	}
}

func (rt *Runtime) arrangeStack(n *node, parent resolvedStyle) {
	style := n.element.Style()
	spacing := 0.0
	if style.HasSpacing {
		spacing = style.Spacing
	}

	totalGrow := 0.0
	used := 0.0
	visible := 0
	for _, child := range n.children {
		if child.kind == ui.ElementKindSpacer {
			props := child.props.(ui.SpacerProps)
			totalGrow += props.Grow
			continue
		}
		if visible > 0 {
			used += spacing
		}
		visible++
		if n.kind == ui.ElementKindVStack {
			used += child.measured.Height
		} else {
			used += child.measured.Width
		}
	}

	remaining := 0.0
	if n.kind == ui.ElementKindVStack {
		remaining = math.Max(0, n.content.Height()-used)
	} else {
		remaining = math.Max(0, n.content.Width()-used)
	}

	cursor := n.content.Min
	for index, child := range n.children {
		if index > 0 && child.kind != ui.ElementKindSpacer {
			if n.kind == ui.ElementKindVStack {
				cursor.Y += spacing
			} else {
				cursor.X += spacing
			}
		}

		size := child.measured
		if child.kind == ui.ElementKindSpacer && totalGrow > 0 {
			props := child.props.(ui.SpacerProps)
			if n.kind == ui.ElementKindVStack {
				size.Height = remaining * (props.Grow / totalGrow)
				size.Width = n.content.Width()
			} else {
				size.Width = remaining * (props.Grow / totalGrow)
				size.Height = n.content.Height()
			}
		}

		rect := ui.R(cursor.X, cursor.Y, size.Width, size.Height)
		if child.kind != ui.ElementKindSpacer {
			rt.arrangeAlignedChild(child, rect, parent, n.childAlignment(ui.AlignTopLeading))
		} else {
			rt.arrangeNode(child, rect, parent)
		}

		if n.kind == ui.ElementKindVStack {
			cursor.Y += size.Height
		} else {
			cursor.X += size.Width
		}
	}
}

func (rt *Runtime) arrangeScrollView(n *node, parent resolvedStyle) {
	if len(n.children) == 0 {
		return
	}
	child := n.children[0]
	rect := n.content
	props := n.props.(ui.ScrollViewProps)
	if props.Direction == ui.DirectionVertical {
		rect.Max.Y = rect.Min.Y + child.measured.Height
		rect = rect.Offset(ui.P(0, -n.scrollOffset.Y))
	} else {
		rect.Max.X = rect.Min.X + child.measured.Width
		rect = rect.Offset(ui.P(-n.scrollOffset.X, 0))
	}
	rt.arrangeNode(child, rect, parent)
}

func (rt *Runtime) arrangeAlignedChild(child *node, slot ui.Rect, parent resolvedStyle, alignment ui.Alignment) {
	size := child.measured
	x := slot.Min.X
	y := slot.Min.Y

	switch alignment {
	case ui.AlignCenter:
		x = slot.Min.X + (slot.Width()-size.Width)/2
		y = slot.Min.Y + (slot.Height()-size.Height)/2
	case ui.AlignTop:
		x = slot.Min.X + (slot.Width()-size.Width)/2
	case ui.AlignTopTrailing:
		x = slot.Max.X - size.Width
	case ui.AlignLeading:
		y = slot.Min.Y + (slot.Height()-size.Height)/2
	case ui.AlignTrailing:
		x = slot.Max.X - size.Width
		y = slot.Min.Y + (slot.Height()-size.Height)/2
	case ui.AlignBottomLeading:
		y = slot.Max.Y - size.Height
	case ui.AlignBottom:
		x = slot.Min.X + (slot.Width()-size.Width)/2
		y = slot.Max.Y - size.Height
	case ui.AlignBottomTrailing:
		x = slot.Max.X - size.Width
		y = slot.Max.Y - size.Height
	}

	rt.arrangeNode(child, ui.R(x, y, size.Width, size.Height), parent)
}

func (rt *Runtime) rebuildFrame() {
	frame := Frame{}
	rt.buildFrame(rt.root, defaultResolvedStyle(), ui.Rect{}, &frame)
	rt.frame = frame
	rt.dirtyFrame = false
}

func (rt *Runtime) buildFrame(n *node, parent resolvedStyle, clip ui.Rect, frame *Frame) {
	if n == nil {
		return
	}
	resolved := n.resolve(parent)
	style := n.element.Style()
	currentClip := clip
	if n.kind == ui.ElementKindScrollView {
		currentClip = n.content
	}

	if style.HasBackground {
		frame.Ops = append(frame.Ops, DrawOp{
			Kind:    DrawKindRect,
			Rect:    n.frame,
			Fill:    style.Background,
			Opacity: resolved.Opacity,
			Clip:    currentClip,
			HasClip: !currentClip.Size().IsZero(),
		})
	}

	switch n.kind {
	case ui.ElementKindRectangle:
		fill := style.Fill
		if !style.HasFill {
			fill = ui.RGBA(200, 200, 200, 255)
		}
		frame.Ops = append(frame.Ops, DrawOp{
			Kind:    DrawKindRect,
			Rect:    n.content,
			Fill:    fill,
			Opacity: resolved.Opacity,
			Clip:    currentClip,
			HasClip: !currentClip.Size().IsZero(),
		})
	case ui.ElementKindCircle:
		fill := style.Fill
		if !style.HasFill {
			fill = ui.RGBA(200, 200, 200, 255)
		}
		frame.Ops = append(frame.Ops, DrawOp{
			Kind:    DrawKindCircle,
			Rect:    n.content,
			Fill:    fill,
			Opacity: resolved.Opacity,
			Clip:    currentClip,
			HasClip: !currentClip.Size().IsZero(),
		})
	case ui.ElementKindText:
		props := n.props.(ui.TextProps)
		frame.Ops = append(frame.Ops, DrawOp{
			Kind:       DrawKindText,
			Rect:       n.content,
			Text:       props.Content,
			Font:       resolved.Font,
			Foreground: resolved.Foreground,
			Opacity:    resolved.Opacity,
			Clip:       currentClip,
			HasClip:    !currentClip.Size().IsZero(),
		})
	case ui.ElementKindImage:
		props := n.props.(ui.ImageProps)
		frame.Ops = append(frame.Ops, DrawOp{
			Kind:    DrawKindImage,
			Rect:    n.content,
			Image:   props.Source,
			Opacity: resolved.Opacity,
			Clip:    currentClip,
			HasClip: !currentClip.Size().IsZero(),
		})
	case ui.ElementKindDivider:
		frame.Ops = append(frame.Ops, DrawOp{
			Kind:    DrawKindRect,
			Rect:    n.content,
			Fill:    ui.RGBA(210, 216, 227, 255),
			Opacity: resolved.Opacity,
			Clip:    currentClip,
			HasClip: !currentClip.Size().IsZero(),
		})
	case ui.ElementKindToggle:
		rt.buildToggleFrame(n, currentClip, resolved, frame)
	case ui.ElementKindSlider:
		rt.buildSliderFrame(n, currentClip, resolved, frame)
	}

	for _, child := range n.children {
		rt.buildFrame(child, resolved, currentClip, frame)
	}
}

func (rt *Runtime) buildToggleFrame(n *node, clip ui.Rect, resolved resolvedStyle, frame *Frame) {
	props := n.props.(ui.ToggleProps)
	value := false
	if props.State != nil {
		value = props.State.Value()
	}
	track := ui.R(n.content.Max.X-52, n.content.Min.Y+(n.content.Height()-28)/2, 52, 28)
	fill := ui.RGBA(198, 207, 225, 255)
	if value {
		fill = ui.RGBA(58, 111, 219, 255)
	}
	frame.Ops = append(frame.Ops, DrawOp{
		Kind:    DrawKindRect,
		Rect:    track,
		Fill:    fill,
		Opacity: resolved.Opacity,
		Clip:    clip,
		HasClip: !clip.Size().IsZero(),
	})
	knobX := track.Min.X + 4
	if value {
		knobX = track.Max.X - 24
	}
	frame.Ops = append(frame.Ops, DrawOp{
		Kind:    DrawKindCircle,
		Rect:    ui.R(knobX, track.Min.Y+4, 20, 20),
		Fill:    ui.White,
		Opacity: resolved.Opacity,
		Clip:    clip,
		HasClip: !clip.Size().IsZero(),
	})
}

func (rt *Runtime) buildSliderFrame(n *node, clip ui.Rect, resolved resolvedStyle, frame *Frame) {
	props := n.props.(ui.SliderProps)
	value := props.Min
	if props.State != nil {
		value = props.State.Value()
	}
	ratio := 0.0
	if props.Max > props.Min {
		ratio = (value - props.Min) / (props.Max - props.Min)
	}
	ratio = clamp(ratio, 0, 1)
	track := ui.R(n.content.Min.X, n.content.Min.Y+n.content.Height()/2-3, n.content.Width(), 6)
	frame.Ops = append(frame.Ops, DrawOp{
		Kind:    DrawKindRect,
		Rect:    track,
		Fill:    ui.RGBA(212, 220, 232, 255),
		Opacity: resolved.Opacity,
		Clip:    clip,
		HasClip: !clip.Size().IsZero(),
	})
	fillWidth := track.Width() * ratio
	frame.Ops = append(frame.Ops, DrawOp{
		Kind:    DrawKindRect,
		Rect:    ui.R(track.Min.X, track.Min.Y, fillWidth, track.Height()),
		Fill:    ui.RGBA(58, 111, 219, 255),
		Opacity: resolved.Opacity,
		Clip:    clip,
		HasClip: !clip.Size().IsZero(),
	})
	knobX := track.Min.X + fillWidth - 8
	frame.Ops = append(frame.Ops, DrawOp{
		Kind:    DrawKindCircle,
		Rect:    ui.R(knobX, track.Min.Y-7, 16, 16),
		Fill:    ui.White,
		Opacity: resolved.Opacity,
		Clip:    clip,
		HasClip: !clip.Size().IsZero(),
	})
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (rt *Runtime) dispatch(event Event) {
	rt.lastPointer = event.Point
	switch event.Kind {
	case EventPointerDown:
		target := rt.hitTest(rt.root, event.Point)
		rt.captured = target
		if target != nil {
			target.pressed = true
			if target.kind == ui.ElementKindSlider {
				target.dragging = true
				rt.updateSlider(target, event.Point)
			}
		}
	case EventPointerMove:
		if rt.captured != nil && rt.captured.dragging && rt.captured.kind == ui.ElementKindSlider {
			rt.updateSlider(rt.captured, event.Point)
		}
	case EventPointerUp:
		target := rt.hitTest(rt.root, event.Point)
		if rt.captured != nil {
			rt.handlePointerRelease(rt.captured, target, event.Point)
			rt.captured.pressed = false
			rt.captured.dragging = false
			rt.captured = nil
		}
	case EventWheel:
		target := rt.hitScroll(rt.root, event.Point)
		if target != nil {
			rt.updateScroll(target, event.Scroll)
		}
	}
}

func (rt *Runtime) handlePointerRelease(captured, target *node, point ui.Point) {
	if captured.resolved.Disabled {
		return
	}
	switch captured.kind {
	case ui.ElementKindButton:
		if captured == target {
			props := captured.props.(ui.ButtonProps)
			if props.OnPress != nil {
				props.OnPress()
			}
		}
	case ui.ElementKindToggle:
		if captured == target {
			props := captured.props.(ui.ToggleProps)
			if props.State != nil {
				props.State.Set(!props.State.Value())
			}
		}
	case ui.ElementKindSlider:
		rt.updateSlider(captured, point)
	}
}

func (rt *Runtime) updateScroll(target *node, scroll ui.Point) {
	if len(target.children) == 0 {
		return
	}
	child := target.children[0]
	props := target.props.(ui.ScrollViewProps)
	if props.Direction == ui.DirectionVertical {
		max := math.Max(0, child.measured.Height-target.content.Height())
		target.scrollOffset.Y = clamp(target.scrollOffset.Y-scroll.Y*24, 0, max)
	} else {
		max := math.Max(0, child.measured.Width-target.content.Width())
		target.scrollOffset.X = clamp(target.scrollOffset.X-scroll.X*24, 0, max)
	}
	target.markLayoutDirty()
}

func (rt *Runtime) updateSlider(target *node, point ui.Point) {
	props := target.props.(ui.SliderProps)
	if props.State == nil {
		return
	}
	width := target.content.Width()
	if width <= 0 {
		return
	}
	ratio := clamp((point.X-target.content.Min.X)/width, 0, 1)
	value := props.Min + ratio*(props.Max-props.Min)
	props.State.Set(value)
}

func (rt *Runtime) hitTest(current *node, point ui.Point) *node {
	if current == nil || !current.frame.Contains(point) {
		return nil
	}

	for i := len(current.children) - 1; i >= 0; i-- {
		if hit := rt.hitTest(current.children[i], point); hit != nil {
			return hit
		}
	}

	switch current.kind {
	case ui.ElementKindButton, ui.ElementKindToggle, ui.ElementKindSlider:
		return current
	default:
		return nil
	}
}

func (rt *Runtime) hitScroll(current *node, point ui.Point) *node {
	if current == nil || !current.frame.Contains(point) {
		return nil
	}
	for i := len(current.children) - 1; i >= 0; i-- {
		if hit := rt.hitScroll(current.children[i], point); hit != nil {
			return hit
		}
	}
	if current.kind == ui.ElementKindScrollView {
		return current
	}
	return nil
}

func FindFirstView(root *node, view ui.View) *node {
	if root == nil || view == nil {
		return nil
	}
	if root.view != nil && root.viewCompatible(view) {
		return root
	}
	for _, child := range root.children {
		if found := FindFirstView(child, view); found != nil {
			return found
		}
	}
	return nil
}
