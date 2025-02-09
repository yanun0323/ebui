package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type GridItem struct {
	View    SomeView
	Column  int
	Row     int
	ColSpan int
	RowSpan int
}

type GridLayout struct {
	*viewContext

	columns   int
	rows      int
	columnGap float64
	rowGap    float64
	items     []GridItem
	frame     image.Rectangle
	cellSize  image.Point
	cache     *ViewCache
}

func Grid(columns, rows int, items ...GridItem) SomeView {
	grid := &GridLayout{
		columns:   columns,
		rows:      rows,
		items:     items,
		columnGap: 8,
		rowGap:    8,
		cache:     NewViewCache(),
	}
	grid.viewContext = NewViewContext(grid)

	return grid
}

func (g *GridLayout) layout(bounds image.Rectangle) image.Rectangle {
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
		item.View.layout(itemBounds)
	}

	return bounds
}

func (g *GridLayout) draw(screen *ebiten.Image) {
	for _, item := range g.items {
		item.View.draw(screen)
	}
}

// 新增的方法，用於設置列間距
func (g *GridLayout) WithColumnGap(gap float64) *GridLayout {
	g.columnGap = gap
	return g
}

// 新增的方法，用於設置行間距
func (g *GridLayout) WithRowGap(gap float64) *GridLayout {
	g.rowGap = gap
	return g
}
