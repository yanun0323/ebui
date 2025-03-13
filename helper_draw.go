package ebui

import (
	"fmt"
)

func logf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func someViews(views ...View) []SomeView {
	someViews := make([]SomeView, 0, len(views))
	for _, view := range views {
		someViews = append(someViews, view.Body())
	}
	return someViews
}
