package ebui

import (
	"sync"
	"time"
)

type animationID uintptr

// animationExecutor 表示一個動畫執行器
type animationExecutor struct {
	onUpdate func(now int64) bool // 每幀更新時調用，返回值表示是否提前完成
	onCancel func()               // 取消動畫時調用
}

// animationManager 管理所有動畫執行器
type animationManager struct {
	mu        sync.RWMutex
	executors map[animationID]animationExecutor
}

var globalAnimationManager = &animationManager{
	executors: make(map[animationID]animationExecutor),
}

// AddExecutor 添加一個動畫執行器
func (am *animationManager) AddExecutor(id animationID, executor animationExecutor) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.executors[id] = executor
}

// Update 更新所有動畫執行器
func (am *animationManager) Update() {
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now().UnixMilli()
	// 執行所有執行器並標記完成的
	completedIndices := make([]animationID, 0)
	for id, executor := range am.executors {
		// 執行動畫
		completed := executor.onUpdate(now)
		if completed {
			completedIndices = append(completedIndices, id)
		}
	}

	for _, id := range completedIndices {
		executor, _ := am.removeExecutor(id)
		executor.onCancel()
	}

	// 如果有動畫正在執行，標記為需要重繪
	if len(am.executors) > 0 {
		globalStateManager.markDirty()
	}
}

// RemoveExecutor 移除一個動畫執行器
func (am *animationManager) RemoveExecutor(id animationID) (animationExecutor, bool) {
	am.mu.Lock()
	defer am.mu.Unlock()

	return am.removeExecutor(id)
}

func (am *animationManager) removeExecutor(id animationID) (animationExecutor, bool) {
	executor, ok := am.executors[id]
	if !ok {
		return animationExecutor{}, false
	}

	delete(am.executors, id)
	return executor, true
}
