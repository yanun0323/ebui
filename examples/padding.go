package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_Padding() View {
	return Rectangle().
		Frame(Const(NewSize(100))).
		Fill(AccentColor).
		Border(Const(NewInset(2)), Const(green)).
		Padding(Const(NewInset(20))).
		Border(Const(NewInset(2)), Const(blue)).
		Offset(Const(NewPoint(20, 20))).
		Debug()
}
