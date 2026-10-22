package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Column is a table column; zero Width shares.
type Column struct {
	Title string
	Width float32
	// End aligns cells to the end.
	End bool
}

// DataTableOptions tune a DataTable.
type DataTableOptions struct {
	// Empty shows while there are no rows.
	Empty string
	// Height scrolls the rows within it; zero shows them all.
	Height float32
}

// DataTable sorts on header clicks; caller orders rows.
func DataTable(c *ui.Context, label string, cols []Column, n int, sort *ui.SortOrder, cell func(row, col int), o DataTableOptions) ui.Element {
	columns("DataTable", label, cols)
	if n < 0 || cell == nil {
		panic(fmt.Sprintf("data: DataTable %q needs a row count and cells", label))
	}
	def(&o.Empty, "No rows")
	t := c.Theme()
	root := frame(c, label)
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	root.Children(func() {
		head(c, root, "DataTable", cols, sort, nil)
		body(c, root, list, n, o.Height, func(i int) {
			line(c, hovered(c, false), cols, nil, func(col int) { cell(i, col) })
		})
		if n == 0 {
			kit.Muted(c, o.Empty).Padding(t.Space(3))
		}
	})
	return root
}
