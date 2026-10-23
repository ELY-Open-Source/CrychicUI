package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SlotDay is one date and available count.
type SlotDay struct {
	Date      string
	Available int
}

// SlotCalendar shows dates with available counts.
func SlotCalendar(c *ui.Context, days []SlotDay) ui.Element {
	if len(days) == 0 {
		panic("schedule: SlotCalendar needs days")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(7).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, day := range days {
			if day.Date == "" || day.Available < 0 {
				panic("schedule: SlotCalendar has an invalid day")
			}
			c.Key(i)
			ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(t.Space(0.5)).Border(1, t.Border).Children(func() {
				ui.Text(c, day.Date).FontSize(t.Rem(0.7))
				kit.Numeric(c, fmt.Sprint(day.Available)).TextColor(t.Accent)
			})
		}
	})
	return root
}
