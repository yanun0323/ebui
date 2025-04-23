package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_Text() View {
	return Text(Const("Hello, World!")).Offset(Const(NewPoint(100, 100))).Debug()
}
