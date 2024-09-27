package ebui

import "image/color"

var _zeroUIViewLayout = uiViewLayout{frame: _zeroSize}

type uiViewLayout struct {
	start      point       // 計算的開始位置
	offset     point       // modifier 設定的 offset
	frame      size        // modifier 設定的 frame
	padding    bounds      // modifier 設定的 padding
	margin     bounds      // 計算的 margin
	background color.Color // modifier 設定的 background
}
