package ebui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHelperCorner(t *testing.T) {
	suite.Run(t, new(HelperCornerSuite))
}

type HelperCornerSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *HelperCornerSuite) SetupSuite() {
	su.ctx = context.Background()
}

func (su *HelperCornerSuite) TestCornerHandler() {
	w, h := 100, 100
	r := 10
	b := 10
	p := r + b
	hl := newCornerHandler(w, h, float64(r), float64(b))
	hl.Execute(func(outside, between bool, x, y int) {
		if x > p && x < w-p && y > p && y < h-p {
			su.Require().False(outside || between,
				"should be inside, x: %d, y: %d, outside: %t, between: %t",
				x, y, outside, between,
			)
		}
	})
}
