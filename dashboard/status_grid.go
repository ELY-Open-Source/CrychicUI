package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Service is one service and status; Detail is a reading.
type Service struct {
	Name, Detail string
	Status       Status
}

// StatusGridOptions word the card; States name each status.
type StatusGridOptions struct {
	Title, AllUp, Issue, Issues string
	States                      [3]string
}

// StatusGrid shows services as colored cells under an overall state.
func StatusGrid(c *ui.Context, services []Service, o StatusGridOptions) ui.Element {
	if len(services) == 0 {
		panic("dashboard: StatusGrid needs services")
	}
	names, issues := make([]string, len(services)), 0
	for i, service := range services {
		if service.Status > Down {
			panic(fmt.Sprintf("dashboard: StatusGrid service %q has an unknown status", service.Name))
		}
		names[i] = service.Name
		if service.Status != Up {
			issues++
		}
	}
	unique("StatusGrid", names)
	def(&o.Title, "Services")
	def(&o.AllUp, "All systems operational")
	def(&o.Issue, "%d issue")
	def(&o.Issues, "%d issues")
	for i, name := range [...]string{"Operational", "Degraded", "Outage"} {
		def(&o.States[i], name)
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "server", o.Title, func() {
			switch {
			case issues == 0:
				kit.Chip(c, o.AllUp, t.Success)
			case issues == 1:
				kit.Chip(c, fmt.Sprintf(o.Issue, issues), t.Warning)
			default:
				kit.Chip(c, fmt.Sprintf(o.Issues, issues), t.Danger)
			}
		})
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		grid.Children(func() {
			for i, service := range services {
				c.Key(service.Name)
				ink := statusInk(t, service.Status)
				cell := ui.Row(c).Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5)).Radius(t.Radius).MinWidth(0)
				cell.Background(cell.AnimateColor("bg", ink.Alpha(0.08+0.06*unit(cell.Hovered())), ui.Snappy))
				kit.Stagger(c, grid, cell, i).Children(func() {
					ui.Icon(c, icon.Must(statusGlyph(service.Status))).Size(16, 16).Shrink(0).TextColor(ink).
						ContentTransition(ui.ContentTransition{})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, service.Name).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						meta := o.States[service.Status]
						if service.Detail != "" {
							meta += " · " + service.Detail
						}
						kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine().ContentTransition(ui.ContentTransition{})
					})
				})
			}
		})
	})
	return root
}
