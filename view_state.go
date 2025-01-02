package ebui

// import "image"

// type ViewID string

// type ViewState struct {
// 	id       ViewID
// 	dirty    bool
// 	frame    image.Rectangle
// 	children map[ViewID]*ViewState
// }

// func (vs *ViewState) MarkDirty() {
// 	vs.dirty = true
// 	// 向上傳播更新
// 	for _, child := range vs.children {
// 		child.MarkDirty()
// 	}
// }

// // 優化的狀態管理器
// type OptimizedStateManager struct {
// 	states map[ViewID]*ViewState
// 	root   *ViewState
// }

// func (sm *OptimizedStateManager) UpdateIfNeeded() {
// 	if sm.root.dirty {
// 		sm.performLayout(sm.root)
// 		sm.root.dirty = false
// 	}
// }
