package ebui

import (
	"sync"
	"time"
)

type animationID uintptr

// animationExecutor represents an animation executor
type animationExecutor struct {
	onUpdate func(now int64) bool // called every frame, returns true if the animation is completed
	onCancel func()               // called when the animation is canceled
}

// animationManager manages all animation executors
type animationManager struct {
	mu        sync.RWMutex
	executors map[animationID]animationExecutor
}

var globalAnimationManager = &animationManager{
	executors: make(map[animationID]animationExecutor),
}

// AddExecutor adds an animation executor
func (am *animationManager) AddExecutor(id animationID, executor animationExecutor) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.executors[id] = executor
}

// Update updates all animation executors
func (am *animationManager) Update() {
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now().UnixMilli()
	// execute all executors and mark completed
	completedIndices := make([]animationID, 0, len(am.executors))
	for id, executor := range am.executors {
		// execute animation
		completed := executor.onUpdate(now)
		if completed {
			completedIndices = append(completedIndices, id)
		}
	}

	for _, id := range completedIndices {
		executor, ok := am.removeExecutor(id)
		if ok {
			executor.onCancel()
		}
	}

	// if there is an animation running, mark as dirty
	if len(am.executors) > 0 {
		globalStateManager.markDirty()
	}
}

// RemoveExecutor removes an animation executor
func (am *animationManager) RemoveExecutor(id animationID) (animationExecutor, bool) {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.removeExecutor(id)
}

func (am *animationManager) removeExecutor(id animationID) (animationExecutor, bool) {
	executor, ok := am.executors[id]
	if ok {
		delete(am.executors, id)
	}

	return executor, ok
}
