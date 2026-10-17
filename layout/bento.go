package layout

import (
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BentoTile spans Cols and Rows; zero spans one.
type BentoTile struct {
	Cols, Rows int
	Content    func()
}

// Bento is a grid of spanning tiles; tiles lift.
func Bento(c *ui.Context, columns int, tiles []BentoTile) ui.Element {
	if columns < 1 {
		panic("layout: Bento needs a column")
	}
	t := c.Theme()
	return ui.Grid(c).FillWidth().Columns(columns).Gap(t.Space(3)).Children(func() {
		for i, tile := range tiles {
			cols, rows := max(tile.Cols, 1), max(tile.Rows, 1)
			if tile.Cols < 0 || tile.Rows < 0 || cols > columns {
				panic("layout: Bento tile span out of range")
			}
			kit.Lift(kit.Card(c).Key(i).Padding(t.Space(4)).Gap(t.Space(2))).ColumnSpan(cols).RowSpan(rows).
				Children(tile.Content)
		}
	})
}
