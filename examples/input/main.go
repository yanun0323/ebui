package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	regularFont font.Face
)

func init() {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	regularFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// TextInput 代表一個文字輸入框
type TextInput struct {
	X, Y, Width, Height int
	Text                string
	CursorPos           int
	SelectionStart      int
	SelectionEnd        int
	Active              bool
	Focused             bool
}

// NewTextInput 創建一個新的輸入框
func NewTextInput(x, y, width, height int) *TextInput {
	return &TextInput{
		X:              x,
		Y:              y,
		Width:          width,
		Height:         height,
		CursorPos:      0,
		SelectionStart: -1,
		SelectionEnd:   -1,
		Active:         true,
		Focused:        false,
	}
}

// Update 更新輸入框狀態
func (t *TextInput) Update() error {
	// 檢查是否點擊輸入框
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= t.X && x < t.X+t.Width && y >= t.Y && y < t.Y+t.Height {
			t.Focused = true

			// 根據點擊位置設置游標
			// textWidth := text.BoundString(regularFont, t.Text[:t.CursorPos]).Dx()
			clickPos := x - t.X

			// 簡單實現：根據點擊位置設置游標
			for i := 0; i <= len(t.Text); i++ {
				width := text.BoundString(regularFont, t.Text[:i]).Dx()
				if width >= clickPos {
					t.CursorPos = i
					t.SelectionStart = -1
					t.SelectionEnd = -1
					break
				}
			}
		} else {
			t.Focused = false
		}
	}

	// 如果輸入框被聚焦，處理鍵盤輸入
	if t.Focused {
		// 處理文字輸入
		runes := ebiten.InputChars()
		if len(runes) > 0 {
			if t.SelectionStart != -1 && t.SelectionEnd != -1 {
				// 刪除選中的文字並在原位置插入新文字
				t.deleteSelection()
			}

			// 插入字符
			if t.CursorPos == len(t.Text) {
				t.Text += string(runes)
			} else {
				t.Text = t.Text[:t.CursorPos] + string(runes) + t.Text[t.CursorPos:]
			}
			t.CursorPos += len(runes)
		}

		// 處理退格鍵
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if t.SelectionStart != -1 && t.SelectionEnd != -1 {
				t.deleteSelection()
			} else if t.CursorPos > 0 {
				t.Text = t.Text[:t.CursorPos-1] + t.Text[t.CursorPos:]
				t.CursorPos--
			}
		}

		// 處理刪除鍵
		if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
			if t.SelectionStart != -1 && t.SelectionEnd != -1 {
				t.deleteSelection()
			} else if t.CursorPos < len(t.Text) {
				t.Text = t.Text[:t.CursorPos] + t.Text[t.CursorPos+1:]
			}
		}

		// 游標左移
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				// 按下Shift鍵開始選取
				if t.SelectionStart == -1 {
					t.SelectionStart = t.CursorPos
					t.SelectionEnd = t.CursorPos
				}

				if t.CursorPos > 0 {
					t.CursorPos--
					t.SelectionEnd = t.CursorPos
				}
			} else {
				// 沒有按Shift鍵，取消選取
				if t.SelectionStart != -1 && t.SelectionEnd != -1 {
					t.CursorPos = t.SelectionStart
					t.SelectionStart = -1
					t.SelectionEnd = -1
				} else if t.CursorPos > 0 {
					t.CursorPos--
				}
			}
		}

		// 游標右移
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				// 按下Shift鍵開始選取
				if t.SelectionStart == -1 {
					t.SelectionStart = t.CursorPos
					t.SelectionEnd = t.CursorPos
				}

				if t.CursorPos < len(t.Text) {
					t.CursorPos++
					t.SelectionEnd = t.CursorPos
				}
			} else {
				// 沒有按Shift鍵，取消選取
				if t.SelectionStart != -1 && t.SelectionEnd != -1 {
					t.CursorPos = t.SelectionEnd
					t.SelectionStart = -1
					t.SelectionEnd = -1
				} else if t.CursorPos < len(t.Text) {
					t.CursorPos++
				}
			}
		}

		// 處理 Ctrl+A (全選)
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyA) {
			t.SelectionStart = 0
			t.SelectionEnd = len(t.Text)
			t.CursorPos = t.SelectionEnd
		}

		// 處理 Ctrl+C (複製)
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
			if t.SelectionStart != -1 && t.SelectionEnd != -1 {
				selectedText := ""
				start, end := t.getSelectionRange()
				selectedText = t.Text[start:end]

				// 在實際應用中，你需要使用系統剪貼簿API
				// 這裡僅作為示範，實際的功能需要額外的庫支援
				log.Printf("複製到剪貼簿: %s", selectedText)
			}
		}

		// 處理 Ctrl+V (貼上)
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
			// 在實際應用中，你需要使用系統剪貼簿API獲取文字
			// 這裡僅作為示範，使用一個假的剪貼簿內容
			clipboardText := "貼上的文字"

			if t.SelectionStart != -1 && t.SelectionEnd != -1 {
				t.deleteSelection()
			}

			if t.CursorPos == len(t.Text) {
				t.Text += clipboardText
			} else {
				t.Text = t.Text[:t.CursorPos] + clipboardText + t.Text[t.CursorPos:]
			}
			t.CursorPos += len(clipboardText)
		}

		// 處理 Ctrl+X (剪切)
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyX) {
			if t.SelectionStart != -1 && t.SelectionEnd != -1 {
				start, end := t.getSelectionRange()
				selectedText := t.Text[start:end]

				// 在實際應用中，你需要使用系統剪貼簿API
				log.Printf("剪切到剪貼簿: %s", selectedText)

				t.deleteSelection()
			}
		}
	}

	return nil
}

// 刪除選中的文字
func (t *TextInput) deleteSelection() {
	if t.SelectionStart != -1 && t.SelectionEnd != -1 {
		start, end := t.getSelectionRange()
		t.Text = t.Text[:start] + t.Text[end:]
		t.CursorPos = start
		t.SelectionStart = -1
		t.SelectionEnd = -1
	}
}

// 獲取選擇範圍（確保開始位置小於結束位置）
func (t *TextInput) getSelectionRange() (int, int) {
	start, end := t.SelectionStart, t.SelectionEnd
	if start > end {
		start, end = end, start
	}
	return start, end
}

// Draw 繪製輸入框
func (t *TextInput) Draw(screen *ebiten.Image) {
	// 繪製輸入框背景
	ebitenutil.DrawRect(screen, float64(t.X), float64(t.Y), float64(t.Width), float64(t.Height), color.RGBA{240, 240, 240, 255})

	// 繪製輸入框邊框
	borderColor := color.RGBA{200, 200, 200, 255}
	if t.Focused {
		borderColor = color.RGBA{0, 120, 215, 255}
	}
	ebitenutil.DrawRect(screen, float64(t.X), float64(t.Y), float64(t.Width), 1, borderColor)
	ebitenutil.DrawRect(screen, float64(t.X), float64(t.Y), 1, float64(t.Height), borderColor)
	ebitenutil.DrawRect(screen, float64(t.X), float64(t.Y+t.Height-1), float64(t.Width), 1, borderColor)
	ebitenutil.DrawRect(screen, float64(t.X+t.Width-1), float64(t.Y), 1, float64(t.Height), borderColor)

	// 繪製選中區域背景
	if t.SelectionStart != -1 && t.SelectionEnd != -1 {
		start, end := t.getSelectionRange()
		startX := t.X + text.BoundString(regularFont, t.Text[:start]).Dx()
		endX := t.X + text.BoundString(regularFont, t.Text[:end]).Dx()
		ebitenutil.DrawRect(screen, float64(startX), float64(t.Y+2), float64(endX-startX), float64(t.Height-4), color.RGBA{173, 214, 255, 255})
	}

	// 繪製文字
	text.Draw(screen, t.Text, regularFont, t.X+5, t.Y+t.Height/2+8, color.Black)

	// 繪製游標（如果輸入框處於活動狀態）
	if t.Focused && t.Active {
		// 讓游標閃爍
		if (int(ebiten.ActualTPS())/2)%2 == 0 {
			cursorX := t.X + text.BoundString(regularFont, t.Text[:t.CursorPos]).Dx() + 5
			ebitenutil.DrawLine(screen, float64(cursorX), float64(t.Y+5), float64(cursorX), float64(t.Y+t.Height-5), color.Black)
		}
	}
}

// Game 是主要的遊戲結構
type Game struct {
	textInput *TextInput
}

// NewGame 創建一個新的遊戲實例
func NewGame() *Game {
	return &Game{
		textInput: NewTextInput(100, 200, 400, 40),
	}
}

// Update 更新遊戲狀態
func (g *Game) Update() error {
	return g.textInput.Update()
}

// Draw 繪製遊戲畫面
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255})
	g.textInput.Draw(screen)

	ebitenutil.DebugPrint(screen, "點擊輸入框並輸入文字\n使用方向鍵+Shift選取文字\n使用Ctrl+A/C/V/X進行全選/複製/貼上/剪切")
}

// Layout 定義遊戲窗口大小
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten 文字輸入框示例")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
