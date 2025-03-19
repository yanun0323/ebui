package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yanun0323/ebui/font"
	layout "github.com/yanun0323/ebui/layout"
)

// Frame: internal bounds (exclude padding, border ...)
// Bounds: external bounds (include padding, border ...)

// View represents the base of all views
type View interface {
	Body() SomeView
}

// ViewModifier represents the base of all view modifiers
type ViewModifier interface {
	Body(SomeView) SomeView
}

// SomeView represents the instance of View
type SomeView interface {
	View
	eventHandler

	// preload sets the environment variables and returns the layout information of the view
	preload(parent *viewCtxEnv, stackTypes ...stackType) (preloadData, layoutFunc)

	// drawOption returns the draw options of the view
	drawOption(rect CGRect, hook ...func(*ebiten.DrawImageOptions)) *ebiten.DrawImageOptions
	// draw draws the view
	draw(screen *ebiten.Image, hook ...func(*ebiten.DrawImageOptions))

	// debugPrint prints debug information in the console
	debugPrint(prefix string, format string, a ...any)

	// setEnv sets the environment variables of the view
	setEnv(env *viewCtxEnv)

	// bytes returns the bytes of the view
	bytes() []byte

	// count returns the count of the view
	count() int

	// userSetFrameSize returns the frame size from the user
	userSetFrameSize() CGSize
	// systemSetFrame returns the frame from the system
	systemSetFrame() CGRect
	// systemSetBounds returns the bounds from the system
	systemSetBounds() CGRect

	// DebugPrint makes the view print debug information in the console
	DebugPrint(tag string) SomeView

	// Frame sets the frame of the view
	Frame(*Binding[CGSize]) SomeView

	// Padding sets the padding of the view
	Padding(inset ...*Binding[CGInset]) SomeView

	// ForegroundColor sets the foreground color of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	ForegroundColor(color *Binding[CGColor]) SomeView

	// BackgroundColor sets the background color of the view
	BackgroundColor(color *Binding[CGColor]) SomeView

	// Fill sets the background color of the view
	Fill(color *Binding[CGColor]) SomeView

	// Font sets the font of the view
	Font(size *Binding[font.Size], weight *Binding[font.Weight]) SomeView

	// FontSize sets the font size of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	FontSize(size *Binding[font.Size]) SomeView

	// FontWeight sets the font weight of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	FontWeight(weight *Binding[font.Weight]) SomeView

	// FontLineHeight sets the line height of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	FontLineHeight(height *Binding[float64]) SomeView

	// FontKerning sets the kerning of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	FontKerning(spacing *Binding[float64]) SomeView

	// FontAlignment sets the alignment of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	FontAlignment(alignment *Binding[font.TextAlign]) SomeView

	// FontItalic sets the italic of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	FontItalic(italic ...*Binding[bool]) SomeView

	// RoundCorner sets the round corner of the view
	RoundCorner(radius ...*Binding[float64]) SomeView

	// ScaleToFit sets the scale to fit of the view
	ScaleToFit(enable ...*Binding[bool]) SomeView

	// KeepAspectRatio sets the keep aspect ratio of the view
	KeepAspectRatio(enable ...*Binding[bool]) SomeView

	// Border sets the border of the view
	Border(border *Binding[CGInset], color ...*Binding[CGColor]) SomeView

	// Opacity sets the opacity of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	Opacity(opacity *Binding[float64]) SomeView

	// Modifier sets the modifier of the view
	Modifier(ViewModifier) SomeView

	// Modify sets the modifier of the view
	Modify(with func(SomeView) SomeView) SomeView

	// Disabled sets the disabled of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	Disabled(disabled ...*Binding[bool]) SomeView

	// Align sets the alignment of the view
	//
	// It is a environment variable, so it will be inherited by subviews.
	Align(alignment *Binding[layout.Align]) SomeView

	// Center is a convenience method for:
	//
	//	 VStack(
	//		Spacer(),
	//		HStack(
	//			Spacer(),
	//			view,
	//			Spacer(),
	//		),
	//		Spacer(),
	//	 )
	Center() SomeView

	// Debug sets a red border around the view
	Debug() SomeView

	// Offset sets the offset of the view
	Offset(offset *Binding[CGPoint]) SomeView

	// Spacing sets the spacing of the view
	Spacing(spacing ...*Binding[float64]) SomeView

	// ScrollViewDirection sets the direction of the view
	ScrollViewDirection(direction *Binding[layout.Direction]) SomeView
}

type alignFunc func(offset CGPoint)

// layoutFunc: used to set the position and size of the child View, and return the actual occupied space
//
//	start: the starting coordinate of this child View
//	childBoundsSize: the external bounds of this child View
//	bounds: the actual occupied space (include padding, border ...)
type layoutFunc func(start CGPoint, childBoundsSize CGSize) (bounds CGRect, alignFunc alignFunc)

func newPreloadData(frameSize CGSize, padding, border CGInset) preloadData {
	return preloadData{
		FrameSize:   frameSize,
		IsInfWidth:  frameSize.IsInfWidth(),
		IsInfHeight: frameSize.IsInfHeight(),
		Padding:     padding,
		Border:      border,
	}
}

type preloadData struct {
	FrameSize   CGSize // the actual internal bounds size
	IsInfWidth  bool
	IsInfHeight bool
	Padding     CGInset // set by Padding
	Border      CGInset // set by Border
}

func (p *preloadData) BoundsSize() CGSize {
	return p.FrameSize.Expand(p.Padding).Expand(p.Border)
}
