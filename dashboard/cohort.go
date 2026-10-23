package dashboard

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// CohortRow is one cohort and retention fractions.
type CohortRow struct {
	Label  string
	Values []float64
}

// Cohort shows a shaded retention grid.
func Cohort(c *ui.Context, rows []CohortRow) ui.Element {
	if len(rows) == 0 {
		panic("dashboard: Cohort needs rows")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(6).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, row := range rows {
			if row.Label == "" || len(row.Values) != 5 {
				panic("dashboard: Cohort rows need five values")
			}
			ui.Text(c, row.Label).FontSize(t.Rem(0.75))
			for j, value := range row.Values {
				if value < 0 || value > 1 || !finite(value) {
					panic("dashboard: Cohort value must be in [0,1]")
				}
				c.Key(i*5 + j)
				cell := ui.Text(c, fmt.Sprintf("%.0f%%", value*100)).Center().Padding(t.Space(1)).Radius(t.Space(0.5)).FontSize(t.Rem(0.75))
				cell.Background(cell.AnimateColor("retention", t.Accent.Alpha(float32(value)*0.7), ui.Snappy))
			}
		}
	})
	return root
}
