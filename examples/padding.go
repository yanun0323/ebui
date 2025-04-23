package examples

import (
	. "github.com/yanun0323/ebui"
)

func Preview_Padding() View {
	return Rectangle().
		Fill(AccentColor).
		Frame(Const(NewSize(100))).
		Debug().
		Padding()
}
