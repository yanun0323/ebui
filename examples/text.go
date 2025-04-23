package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_Text() View {
	return Text(Const("Hello, World 2!")).Offset(Const(NewPoint(100, 100)))
}
