package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TimesheetOptions are the timesheet labels.
type TimesheetOptions struct{ Total string }

// Timesheet shows hours per day and total.
func Timesheet(c *ui.Context, hours []float64, o TimesheetOptions) ui.Element {
	if len(hours) != 7 {
		panic("schedule: Timesheet needs seven days")
	}
	def(&o.Total, "Total %.1f h")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	total := 0.0
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			for i, value := range hours {
				if value < 0 || !finite(value) {
					panic("schedule: Timesheet hours must be nonnegative")
				}
				c.Key(i)
				ui.Column(c).Grow(1).Basis(0).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, fmt.Sprintf("%.1f", value)).Bold()
				})
				total += value
				if !finite(total) {
					panic("schedule: Timesheet total overflow")
				}
			}
		})
		kit.Numeric(c, fmt.Sprintf(o.Total, total)).Bold()
	})
	return root
}
