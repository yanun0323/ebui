package ebiten

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"

	ebitenlib "github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	ui "github.com/yanun0323/ebui"
	"github.com/yanun0323/ebui/internal/runtime"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

var defaultFontCandidates = []string{
	"resource/Taipei Sans TC Regular.ttf",
	"resource/NotoSansTC.ttf",
}

type Renderer struct {
	fontBytes  []byte
	faceCache  map[faceKey]font.Face
	imageCache map[string]*ebitenlib.Image
}

type faceKey struct {
	size   float64
	weight ui.FontWeight
}

func NewRenderer(preferredFontPath string) (*Renderer, error) {
	fontBytes, err := loadFontBytes(preferredFontPath)
	if err != nil {
		return nil, err
	}
	return &Renderer{
		fontBytes:  fontBytes,
		faceCache:  make(map[faceKey]font.Face),
		imageCache: make(map[string]*ebitenlib.Image),
	}, nil
}

func loadFontBytes(preferred string) ([]byte, error) {
	candidates := make([]string, 0, len(defaultFontCandidates)+1)
	if preferred != "" {
		candidates = append(candidates, preferred)
	}
	candidates = append(candidates, defaultFontCandidates...)
	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate)
		if err == nil {
			return data, nil
		}
	}
	return goregular.TTF, nil
}

func (r *Renderer) MeasureText(content string, style ui.TextStyle, lineLimit int) ui.Size {
	lines := splitLines(content, lineLimit)
	face := r.fontFace(style)
	maxWidth := 0
	lineHeight := face.Metrics().Height.Ceil()
	for _, line := range lines {
		bounds := ebitentext.BoundString(face, line)
		if bounds.Dx() > maxWidth {
			maxWidth = bounds.Dx()
		}
	}
	if len(lines) == 0 {
		return ui.S(0, float64(lineHeight))
	}
	height := float64(lineHeight*len(lines)) + math.Max(0, float64(len(lines)-1))*style.LineHeight
	return ui.S(float64(maxWidth), height)
}

func (r *Renderer) ImageSize(source ui.ImageSource) ui.Size {
	if source.Path == "" {
		return ui.S(0, 0)
	}
	img := r.loadImage(source.Path)
	if img == nil {
		return ui.S(0, 0)
	}
	size := img.Bounds().Size()
	return ui.S(float64(size.X), float64(size.Y))
}

func (r *Renderer) Draw(screen *ebitenlib.Image, frame runtime.Frame) {
	for _, op := range frame.Ops {
		if op.HasClip && !op.Rect.Intersects(op.Clip) {
			continue
		}
		switch op.Kind {
		case runtime.DrawKindRect:
			vector.DrawFilledRect(screen, float32(op.Rect.Min.X), float32(op.Rect.Min.Y), float32(op.Rect.Width()), float32(op.Rect.Height()), op.Fill.NRGBA(), false)
			if op.Border != nil {
				r.drawBorder(screen, op.Rect, *op.Border)
			}
		case runtime.DrawKindCircle:
			radius := float32(math.Min(op.Rect.Width(), op.Rect.Height()) / 2)
			vector.DrawFilledCircle(screen, float32(op.Rect.Min.X)+radius, float32(op.Rect.Min.Y)+radius, radius, op.Fill.NRGBA(), false)
		case runtime.DrawKindText:
			r.drawText(screen, op)
		case runtime.DrawKindImage:
			r.drawImage(screen, op)
		}
	}
}

func (r *Renderer) drawBorder(screen *ebitenlib.Image, rect ui.Rect, stroke ui.Stroke) {
	top := rect.Min.Y
	left := rect.Min.X
	right := rect.Max.X
	bottom := rect.Max.Y
	color := stroke.Color.NRGBA()
	if stroke.Insets.Top > 0 {
		vector.DrawFilledRect(screen, float32(left), float32(top), float32(rect.Width()), float32(stroke.Insets.Top), color, false)
	}
	if stroke.Insets.Bottom > 0 {
		vector.DrawFilledRect(screen, float32(left), float32(bottom-stroke.Insets.Bottom), float32(rect.Width()), float32(stroke.Insets.Bottom), color, false)
	}
	if stroke.Insets.Left > 0 {
		vector.DrawFilledRect(screen, float32(left), float32(top), float32(stroke.Insets.Left), float32(rect.Height()), color, false)
	}
	if stroke.Insets.Right > 0 {
		vector.DrawFilledRect(screen, float32(right-stroke.Insets.Right), float32(top), float32(stroke.Insets.Right), float32(rect.Height()), color, false)
	}
}

func (r *Renderer) drawText(screen *ebitenlib.Image, op runtime.DrawOp) {
	face := r.fontFace(op.Font)
	lines := splitLines(op.Text, 0)
	lineHeight := face.Metrics().Height.Ceil()
	y := int(op.Rect.Min.Y) + face.Metrics().Ascent.Ceil()
	for _, line := range lines {
		ebitentext.Draw(screen, line, face, int(op.Rect.Min.X), y, op.Foreground.NRGBA())
		y += lineHeight + int(op.Font.LineHeight)
	}
}

func (r *Renderer) drawImage(screen *ebitenlib.Image, op runtime.DrawOp) {
	img := r.loadImage(op.Image.Path)
	if img == nil {
		return
	}
	bounds := img.Bounds()
	options := &ebitenlib.DrawImageOptions{}
	options.GeoM.Scale(op.Rect.Width()/float64(bounds.Dx()), op.Rect.Height()/float64(bounds.Dy()))
	options.GeoM.Translate(op.Rect.Min.X, op.Rect.Min.Y)
	options.ColorScale.ScaleAlpha(float32(op.Opacity))
	screen.DrawImage(img, options)
}

func (r *Renderer) fontFace(style ui.TextStyle) font.Face {
	if style.Size <= 0 {
		style.Size = 16
	}
	key := faceKey{size: style.Size, weight: style.Weight}
	if face, ok := r.faceCache[key]; ok {
		return face
	}
	tt, err := opentype.Parse(r.fontBytes)
	if err != nil {
		return basicfont.Face7x13
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    style.Size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return basicfont.Face7x13
	}
	r.faceCache[key] = face
	return face
}

func (r *Renderer) loadImage(path string) *ebitenlib.Image {
	if path == "" {
		return nil
	}
	if image, ok := r.imageCache[path]; ok {
		return image
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	img := ebitenlib.NewImageFromImage(decoded)
	r.imageCache[path] = img
	return img
}

func splitLines(content string, lineLimit int) []string {
	if content == "" {
		return []string{""}
	}
	lines := strings.Split(content, "\n")
	if lineLimit > 0 && len(lines) > lineLimit {
		return lines[:lineLimit]
	}
	return lines
}
