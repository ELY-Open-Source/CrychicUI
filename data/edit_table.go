package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EditTableOptions tune an EditTable.
type EditTableOptions struct {
	// Edit takes the column and row number.
	Edit string
	// Height scrolls the rows within it; zero shows them all.
	Height float32
}

// EditTable edits cells in place; at is -1, -1 idle.
func EditTable(c *ui.Context, label string, cols []Column, rows [][]string, at *[2]int, o EditTableOptions) ui.Element {
	columns("EditTable", label, cols)
	if at == nil || at[0] < -1 || at[0] >= len(rows) || at[1] < -1 || at[1] >= len(cols) || (at[0] < 0) != (at[1] < 0) {
		panic(fmt.Sprintf("data: EditTable %q needs a cell in the table or -1, -1", label))
	}
	for _, r := range rows {
		if len(r) != len(cols) {
			panic(fmt.Sprintf("data: EditTable %q rows need a cell per column", label))
		}
	}
	def(&o.Edit, "Edit %s, row %d")
	t := c.Theme()
	root := frame(c, label)
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	draft := ui.Local(root, "draft", func() string { return "" })
	drafted := ui.Local(root, "drafted", func() [2]int { return [2]int{-1, -1} })
	if *at != *drafted {
		*drafted = *at
		if at[0] >= 0 {
			*draft = rows[at[0]][at[1]]
		}
	}
	root.Children(func() {
		head(c, root, "EditTable", cols, nil, nil)
		body(c, root, list, len(rows), o.Height, func(i int) {
			line(c, hovered(c, false), cols, nil, func(col int) {
				if *at == [2]int{i, col} {
					c.Key(*at)
					in := ui.TextInput(c, draft).Label(cols[col].Title).FillWidth().AutoFocus()
					switch {
					case in.Submitted():
						rows[i][col], *at = *draft, [2]int{-1, -1}
						root.SetChanged()
						sound.Play(sound.Success)
						kit.Log("EditTable", "saved", "row", i, "column", col)
					case in.Shortcut(0, ui.KeyEscape):
						*at = [2]int{-1, -1}
						root.SetChanged()
						kit.Log("EditTable", "cancelled")
					}
					return
				}
				side := ui.Start
				if cols[col].End {
					side = ui.End
				}
				b := ui.ButtonBase(c).Label(fmt.Sprintf(o.Edit, cols[col].Title, i+1)).FillWidth().Justify(side).Padding(t.Space(1), 0).Radius(t.Radius / 2)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Text(c, rows[i][col]).SingleLine()
				})
				if b.Clicked() {
					*at, *draft, *drafted = [2]int{i, col}, rows[i][col], [2]int{i, col}
					root.SetChanged()
					kit.Log("EditTable", "editing", "row", i, "column", col)
				}
			})
		})
	})
	return root
}
