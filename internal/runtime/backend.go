package runtime

import ui "github.com/yanun0323/ebui"

type Measurer interface {
	MeasureText(content string, style ui.TextStyle, lineLimit int) ui.Size
	ImageSize(source ui.ImageSource) ui.Size
}

type DrawKind uint8

const (
	DrawKindRect DrawKind = iota
	DrawKindCircle
	DrawKindText
	DrawKindImage
)

type DrawOp struct {
	Kind       DrawKind
	Rect       ui.Rect
	Fill       ui.Color
	Border     *ui.Stroke
	Text       string
	Font       ui.TextStyle
	Foreground ui.Color
	Image      ui.ImageSource
	Opacity    float64
	Clip       ui.Rect
	HasClip    bool
}

type Frame struct {
	Ops []DrawOp
}
