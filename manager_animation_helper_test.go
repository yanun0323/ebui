package ebui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimateValue(t *testing.T) {
	{
		res := animateValue(0, 100, 0.25)
		assert.Equal(t, 25, res)
	}

	{
		res := animateValue(0.0, 100.0, 0.75)
		assert.Equal(t, 75.0, res)
	}

	{
		res := animateValue(NewPoint(0, 100), NewPoint(100, 0), 0.25)
		assert.Equal(t, NewPoint(25, 75), res)
	}

	{
		res := animateValue(NewRect(0, 100, 100, 200), NewRect(100, 0, 200, 100), 0.25)
		assert.Equal(t, NewRect(25, 75, 125, 175), res)
	}

	{
		res := animateValue(NewSize(0, 100), NewSize(100, 0), 0.25)
		assert.Equal(t, NewSize(25, 75), res)
	}

	{
		res := animateValue(NewInset(0, 100, 100, 0), NewInset(100, 0, 0, 100), 0.25)
		assert.Equal(t, NewInset(25, 75, 75, 25), res)
	}

	{
		res := animateValue(NewColor(0, 0, 0, 0), NewColor(100, 100, 100, 100), 0.25)
		assert.Equal(t, NewColor(25, 25, 25, 25), res)
	}
}
