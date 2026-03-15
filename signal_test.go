package ebui

import "testing"

func TestStateAndComputed(t *testing.T) {
	count := NewState(1)
	double := Map[int, int](count, func(v int) int { return v * 2 })

	if got := double.Value(); got != 2 {
		t.Fatalf("unexpected initial value: %d", got)
	}

	count.Set(3)
	if got := double.Value(); got != 6 {
		t.Fatalf("unexpected computed value: %d", got)
	}
}
