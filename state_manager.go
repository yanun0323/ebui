package ebui

type StateManager struct {
	dirty  bool
	bounds CGRect
}

var globalStateManager = &StateManager{
	bounds: rect(0, 0, 0, 0),
}

func (sm *StateManager) markDirty() {
	sm.dirty = true
}

func (sm *StateManager) clearDirty() {
	sm.dirty = false
}

func (sm *StateManager) isDirty() bool {
	return sm.dirty
}

func (sm *StateManager) SetBounds(bounds CGRect) {
	if sm.bounds != bounds {
		sm.bounds = bounds
		sm.markDirty()
	}
}
