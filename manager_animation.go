package ebui

import (
	"slices"
	"time"
)

type animationExecutor struct {
	execute   func()
	expiredAt time.Time
}

type animationManager struct {
	executors []animationExecutor
}

var globalAnimationManager = &animationManager{}

func (am *animationManager) AddExecutor(executor animationExecutor) {
	am.executors = append(am.executors, executor)
}

func (am *animationManager) Update() {
	now := time.Now()
	for _, executor := range am.executors {
		if executor.expiredAt.Before(now) {
			executor.execute()
			globalStateManager.markDirty()
		}
	}

	am.executors = slices.DeleteFunc(am.executors, func(executor animationExecutor) bool {
		return executor.expiredAt.Before(now)
	})
}
