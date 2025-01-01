package ebui

import (
	"testing"

	"go.uber.org/goleak"
)

func TestEbui(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Log("ebui test")
}
