package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EventCard shows title, time, place, and a color bar.
func EventCard(c *ui.Context, event Event) ui.Element {
	validEvent("EventCard", event)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{X: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Width(t.Space(1)).Height(t.Space(8)).Radius(999).Background(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, event.Title).Bold()
			kit.Muted(c, eventText(event)).FontSize(t.Rem(0.75))
		})
		if event.Place != "" {
			kit.Muted(c, event.Place).FontSize(t.Rem(0.75))
		}
	})
	return root
}
