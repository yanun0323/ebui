package ebui

type stateManager struct {
	dirty  bool
	bounds Rect
}

var globalStateManager = &stateManager{
	bounds: CGRect(0, 0, 0, 0),
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

func (sm *stateManager) SetBounds(bounds Rect) {
	if sm.bounds != bounds {
		sm.bounds = bounds
		sm.markDirty()
	}
}
