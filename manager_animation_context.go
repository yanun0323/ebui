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

// WithAnimation uses specific animation style to execute the operation
// similar to SwiftUI's withAnimation function
func WithAnimation(body func(), style ...animation.Style) {
	s := animation.EaseInOut()
	if len(style) != 0 {
		s = style[0]
	}

	popStyle := globalContext.PushStyle(s)
	defer popStyle()
	body()
}

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

func newAnimationContext() *animationContext {
	return &animationContext{
		styleStack: []animation.Style{},
		animations: make(map[string]*animations),
	}
}

// CurrentStyle 返回當前活動的動畫風格
func (ctx *animationContext) CurrentStyle() animation.Style {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if len(ctx.styleStack) == 0 {
		return nil
	}
	return ctx.styleStack[len(ctx.styleStack)-1]
}

// PushStyle 將動畫風格添加到堆疊
func (ctx *animationContext) PushStyle(style animation.Style) (popStyle func() animation.Style) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.styleStack = append(ctx.styleStack, style)

	return func() animation.Style {
		ctx.mu.Lock()
		defer ctx.mu.Unlock()

		if len(ctx.styleStack) == 0 {
			return animation.None()
		}

		style := ctx.styleStack[len(ctx.styleStack)-1]
		ctx.styleStack = ctx.styleStack[:len(ctx.styleStack)-1]
		return style
	}
}
