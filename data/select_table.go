package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SelectTableOptions tune a SelectTable.
type SelectTableOptions struct {
	// All names the header check; Row takes the number.
	All, Row string
	// Height scrolls the rows within it; zero shows them all.
	Height float32
}

// SelectTable is a table with row checks.
func SelectTable(c *ui.Context, label string, cols []Column, n int, chosen *ui.Selection[int], cell func(row, col int), o SelectTableOptions) ui.Element {
	columns("SelectTable", label, cols)
	if n < 0 || chosen == nil || cell == nil {
		panic(fmt.Sprintf("data: SelectTable %q needs a row count, a selection and cells", label))
	}
	def(&o.All, "Select all")
	def(&o.Row, "Select row %d")
	root := frame(c, label)
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	count := 0
	for i := range n {
		if chosen.Has(i) {
			count++
		}
	}
	changed := func(event string, args ...any) {
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("SelectTable", event, args...)
	}
	root.Children(func() {
		head(c, root, "SelectTable", cols, nil, func() {
			all := n > 0 && count == n
			if checkbox(c, o.All, &all, count > 0 && count < n).Changed() {
				for i := range n {
					if all {
						chosen.Add(i)
					} else {
						chosen.Remove(i)
					}
				}
				changed("all", "on", all)
			}
		})
		body(c, root, list, n, o.Height, func(i int) {
			on := chosen.Has(i)
			line(c, hovered(c, on), cols, func() {
				if checkbox(c, fmt.Sprintf(o.Row, i+1), &on, false).Changed() {
					toggle(chosen, i)
					changed("row", "index", i, "on", on)
				}
			}, func(col int) { cell(i, col) })
		})
	})
	return root
}
