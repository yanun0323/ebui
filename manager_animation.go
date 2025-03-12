package ebui

import (
	"slices"
	"sync"
	"time"

	"github.com/yanun0323/ebui/animation"
)

// animationExecutor 表示一個動畫執行器
type animationExecutor struct {
	execute   func() bool     // 執行函數，返回是否完成
	style     animation.Style // 動畫風格
	startTime time.Time       // 開始時間
}

// animationManager 管理所有動畫執行器
type animationManager struct {
	mu        sync.RWMutex
	executors []animationExecutor
	nextID    int
}

var globalAnimationManager = &animationManager{}

// AddExecutor 添加一個動畫執行器
func (am *animationManager) AddExecutor(executor animationExecutor) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.nextID++
	am.executors = append(am.executors, executor)
}

// Update 更新所有動畫執行器
func (am *animationManager) Update() {
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now()

	// 執行所有執行器並標記完成的
	completedIndices := make([]int, 0)
	for i, executor := range am.executors {
		// 檢查是否動畫已完成
		elapsed := now.Sub(executor.startTime)
		duration := executor.style.Duration()

		// 執行動畫
		completed := executor.execute()

		// 如果已完成或超過時長，標記為完成
		if completed || elapsed >= duration {
			completedIndices = append(completedIndices, i)
		}
	}

	// 從後向前移除已完成的執行器
	for i := len(completedIndices) - 1; i >= 0; i-- {
		idx := completedIndices[i]
		am.executors = slices.Delete(am.executors, idx, idx+1)
	}

	// 如果有動畫正在執行，標記為需要重繪
	if len(am.executors) > 0 {
		globalStateManager.markDirty()
	}
}

// CreateAnimatedExecutor 創建一個基於動畫的執行器
// 它會在每一幀調用 onUpdate 直到動畫完成
func (am *animationManager) CreateAnimatedExecutor(
	style animation.Style,
	onUpdate func(progress float64) bool, // 每幀更新時調用，返回值表示是否提前完成
) time.Time {
	startTime := time.Now()
	duration := style.Duration()

	executor := animationExecutor{
		style:     style,
		startTime: startTime,
		execute: func() bool {
			elapsed := time.Since(startTime)
			if elapsed >= duration {
				// 動畫完成，使用最終進度值
				return onUpdate(1.0)
			}

			// 計算當前進度
			progress := float64(elapsed) / float64(duration)
			curveProgress := style.Value(progress)

			// 調用更新函數，它可能返回提前完成
			return onUpdate(curveProgress)
		},
	}

	am.AddExecutor(executor)
	return startTime
}
