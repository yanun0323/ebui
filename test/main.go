package main

import (
	"github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/test/component"
)

func main() {
	contentView := component.TestView("title", "content\ncontent")
	println("Hi")
	ebui.Run("Windows Title", contentView, false)
}
