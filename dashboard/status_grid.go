package dashboard

import "github.com/egoist/mygo/ui"

// Service is one service and health status.
type Service struct {
	Name   string
	Status Status
}

// StatusGrid shows services as colored cells.
func StatusGrid(c *ui.Context, services []Service) ui.Element {
	if len(services) == 0 {
		panic("dashboard: StatusGrid needs services")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, service := range services {
			if service.Name == "" || service.Status > Down {
				panic("dashboard: StatusGrid has an invalid service")
			}
			c.Key(i)
			ink := t.Success
			if service.Status == Degraded {
				ink = t.Warning
			} else if service.Status == Down {
				ink = t.Danger
			}
			row := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
			row.Background(row.AnimateColor("status", ink.Alpha(0.15), ui.Snappy))
			row.Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
				ui.Text(c, service.Name).FontSize(t.Rem(0.8))
			})
		}
	})
	return root
}
