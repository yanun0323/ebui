package ebui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/yanun0323/ebui/animation"
)

func TestBinding(t *testing.T) {
	suite.Run(t, new(BindingSuite))
}

type BindingSuite struct {
	suite.Suite
}

func (su *BindingSuite) Test() {
	res := 100000.0
	b := Bind(0.0)
	b.Set(res, animation.Linear(10*time.Second))
	animRes := b.animResult.Load()
	su.NotNil(animRes)
	su.Equal(res, *animRes)
}
