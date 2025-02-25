package ebui

type stateManager struct {
	dirty  bool
	bounds CGRect
}

var globalStateManager = &stateManager{
	bounds: NewRect(0, 0, 0, 0),
}

func (sm *stateManager) markDirty() {
	sm.dirty = true
}

func (sm *stateManager) clearDirty() {
	sm.dirty = false
}

func (sm *stateManager) isDirty() bool {
	return sm.dirty
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
