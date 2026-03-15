package ebui

import "sync/atomic"

type stateManager struct {
	dirty  atomic.Bool
	bounds CGRect
}

var globalStateManager = &stateManager{
	bounds: NewRect(0, 0, 0, 0),
}

func (sm *stateManager) markDirty() {
	sm.dirty.Store(true)
}

func (sm *stateManager) clearDirty() {
	sm.dirty.Store(false)
}

func (sm *stateManager) isDirty() bool {
	return sm.dirty.Load()
}

func (sm *stateManager) GetBounds() CGRect {
	return sm.bounds
}

func (sm *stateManager) SetBounds(bounds CGRect) {
	if sm.bounds != bounds {
		sm.bounds = bounds
		sm.markDirty()
	}
}
