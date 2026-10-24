package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// YearViewOptions name the months and weekdays; Now overrides the clock.
type YearViewOptions struct {
	Weekdays [7]string
	Months   [12]string
	Now      time.Time
}

// YearView shows twelve mini months; today is ringed.
func YearView(c *ui.Context, year int, o YearViewOptions) ui.Element {
	if year < 1 {
		panic("schedule: YearView needs a year")
	}
	if o.Weekdays == [7]string{} {
		o.Weekdays = [7]string{"S", "M", "T", "W", "T", "F", "S"}
	}
	if o.Months == [12]string{} {
		o.Months = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	}
	now := o.Now
	if now.IsZero() {
		now = c.Now()
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(4).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for month := time.January; month <= time.December; month++ {
			c.Key(int(month))
			panel := kit.Card(c).Gap(t.Space(1)).Padding(t.Space(2))
			current := now.Year() == year && now.Month() == month
			kit.Stagger(c, root, panel, int(month)-1).Children(func() {
				ui.Text(c, o.Months[month-1]).FontSize(t.Rem(0.75)).Bold().TextColor(t.Text.Mix(t.Accent, unit(current)))
				ui.Grid(c).Columns(7).FillWidth().Children(func() {
					for _, day := range o.Weekdays {
						ui.Box(c).Height(t.Space(4)).Center().Children(func() { kit.Muted(c, day).FontSize(t.Rem(0.5)).Bold() })
					}
					first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
					for range int(first.Weekday()) {
						ui.Box(c).Height(t.Space(4))
					}
					days := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
					for day := 1; day <= days; day++ {
						today := current && now.Day() == day
						weekend := (int(first.Weekday())+day-1)%7 == 0 || (int(first.Weekday())+day-1)%7 == 6
						ink := t.Text
						if weekend {
							ink = t.TextMuted
						}
						if today {
							ink = t.AccentText
						}
						ui.Box(c).Height(t.Space(4)).Center().Children(func() {
							ui.Box(c).Size(t.Space(4), t.Space(4)).Center().Radius(999).Background(t.Accent.Alpha(unit(today))).Children(func() {
								ui.Text(c, fmt.Sprint(day)).FontSize(t.Rem(0.5)).TextColor(ink)
							})
						})
					}
				})
			})
		}
	})
	return root
}
