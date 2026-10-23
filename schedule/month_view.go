package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MonthView shows month cells with event chips; cue Select.
func MonthView(c *ui.Context, month time.Time, events []Event, selected *string) ui.Element {
	if selected == nil {
		panic("schedule: MonthView needs selection")
	}
	validEvents("MonthView", events)
	t := c.Theme()
	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	start := first.AddDate(0, 0, -int(first.Weekday()))
	root := ui.Grid(c).Columns(7).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i := range 42 {
			day := start.AddDate(0, 0, i)
			ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(0.5)).Radius(t.Space(0.5)).Border(1, t.Border).Children(func() {
				ui.Text(c, fmt.Sprint(day.Day())).FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
				for _, event := range events {
					if !sameDate(event.Start, day) {
						continue
					}
					c.Key(event.ID)
					b := kit.Press(ui.ButtonBase(c).Label(event.Title).Padding(t.Space(0.25)).Radius(t.Space(0.25)).
						Background(t.Accent.Alpha(0.15)).TextColor(t.Accent).Checked(*selected == event.ID))
					b.Children(func() { ui.Text(c, event.Title).FontSize(t.Rem(0.55)).SingleLine() })
					if b.Clicked() && *selected != event.ID {
						*selected = event.ID
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("MonthView", "selected", "event", event.ID)
					}
				}
			})
		}
	})
	return root
}
