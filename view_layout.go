package ebui

import (
	"image/color"

	"github.com/yanun0323/pkg/sys"
)

var _zeroUIViewLayout = uiViewLayout{frame: _zeroSize}

type uiViewLayout struct {
	start      point       // 計算的開始位置
	offset     point       // modifier 設定的 offset
	frame      size        // modifier 設定的 frame
	padding    bounds      // modifier 設定的 padding
	margin     bounds      // modifier 設定的 外部 padding
	background color.Color // modifier 設定的 background
}

func (v *uiViewLayout) getDrawSize(cachedSize size) size {
	return size{
		w: sys.If(v.frame.w == -1, cachedSize.w, v.frame.w) - v.padding.left - v.padding.right,
		h: sys.If(v.frame.h == -1, cachedSize.h, v.frame.h) - v.padding.top - v.padding.bottom,
	}
}
