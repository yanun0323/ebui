package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_ZStack() View {
	return ZStack(
		Spacer(),
		Text("Hello, World"),
		Text("_____"),
		Spacer(),
	)
}
