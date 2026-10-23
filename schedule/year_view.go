package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// YearViewOptions are the mini calendar names.
type YearViewOptions struct {
	Weekdays [7]string
	Months   [12]string
}

// YearView shows twelve mini months.
func YearView(c *ui.Context, year int, o YearViewOptions) ui.Element {
	if year < 1 {
		panic("schedule: YearView needs a year")
	}
	if o.Weekdays == [7]string{} {
		o.Weekdays = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	}
	if o.Months == [12]string{} {
		o.Months = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(4).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for month := time.January; month <= time.December; month++ {
			panel := ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(panel, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			panel.Children(func() {
				ui.Text(c, o.Months[month-1]).FontSize(t.Rem(0.75)).Bold()
				ui.Grid(c).Columns(7).Gap(t.Space(0.25)).Children(func() {
					for _, day := range o.Weekdays {
						ui.Text(c, day).FontSize(t.Rem(0.55)).TextColor(t.TextMuted)
					}
					first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
					for range int(first.Weekday()) {
						ui.Text(c, "")
					}
					days := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
					for day := 1; day <= days; day++ {
						ui.Text(c, fmt.Sprint(day)).FontSize(t.Rem(0.55)).TextColor(t.TextMuted)
					}
				})
			})
		}
	})
	return root
}
