package main

import (
	"github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/example/component"
)

func main() {
	contentView := component.ContentView("title", "content")

	ebui.Run("Windows Title", contentView)
}
