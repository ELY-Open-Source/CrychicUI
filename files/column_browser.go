package files

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ColumnBrowser shows Finder-style columns; cue Select.
func ColumnBrowser(c *ui.Context, label string, columns [][]string, selected *[]int) ui.Element {
	if label == "" || len(columns) == 0 || selected == nil || len(*selected) != len(columns) {
		panic("files: ColumnBrowser needs columns and selections")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(0).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, column := range columns {
			if len(column) == 0 {
				panic("files: ColumnBrowser columns must not be empty")
			}
			c.Key(i)
			col := ui.Column(c).Width(t.Space(48)).Shrink(0).BorderWidth(0, 1, 0, 0).BorderColor(t.Border)
			kit.Arrive(col, ui.Pose{X: t.Space(4), Fade: 1}, 0)
			col.Children(func() {
				for j, name := range column {
					if name == "" {
						panic("files: ColumnBrowser item must not be empty")
					}
					if (*selected)[i] < -1 || (*selected)[i] >= len(column) {
						panic("files: ColumnBrowser selection out of range")
					}
					row := fileButton(c, name, (*selected)[i] == j).Radius(0)
					row.Children(func() { ui.Text(c, name).SingleLine() })
					if row.Clicked() && (*selected)[i] != j {
						(*selected)[i] = j
						for k := i + 1; k < len(*selected); k++ {
							(*selected)[k] = -1
						}
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("ColumnBrowser", "selected", "column", i, "index", j)
					}
				}
			})
		}
	})
	return root
}
