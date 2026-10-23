package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WeekViewOptions are the week labels.
type WeekViewOptions struct{ DayFormat string }

// WeekView shows days with hours and event blocks; cue Select.
func WeekView(c *ui.Context, days []time.Time, events []Event, selected *string, o WeekViewOptions) ui.Element {
	if len(days) != 7 || selected == nil {
		panic("schedule: WeekView needs seven days and selection")
	}
	def(&o.DayFormat, "Mon 02")
	validEvents("WeekView", events)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			ui.Box(c).Width(t.Space(8))
			for _, day := range days {
				ui.Text(c, day.Format(o.DayFormat)).FontSize(t.Rem(0.75)).Grow(1).Basis(0)
			}
		})
		for hour := range 24 {
			ui.Row(c).Gap(t.Space(1)).Children(func() {
				kit.Muted(c, fmt.Sprintf("%02d:00", hour)).FontSize(t.Rem(0.7)).Width(t.Space(8))
				for _, day := range days {
					ui.Column(c).Grow(1).Basis(0).Children(func() {
						for _, event := range events {
							local := event.Start.In(day.Location())
							if !sameDate(local, day) || local.Hour() != hour {
								continue
							}
							c.Key(event.ID)
							b := kit.Press(ui.ButtonBase(c).Label(event.Title).FillWidth().Padding(t.Space(0.5)).Radius(t.Space(0.5)).
								Background(t.Accent.Alpha(0.15)).TextColor(t.Accent).Checked(*selected == event.ID))
							b.Children(func() { ui.Text(c, event.Title).FontSize(t.Rem(0.65)).SingleLine() })
							if b.Clicked() && *selected != event.ID {
								*selected = event.ID
								root.SetChanged()
								sound.Play(sound.Select)
								kit.Log("WeekView", "selected", "event", event.ID)
							}
						}
					})
				}
			})
		}
	})
	return root
}
