package view

import (
	"strconv"

	. "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/layout"
)

func PageScrollView() View {
	return &pageScrollView{}
}

type pageScrollView struct {
}

func (v *pageScrollView) Body() SomeView {
	rect := func(i int) SomeView {
		return Text(Const(strconv.Itoa(i))).
			Align(Bind(layout.AlignCenter)).
			Frame(Bind(NewSize(100))).
			BackgroundColor(Bind(NewColor(128, 0, 0)))
	}

	enum := func(count int) []View {
		res := make([]View, 0, count)
		for i := range count {
			res = append(res, rect(i))
		}
		return res
	}

	return HStack(
		ScrollView(
			VStack(
				enum(10)...,
			).Spacing().Align(Bind(layout.AlignTrailing)),
		).ScrollViewDirection(Const(layout.DirectionVertical)),
		ScrollView(
			HStack(
				enum(10)...,
			).Spacing().Align(Bind(layout.AlignCenter)),
		).ScrollViewDirection(Const(layout.DirectionHorizontal)),
	).Spacing()
}
