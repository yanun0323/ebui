package ebui

import (
	"sync"
	"time"

	"github.com/yanun0323/ebui/animation"
)

var (
	// 全局動畫上下文，管理當前活動的動畫
	globalContext = newAnimationContext()
)

// animationContext 管理動畫上下文和狀態
type animationContext struct {
	mu         sync.RWMutex
	styleStack []animation.Style
	animations map[string]*animations
}

// animations 表示一個正在進行的動畫
type animations struct {
	ID         string
	Style      animation.Style
	StartTime  time.Time
	Completed  bool
	OnComplete func()
}

// newAnimationContext 創建新的動畫上下文
func newAnimationContext() *animationContext {
	return &animationContext{
		styleStack: []animation.Style{},
		animations: make(map[string]*animations),
	}
}

// Currentanimation.Style 返回當前活動的動畫風格
func (ctx *animationContext) CurrentStyle() animation.Style {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if len(ctx.styleStack) == 0 {
		return animation.None()
	}
	return ctx.styleStack[len(ctx.styleStack)-1]
}

// Pushanimation.Style 將動畫風格添加到堆疊
func (ctx *animationContext) PushStyle(style animation.Style) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.styleStack = append(ctx.styleStack, style)
}

// Popanimation.Style 從堆疊中彈出動畫風格
func (ctx *animationContext) PopStyle() animation.Style {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if len(ctx.styleStack) == 0 {
		return animation.None()
	}

	style := ctx.styleStack[len(ctx.styleStack)-1]
	ctx.styleStack = ctx.styleStack[:len(ctx.styleStack)-1]
	return style
}

// RegisterAnimation 註冊一個正在進行的動畫
func (ctx *animationContext) RegisterAnimation(id string, style animation.Style, onComplete func()) *animations {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	anim := &animations{
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
func (ctx *animationContext) RemoveAnimation(id string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	delete(ctx.animations, id)
}

// UpdateAnimations 更新所有動畫的狀態
func (ctx *animationContext) UpdateAnimations() {
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
func WithAnimation(style animation.Style, body func()) {
	globalContext.PushStyle(style)
	defer globalContext.PopStyle()
	body()
}

// getCurrentStyle 獲取當前的動畫風格
func getCurrentStyle() animation.Style {
	return globalContext.CurrentStyle()
}
