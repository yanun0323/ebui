package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type GridItem struct {
	View    View
	Column  int
	Row     int
	ColSpan int
	RowSpan int
}

type GridLayout struct {
	columns   int
	rows      int
	columnGap float64
	rowGap    float64
	items     []GridItem
	frame     image.Rectangle
	cellSize  image.Point
	cache     *ViewCache
}

func Grid(columns, rows int, items ...GridItem) ViewBuilder {
	return ViewBuilder{
		build: func() View {
			grid := &GridLayout{
				columns:   columns,
				rows:      rows,
				items:     items,
				columnGap: 8,
				rowGap:    8,
				cache:     NewViewCache(),
			}

			// 構建子視圖
			for i := range grid.items {
				grid.items[i].View = grid.items[i].View.Build()
			}

			return grid
		},
	}
}

func (g *GridLayout) Layout(bounds image.Rectangle) image.Rectangle {
	g.frame = bounds

	// 計算單元格大小
	availWidth := float64(bounds.Dx()) - g.columnGap*float64(g.columns-1)
	availHeight := float64(bounds.Dy()) - g.rowGap*float64(g.rows-1)

	g.cellSize = image.Point{
		X: int(availWidth / float64(g.columns)),
		Y: int(availHeight / float64(g.rows)),
	}

	// 為每個項目計算佈局
	for i := range g.items {
		item := &g.items[i]

		// 計算項目位置
		x := bounds.Min.X + item.Column*(g.cellSize.X+int(g.columnGap))
		y := bounds.Min.Y + item.Row*(g.cellSize.Y+int(g.rowGap))

		// 計算項目大小
		width := item.ColSpan*g.cellSize.X + int(float64(item.ColSpan-1)*g.columnGap)
		height := item.RowSpan*g.cellSize.Y + int(float64(item.RowSpan-1)*g.rowGap)

		// 設置項目邊界
		itemBounds := image.Rect(x, y, x+width, y+height)
		item.View.Layout(itemBounds)
	}

	return bounds
}

func (g *GridLayout) Draw(screen *ebiten.Image) {
	for _, item := range g.items {
		item.View.Draw(screen)
	}
}

// 實現 View 介面
func (g *GridLayout) Build() View {
	return g
}
