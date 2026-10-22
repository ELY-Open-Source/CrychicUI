package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PinnedTable fixes the first column; the rest scroll sideways.
func PinnedTable(c *ui.Context, label string, cols []Column, n int, cell func(row, col int)) ui.Element {
	columns("PinnedTable", label, cols)
	if len(cols) < 2 || n < 0 || cell == nil {
		panic(fmt.Sprintf("data: PinnedTable %q needs two columns, a row count and cells", label))
	}
	for _, col := range cols {
		if col.Width == 0 {
			panic(fmt.Sprintf("data: PinnedTable %q columns need widths", label))
		}
	}
	t := c.Theme()
	h := t.Space(10)
	root := frame(c, label).Row().AlignItems(ui.Start)
	part := func(cs []Column, first int) {
		line(c, ui.Row(c).Background(t.Surface).Height(h), cs, nil, func(i int) {
			kit.Muted(c, cs[i].Title).Bold().FontSize(t.Rem(0.8)).SingleLine()
		})
		for r := range n {
			c.Key(r)
			line(c, hovered(c, false).Height(h), cs, nil, func(i int) { cell(r, first+i) })
		}
	}
	return root.Children(func() {
		ui.Column(c).Shrink(0).BorderEndWidth(1).BorderColor(t.Border).Children(func() { part(cols[:1], 0) })
		ui.ScrollHorizontal(c).Grow(1).Children(func() {
			ui.Column(c).Children(func() { part(cols[1:], 1) })
		})
	})
}
