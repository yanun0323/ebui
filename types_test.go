package ebui

import "testing"

func TestConstraintsConstrain(t *testing.T) {
	size := Flexible(S(100, 50)).Constrain(S(120, 20))
	if size.Width != 100 || size.Height != 20 {
		t.Fatalf("unexpected constrained size: %#v", size)
	}
}
