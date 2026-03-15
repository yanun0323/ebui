package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_HStack() View {
	return HStack(
		Spacer(),
		Text("Hello, World!"),
		Text("Hello, World!"),
		Spacer(),
	)
}
