package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Activity is one recent event.
type Activity struct {
	ID, Text, Time string
}

// ActivityFeed shows recent events; rows slide in.
func ActivityFeed(c *ui.Context, events []Activity) ui.Element {
	if len(events) == 0 {
		panic("dashboard: ActivityFeed needs events")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, event := range events {
			if event.ID == "" || event.Text == "" || event.Time == "" {
				panic("dashboard: ActivityFeed has an invalid event")
			}
			c.Key(event.ID)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(row, ui.Pose{X: -t.Space(3), Fade: 1}, 0)
			row.Children(func() {
				ui.Text(c, event.Text).Grow(1).Wrap()
				kit.Muted(c, event.Time).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
