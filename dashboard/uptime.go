package dashboard

import "github.com/egoist/mygo/ui"

// Status is a service health state.
type Status uint8

const (
	Up Status = iota
	Degraded
	Down
)

// UptimeDay is one day's service status.
type UptimeDay struct {
	Label  string
	Status Status
}

// Uptime shows day bars colored by status.
func Uptime(c *ui.Context, days []UptimeDay) ui.Element {
	if len(days) == 0 {
		panic("dashboard: Uptime needs days")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(0.5)).AlignItems(ui.End).Role(ui.RoleGroup)
	root.Children(func() {
		for i, day := range days {
			if day.Label == "" || day.Status > Down {
				panic("dashboard: Uptime has an invalid day")
			}
			c.Key(i)
			ink := t.Success
			if day.Status == Degraded {
				ink = t.Warning
			} else if day.Status == Down {
				ink = t.Danger
			}
			ui.Column(c).Grow(1).Basis(0).Gap(t.Space(0.5)).AlignItems(ui.Center).Children(func() {
				bar := ui.Box(c).FillWidth().Height(t.Space(10)).Radius(t.Space(0.5)).Tooltip(day.Label)
				bar.Background(bar.AnimateColor("status", ink, ui.Snappy))
			})
		}
	})
	return root
}
