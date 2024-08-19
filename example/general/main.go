package main

import (
	"time"

	"github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/example/general/component"
	"github.com/yanun0323/pkg/logs"
)

func main() {

	t := "Hello"

	logs.SetDefault(logs.New(logs.LevelInfo))
	// logs.SetDefault(logs.New(logs.LevelDebug))

	go func() {
		for {
			time.Sleep(time.Second)
			t += "."
		}
	}()

	ebui.Run("Windows Title", component.ContentView("content text"), true)
}
