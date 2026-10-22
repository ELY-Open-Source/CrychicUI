package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/egoist/mygo/ui"
)

// LoadingTable is a table head over sheen rows.
func LoadingTable(c *ui.Context, label string, cols []Column, rows int) ui.Element {
	columns("LoadingTable", label, cols)
	if rows <= 0 {
		panic(fmt.Sprintf("data: LoadingTable %q needs rows", label))
	}
	t := c.Theme()
	root := frame(c, label).Role(ui.RoleProgress).Range(0, 1, -1)
	return root.Children(func() {
		head(c, root, "LoadingTable", cols, nil, nil)
		for i := range rows {
			c.Key(i)
			line(c, ui.Row(c), cols, nil, func(int) { feedback.Skeleton(c).Height(t.Space(3)) })
		}
	})
}
