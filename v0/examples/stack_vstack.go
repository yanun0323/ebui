package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_VStack() View {
	return VStack(
		Spacer(),
		Text("Hello, World!"),
		Text("Hello, World!"),
		Spacer(),
	)
}
