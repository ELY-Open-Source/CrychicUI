package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Shift is one staff shift; Role is optional.
type Shift struct {
	Name, Role, Day, Start, End string
}

// ShiftsOptions word the roster.
type ShiftsOptions struct{ Title string }

// Shifts shows staff shifts grouped by day.
func Shifts(c *ui.Context, shifts []Shift, o ShiftsOptions) ui.Element {
	if len(shifts) == 0 {
		panic("schedule: Shifts needs shifts")
	}
	def(&o.Title, "Shifts")
	t := c.Theme()
	var days []string
	byDay := map[string][]int{}
	for i, shift := range shifts {
		if shift.Name == "" || shift.Day == "" || shift.Start == "" || shift.End == "" {
			panic(fmt.Sprintf("schedule: Shifts has an invalid shift for %q", shift.Name))
		}
		if _, ok := byDay[shift.Day]; !ok {
			days = append(days, shift.Day)
		}
		byDay[shift.Day] = append(byDay[shift.Day], i)
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "briefcase", o.Title, func() { kit.Chip(c, fmt.Sprint(len(shifts)), t.Accent) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), 0, t.Space(2), 0)
		list.Children(func() {
			for _, day := range days {
				c.Key("day." + day)
				kit.Muted(c, day).FontSize(t.Rem(0.7)).Bold().Padding(t.Space(2), t.Space(3), t.Space(1), t.Space(3))
				for _, i := range byDay[day] {
					shift := shifts[i]
					c.Key(i)
					row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3))
					row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(0.6*unit(row.Hovered())), ui.Snappy))
					kit.Stagger(c, list, row, i).Children(func() {
						display.Avatar(c, shift.Name, nil).Size(t.Space(8), t.Space(8))
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							ui.Text(c, shift.Name).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
							if shift.Role != "" {
								kit.Muted(c, shift.Role).FontSize(t.Rem(0.7)).SingleLine().MinWidth(0)
							}
						})
						ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).
							Background(t.Accent.Alpha(0.1)).TextColor(t.Accent).Children(func() {
							glyph(c, "clock", 12)
							kit.Numeric(c, shift.Start+"–"+shift.End).FontSize(t.Rem(0.7)).Bold()
						})
					})
				}
			}
		})
	})
	return root
}
