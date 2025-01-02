package ebui

import "image"

type StateManager struct {
	dirty  bool
	bounds image.Rectangle
}

func NewStateManager(bounds image.Rectangle) *StateManager {
	return &StateManager{
		bounds: bounds,
	}
}

func (sm *StateManager) MarkDirty() {
	sm.dirty = true
}

func (sm *StateManager) clearDirty() {
	sm.dirty = false
}

func (sm *StateManager) isDirty() bool {
	return sm.dirty
}

func (sm *StateManager) SetBounds(bounds image.Rectangle) {
	if sm.bounds != bounds {
		sm.bounds = bounds
		sm.MarkDirty()
	}
}

var defaultStateManager = NewStateManager(image.Rect(0, 0, 0, 0))
