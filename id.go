package ebui

import (
	"sync/atomic"
)

var (
	idPool atomic.Int64
)

func getID() int64 {
	return idPool.Add(1)
}
