package view

import (
	. "github.com/yanun0323/ebui"
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
			KeepAspectRatio().
			Frame(blockSizeS).Debug()
	}

	imgM := func() SomeView {
		return Image(Const("gopher.png")).
			ScaleToFit().
			KeepAspectRatio().
			Frame(blockSizeM).Debug()
	}

	imgL := func() SomeView {
		return Image(Const("gopher.png")).
			ScaleToFit().
			KeepAspectRatio().
			Frame(blockSizeL).Debug()
	}

	bt := func(key string, align layout.Align) SomeView {
		return Button(key, func() {
			v.alignment.Set(align)
		}).Padding()
	}

	sp := Bind(0.0).Animated()

	return HStack(
		VStack(
			Text("VStack"),
			HStack(
				VStack(
					imgS(),
					imgM(),
					imgL(),
				).Spacing(sp).Modify(v.stackWrapper),
			),
		),

		VStack(
			Text("HStack"),
			HStack(
				imgS(),
				imgM(),
				imgL(),
			).Spacing(sp).Modify(v.stackWrapper),
		),

		VStack(
			Text("Align"),
			HStack(
				bt("Top", layout.AlignTop),
				bt("Bottom", layout.AlignBottom),
			),
			HStack(
				bt("Leading", layout.AlignLeading),
				bt("Trailing", layout.AlignTrailing),
			),

			Rectangle().Frame(Bind(NewSize(5))),

			bt("Center", layout.AlignCenter),
			bt("CenterHorizontal", layout.AlignCenterHorizontal),
			bt("CenterVertical", layout.AlignCenterVertical),

			Rectangle().Frame(Bind(NewSize(5))),

			bt("Reset", layout.AlignDefault),
		),
	).Center().Align(Const(layout.AlignCenter))
}

func (v *layoutView) stackWrapper(vv SomeView) SomeView {
	return vv.
		Frame(Bind(NewSize(300))).
		Border(Bind(NewInset(1))).
		Align(v.alignment)
}
