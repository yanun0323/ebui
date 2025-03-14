package view

import (
	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/layout"
)

var _ View = (*layoutView)(nil)

func LayoutView() View {
	return &layoutView{
		alignment: Bind(layout.AlignDefault).Animated(),
	}
}

type layoutView struct {
	alignment *Binding[layout.Align]
}

func (v *layoutView) Body() SomeView {
	blockSizeS := Bind(NewSize(40))
	blockSizeM := Bind(NewSize(60))
	blockSizeL := Bind(NewSize(80))

	imgS := func() SomeView {
		return Image(Const("gopher.png")).
			ScaleToFit().
			Frame(blockSizeS).Debug()
	}

	imgM := func() SomeView {
		return Image(Const("gopher.png")).
			ScaleToFit().
			Frame(blockSizeM).Debug()
	}

	imgL := func() SomeView {
		return Image(Const("gopher.png")).
			ScaleToFit().
			Frame(blockSizeL).Debug()
	}

	bt := func(align layout.Align) SomeView {
		return Button("", func() {
			v.alignment.Set(align)
		}, func() SomeView {
			return Text(align.String()).
				Center().
				Frame(Bind(NewSize(150, 33))).
				BackgroundColor(Bind(NewColor(64))).
				RoundCorner(Bind(10.0))
		}).Padding(Bind(NewInset(5)))
	}

	return VStack(
		HStack(
			VStack(
				Text("VStack"),
				HStack(
					VStack(
						imgS(),
						imgM(),
						imgL(),
					).Modify(v.stackWrapper),
				),
			),
			VStack(
				Text("HStack"),
				HStack(
					imgS(),
					imgM(),
					imgL(),
				).Modify(v.stackWrapper),
			),
		).Spacing(Const(60.0)).FontSize(Const(font.Title)),

		VStack(
			Text("Align").FontSize(Const(font.Headline)),
			VStack(
				HStack(
					bt(layout.AlignTopLeading),
					bt(layout.AlignTopCenter),
					bt(layout.AlignTopTrailing),
				),
				HStack(
					bt(layout.AlignLeadingCenter),
					bt(layout.AlignCenter),
					bt(layout.AlignTrailingCenter),
				),
				HStack(
					bt(layout.AlignBottomLeading),
					bt(layout.AlignBottomCenter),
					bt(layout.AlignBottomTrailing),
				),
			).Spacing(),
		).FontSize(Const(font.SubHeadline)),
	).Spacing(Bind(60.0)).Center().Align(Const(layout.AlignCenter))
}

func (v *layoutView) stackWrapper(vv SomeView) SomeView {
	return vv.
		Frame(Bind(NewSize(300))).
		Border(Bind(NewInset(1))).
		Align(v.alignment)
}
