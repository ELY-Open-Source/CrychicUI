package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DayViewOptions are the day view options.
type DayViewOptions struct{ Now time.Time }

// DayView shows hours with events and a moving now line.
func DayView(c *ui.Context, day time.Time, events []Event, o DayViewOptions) ui.Element {
	validEvents("DayView", events)
	t := c.Theme()
	now := o.Now
	if now.IsZero() {
		now = c.Now()
	}
	now = now.In(day.Location())
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	c.After(time.Second)
	root.Children(func() {
		ui.Text(c, day.Format("2006-01-02")+" "+now.Format("15:04:05")).Bold()
		for hour := range 24 {
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				kit.Muted(c, fmt.Sprintf("%02d:00", hour)).FontSize(t.Rem(0.7))
				ui.Box(c).Grow(1).Height(t.Space(4)).Children(func() {
					for _, event := range events {
						local := event.Start.In(day.Location())
						if sameDate(local, day) && local.Hour() == hour {
							ui.Text(c, event.Title).FontSize(t.Rem(0.75)).TextColor(t.Accent)
						}
					}
					if sameDate(now, day) && now.Hour() == hour {
						y := float32(now.Minute()) / 60 * t.Space(4)
						ui.Box(c).Absolute().Left(0).Right(0).Top(y).Height(1).Background(t.Danger)
					}
				})
			})
		}
	})
	return root
}
