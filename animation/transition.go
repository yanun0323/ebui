package animation

import (
	"time"
)

// Transition 定義一個視圖過渡動畫
type Transition interface {
	// Apply 應用過渡動畫到視圖，返回修改後的值
	Apply(progress float64, isAppearing bool) TransformModifier
}

// TransformModifier 表示一個變換修飾器，用於修改視圖的外觀
type TransformModifier struct {
	Opacity     float64
	Scale       float64
	Translation [2]float64 // x, y
	Rotation    float64    // 旋度，弧度制
}

// Identity 返回一個不進行任何變換的修飾器
func Identity() TransformModifier {
	return TransformModifier{
		Opacity:     1.0,
		Scale:       1.0,
		Translation: [2]float64{0, 0},
		Rotation:    0,
	}
}

// 預定義的過渡動畫

// FadeTransition 淡入淡出過渡
type FadeTransition struct {
	Style Style
}

func (t FadeTransition) Apply(progress float64, isAppearing bool) TransformModifier {
	mod := Identity()
	if isAppearing {
		mod.Opacity = t.Style.Value(progress)
	} else {
		mod.Opacity = 1 - t.Style.Value(progress)
	}
	return mod
}

// ScaleTransition 縮放過渡
type ScaleTransition struct {
	Style Style
	Scale float64 // 目標縮放比例
}

func (t ScaleTransition) Apply(progress float64, isAppearing bool) TransformModifier {
	mod := Identity()
	value := t.Style.Value(progress)

	if isAppearing {
		// 從小到大
		mod.Scale = 1 - (1-t.Scale)*(1-value)
	} else {
		// 從大到小
		mod.Scale = 1 - (1-t.Scale)*value
	}
	return mod
}

// SlideTransition 滑動過渡
type SlideTransition struct {
	Style Style
	Edge  Edge // 從哪個邊緣滑入/滑出
}

// Edge 表示邊緣方向
type Edge int

const (
	EdgeLeading Edge = iota
	EdgeTrailing
	EdgeTop
	EdgeBottom
)

func (t SlideTransition) Apply(progress float64, isAppearing bool) TransformModifier {
	mod := Identity()
	value := t.Style.Value(progress)
	offset := 0.0

	if !isAppearing {
		value = 1 - value
	}

	switch t.Edge {
	case EdgeLeading:
		offset = -1.0 + value
		mod.Translation = [2]float64{offset, 0}
	case EdgeTrailing:
		offset = 1.0 - value
		mod.Translation = [2]float64{offset, 0}
	case EdgeTop:
		offset = -1.0 + value
		mod.Translation = [2]float64{0, offset}
	case EdgeBottom:
		offset = 1.0 - value
		mod.Translation = [2]float64{0, offset}
	}

	return mod
}

// 預設轉場效果
func Fade() Transition {
	return FadeTransition{Style: EaseInOut(time.Millisecond * 300)}
}

func Scale(scale float64) Transition {
	return ScaleTransition{
		Style: EaseInOut(time.Millisecond * 300),
		Scale: scale,
	}
}

func Slide(edge Edge) Transition {
	return SlideTransition{
		Style: EaseInOut(time.Millisecond * 300),
		Edge:  edge,
	}
}

// 組合過渡效果
type CombinedTransition struct {
	Transitions []Transition
}

func (t CombinedTransition) Apply(progress float64, isAppearing bool) TransformModifier {
	if len(t.Transitions) == 0 {
		return Identity()
	}

	result := t.Transitions[0].Apply(progress, isAppearing)

	for i := 1; i < len(t.Transitions); i++ {
		mod := t.Transitions[i].Apply(progress, isAppearing)
		result.Opacity *= mod.Opacity
		result.Scale *= mod.Scale
		result.Translation[0] += mod.Translation[0]
		result.Translation[1] += mod.Translation[1]
		result.Rotation += mod.Rotation
	}

	return result
}

// Combined 組合多個過渡效果
func Combined(transitions ...Transition) Transition {
	return CombinedTransition{Transitions: transitions}
}
