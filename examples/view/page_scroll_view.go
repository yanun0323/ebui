package view

import (
	"strconv"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/layout"
)

func PageScrollView() View {
	return &pageScrollView{}
}

type pageScrollView struct{}

func (v *pageScrollView) Body() SomeView {
	rect := func(i int) SomeView {
		return Text(Const(strconv.Itoa(i))).
			Align(Bind(layout.AlignCenter)).
			Frame(Bind(NewSize(100))).
			BackgroundColor(Bind(NewColor(128, 0, 0))).
			Shadow()
	}

	enum := func(count int) []View {
		res := make([]View, 0, count)
		for i := range count {
			res = append(res, rect(i))
		}
		return res
	}

	value := Bind(10.0)
	return HStack(
		ScrollView(
			VStack(
				enum(10)...,
			).Spacing().Padding(),
		).Align(Bind(layout.AlignCenter)).
			ScrollViewDirection(Const(layout.DirectionVertical)).Debug(),
		ScrollView(
			HStack(
				enum(10)...,
			).Spacing().Padding(),
		).Align(Bind(layout.AlignCenter)).
			ScrollViewDirection(Const(layout.DirectionHorizontal)).Debug(),

		VStack(
			Text(BindOneWay(value, func(f float64) string {
				return strconv.FormatFloat(f, 'f', 0, 64)
			})).Debug().Center().Frame(Const(NewSize(Inf, 30))).Fill(Const(NewColor(128, 0, 0))),
			Slider(value, Bind(100.0), Bind(0.0)).Debug(),
			Circle().
				Fill(Const(NewColor(255))).
				Frame(Const(NewSize(100))).
				Shadow().
				Padding(),
			Toggle(Bind(false)).Debug(),
		).Debug(),
	).Spacing().Frame(Bind(NewSize(1000, 500)))
}
