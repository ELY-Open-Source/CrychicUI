package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// IncidentEvent is one incident timeline entry.
type IncidentEvent struct {
	Text, Time string
}

// IncidentOptions are the incident labels.
type IncidentOptions struct{ Timeline string }

// Incident shows an incident card with timeline.
func Incident(c *ui.Context, title string, events []IncidentEvent, o IncidentOptions) ui.Element {
	if title == "" || len(events) == 0 {
		panic("dashboard: Incident needs a title and events")
	}
	def(&o.Timeline, "Timeline")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).
		Background(t.Danger.Alpha(0.1)).Border(1, t.Danger)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, title).Bold()
		kit.Muted(c, o.Timeline).FontSize(t.Rem(0.75))
		for i, event := range events {
			if event.Text == "" || event.Time == "" {
				panic("dashboard: Incident has an invalid event")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Danger)
				ui.Text(c, event.Text).Grow(1).Wrap()
				kit.Muted(c, event.Time).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
