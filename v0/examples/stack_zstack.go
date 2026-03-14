package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_ZStack() View {
	return ZStack(
		Spacer(),
		Rectangle().Fill(Const(red)),
		Text("Hello, World"),
		Text("_____"),
		Spacer(),
	).Frame(Const(NewSize(50))).Debug()
}
