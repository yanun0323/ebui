package ebui

import "image"

type StateManager struct {
	dirty bool
	root  View
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

func (sm *StateManager) Update() error {
	if sm.dirty {
		sm.root.Layout(sm.bounds)
		sm.dirty = false
	}
	return nil
}

func (sm *StateManager) isDirty() bool {
	return sm.dirty
}
