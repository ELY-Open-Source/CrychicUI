package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Timeline staggers vertical events with dots and a line.
func Timeline(c *ui.Context, events []string) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
		for i, event := range events {
			row := ui.Row(c).Key(i).Gap(t.Space(2)).AlignItems(ui.Start)
			kit.Arrive(row, ui.Pose{Y: 4, Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Column(c).Width(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Accent)
					if i < len(events)-1 {
						ui.Box(c).Width(1).Height(t.Space(6)).Background(t.Border)
					}
				})
				ui.Text(c, event).Grow(1)
			})
		}
	})
}
