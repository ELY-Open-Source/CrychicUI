package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ExpandTableOptions tune an ExpandTable.
type ExpandTableOptions struct {
	// Expand names an arrow; it takes the row.
	Expand string
	// Height scrolls the rows within it; zero shows them all.
	Height float32
}

// ExpandTable rows open a springing detail row.
func ExpandTable(c *ui.Context, label string, cols []Column, n int, open *ui.Selection[int], cell func(row, col int), detail func(row int), o ExpandTableOptions) ui.Element {
	columns("ExpandTable", label, cols)
	if n < 0 || open == nil || cell == nil || detail == nil {
		panic(fmt.Sprintf("data: ExpandTable %q needs a row count, open rows, cells and details", label))
	}
	def(&o.Expand, "Details for row %d")
	t := c.Theme()
	root := frame(c, label)
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	root.Children(func() {
		head(c, root, "ExpandTable", cols, nil, func() { ui.Box(c).Width(t.Space(6)).Shrink(0) })
		body(c, root, list, n, o.Height, func(i int) {
			shown := open.Has(i)
			ui.Column(c).FillWidth().Children(func() {
				line(c, hovered(c, false), cols, func() {
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Expand, i+1)).Size(t.Space(6), t.Space(6)).Shrink(0).Center().
						Radius(t.Radius)).Expanded(shown)
					b.Children(func() { chevron(c, shown) })
					if b.Clicked() {
						toggle(open, i)
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("ExpandTable", "toggled", "index", i, "open", !shown)
					}
				}, func(col int) { cell(i, col) })
				if shown {
					ui.Column(c).Key("detail").FillWidth().Clip().Padding(t.Space(3), t.Space(3), t.Space(3), t.Space(12)).
						Background(t.Surface).Transition(spring).Children(func() { detail(i) })
				}
			})
		})
	})
	return root
}
