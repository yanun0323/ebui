package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_Rectangle() View {
	return Rectangle().
		Fill(Const(NewColor(255, 255, 255))).
		Frame(Bind(NewSize(100, 100))).
		Offset(Const(NewPoint(50, 50))).
		Debug()
}
