package animation

import (
	"sync"
	"time"
)

var (
	// 全局動畫上下文，管理當前活動的動畫
	globalContext = newAnimationContext()
)

// AnimationContext 管理動畫上下文和狀態
type AnimationContext struct {
	mu         sync.RWMutex
	styleStack []Style
	animations map[string]*Animation
}

// Animation 表示一個正在進行的動畫
type Animation struct {
	ID         string
	Style      Style
	StartTime  time.Time
	Completed  bool
	OnComplete func()
}

// newAnimationContext 創建新的動畫上下文
func newAnimationContext() *AnimationContext {
	return &AnimationContext{
		styleStack: []Style{},
		animations: make(map[string]*Animation),
	}
}

// CurrentStyle 返回當前活動的動畫風格
func (ctx *AnimationContext) CurrentStyle() Style {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if len(ctx.styleStack) == 0 {
		return None()
	}
	return ctx.styleStack[len(ctx.styleStack)-1]
}

// PushStyle 將動畫風格添加到堆疊
func (ctx *AnimationContext) PushStyle(style Style) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.styleStack = append(ctx.styleStack, style)
}

// PopStyle 從堆疊中彈出動畫風格
func (ctx *AnimationContext) PopStyle() Style {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if len(ctx.styleStack) == 0 {
		return None()
	}

	style := ctx.styleStack[len(ctx.styleStack)-1]
	ctx.styleStack = ctx.styleStack[:len(ctx.styleStack)-1]
	return style
}

// RegisterAnimation 註冊一個正在進行的動畫
func (ctx *AnimationContext) RegisterAnimation(id string, style Style, onComplete func()) *Animation {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	anim := &Animation{
		ID:         id,
		Style:      style,
		StartTime:  time.Now(),
		Completed:  false,
		OnComplete: onComplete,
	}

	ctx.animations[id] = anim
	return anim
}

// RemoveAnimation 移除一個動畫
func (ctx *AnimationContext) RemoveAnimation(id string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	delete(ctx.animations, id)
}

// UpdateAnimations 更新所有動畫的狀態
func (ctx *AnimationContext) UpdateAnimations() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	now := time.Now()
	for id, anim := range ctx.animations {
		if anim.Completed {
			continue
		}

		duration := anim.Style.Duration()
		elapsed := now.Sub(anim.StartTime)

		if elapsed >= duration {
			anim.Completed = true
			if anim.OnComplete != nil {
				go anim.OnComplete()
			}
			delete(ctx.animations, id)
		}
	}
}

// WithAnimation 使用特定動畫風格執行操作
// 類似於 SwiftUI 的 withAnimation 函數
func WithAnimation(style Style, body func()) {
	globalContext.PushStyle(style)
	defer globalContext.PopStyle()
	body()
}

// GetCurrentStyle 獲取當前的動畫風格
func GetCurrentStyle() Style {
	return globalContext.CurrentStyle()
}
