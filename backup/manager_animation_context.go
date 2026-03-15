package ebui

import (
	"sync"
	"time"

	"github.com/yanun0323/ebui/animation"
)

var (
	// global animation context, manages the current active animations
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

// animationContext manages the animation context and state
type animationContext struct {
	mu         sync.RWMutex
	styleStack []animation.Style
	animations map[string]*animations
}

// animations represents an ongoing animation
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

// CurrentStyle returns the current active animation style
func (ctx *animationContext) CurrentStyle() animation.Style {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if len(ctx.styleStack) == 0 {
		return nil
	}
	return ctx.styleStack[len(ctx.styleStack)-1]
}

// PushStyle adds an animation style to the stack
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
