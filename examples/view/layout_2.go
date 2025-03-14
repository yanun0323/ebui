package view

import (
	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/font"
	"github.com/yanun0323/ebui/layout"
)

var _ View = (*layoutView2)(nil)

func LayoutView2() View {
	return &layoutView2{
		alignment: Bind(layout.AlignDefault).Animated(),
	}
}

type layoutView2 struct {
	alignment *Binding[layout.Align]
}

func (v *layoutView2) Body() SomeView {
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

	tg := func(align layout.Align) SomeView {
		return HStack(
			Spacer(),
			Text(align.String()),
			Toggle(
				BindTwoWay(
					v.alignment,
					func(a layout.Align) bool {
						if a == layout.AlignDefault {
							a = layout.AlignTopLeading
						}

						return align == a
					}, func(b bool) layout.Align {
						if b {
							return align
						}
						a := v.alignment.Get()
						return a &^ align
					},
				),
			),
		).Spacing(Const(5.0)).Frame(Bind(NewSize(230, 30))).Padding(Bind(NewInset(5)))
	}

	return VStack(
		HStack(
			VStack(
				Text("VStack").Padding(),
				HStack(
					VStack(
						imgS(),
						imgM(),
						imgL(),
					).Modify(v.stackWrapper),
				),
			),
			VStack(
				Text("HStack").Padding(),
				HStack(
					imgS(),
					imgM(),
					imgL(),
				).Modify(v.stackWrapper),
			),
		).Spacing(Const(60.0)),

		VStack(
			Text("Align").Padding(),
			VStack(
				HStack(
					tg(layout.AlignTopLeading),
					tg(layout.AlignTopCenter),
					tg(layout.AlignTopTrailing),
				),
				HStack(
					tg(layout.AlignLeadingCenter),
					tg(layout.AlignCenter),
					tg(layout.AlignTrailingCenter),
				),
				HStack(
					tg(layout.AlignBottomLeading),
					tg(layout.AlignBottomCenter),
					tg(layout.AlignBottomTrailing),
				),
			).FontSize(Const(font.SubHeadline)).Spacing(),
		),
	).Spacing(Bind(30.0)).Center().Align(Const(layout.AlignCenter)).FontSize(Const(font.Headline))
}

func (v *layoutView2) stackWrapper(vv SomeView) SomeView {
	return vv.
		Frame(Bind(NewSize(300))).
		Border(Bind(NewInset(1))).
		Align(v.alignment)
}
