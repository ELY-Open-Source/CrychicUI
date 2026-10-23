package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Shift is one staff shift.
type Shift struct {
	Name, Day, Start, End string
}

// Shifts shows staff shifts by day.
func Shifts(c *ui.Context, shifts []Shift) ui.Element {
	if len(shifts) == 0 {
		panic("schedule: Shifts needs shifts")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, shift := range shifts {
			if shift.Name == "" || shift.Day == "" || shift.Start == "" || shift.End == "" {
				panic("schedule: Shifts has an invalid shift")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, shift.Name).Grow(1)
				kit.Muted(c, shift.Day+" "+shift.Start+"–"+shift.End).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
