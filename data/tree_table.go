package data

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TreeRow is cells with nested rows.
type TreeRow struct {
	Cells    []string
	Children []TreeRow
}

// TreeTableOptions tune a TreeTable.
type TreeTableOptions struct {
	// Expand takes the row's first cell.
	Expand string
	// Height scrolls the rows within it; zero shows them all.
	Height float32
}

// shownRow is a row in view, with depth and path.
type shownRow struct {
	row   *TreeRow
	depth int
	path  string
}

// TreeTable is an outline table; paths join with "/".
func TreeTable(c *ui.Context, label string, cols []Column, rows []TreeRow, open *ui.Selection[string], o TreeTableOptions) ui.Element {
	columns("TreeTable", label, cols)
	if open == nil {
		panic(fmt.Sprintf("data: TreeTable %q needs open rows", label))
	}
	def(&o.Expand, "Expand %s")
	t := c.Theme()
	var shown []shownRow
	var walk func(rs []TreeRow, depth int, at string)
	walk = func(rs []TreeRow, depth int, at string) {
		seen := map[string]bool{}
		for i := range rs {
			r := &rs[i]
			name := ""
			if len(r.Cells) == len(cols) {
				name = r.Cells[0]
			}
			if name == "" || strings.Contains(name, "/") || seen[name] {
				panic(fmt.Sprintf("data: TreeTable %q rows need a cell per column and unique first cells without /", label))
			}
			seen[name] = true
			path := at + name
			shown = append(shown, shownRow{r, depth, path})
			if open.Has(path) {
				walk(r.Children, depth+1, path+"/")
			}
		}
	}
	walk(rows, 0, "")
	root := frame(c, label)
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	list.Key = func(i int) any { return shown[i].path }
	root.Children(func() {
		head(c, root, "TreeTable", cols, nil, nil)
		body(c, root, list, len(shown), o.Height, func(i int) {
			s := shown[i]
			line(c, hovered(c, false), cols, nil, func(col int) {
				if col > 0 {
					ui.Text(c, s.row.Cells[col]).SingleLine()
					return
				}
				ui.Box(c).Width(t.Space(4) * float32(s.depth)).Shrink(0)
				if len(s.row.Children) == 0 {
					ui.Box(c).Width(t.Space(6)).Shrink(0)
				} else {
					on := open.Has(s.path)
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Expand, s.row.Cells[0])).Size(t.Space(6), t.Space(6)).Shrink(0).
						Center().Radius(t.Radius)).Expanded(on)
					b.Children(func() { chevron(c, on) })
					if b.Clicked() {
						toggle(open, s.path)
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("TreeTable", "toggled", "depth", s.depth, "open", !on)
					}
				}
				ui.Text(c, s.row.Cells[0]).SingleLine()
			})
		})
	})
	return root
}
