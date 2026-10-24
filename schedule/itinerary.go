package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stop is one itinerary stop; Icon names a Lucide glyph.
type Stop struct {
	Time, Title, Place, Icon string
}

// TripDay is one itinerary day and stops.
type TripDay struct {
	Label, Date string
	Stops       []Stop
}

// Itinerary shows trip days with stops on a rail.
func Itinerary(c *ui.Context, days []TripDay) ui.Element {
	if len(days) == 0 {
		panic("schedule: Itinerary needs days")
	}
	for _, day := range days {
		if day.Label == "" || len(day.Stops) == 0 {
			panic(fmt.Sprintf("schedule: Itinerary day %q needs stops", day.Label))
		}
		for _, stop := range day.Stops {
			if stop.Title == "" || stop.Icon == "" {
				panic(fmt.Sprintf("schedule: Itinerary day %q has a stop without title or icon", day.Label))
			}
			icon.Must(stop.Icon)
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		n := 0
		for i, day := range days {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					kit.Chip(c, day.Label, t.Accent)
					if day.Date != "" {
						kit.Muted(c, day.Date).FontSize(t.Rem(0.75)).Shrink(0)
					}
				})
				for k, stop := range day.Stops {
					last := k == len(day.Stops)-1
					c.Key(k)
					row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
					kit.Stagger(c, root, row, n).Children(func() {
						kit.Numeric(c, stop.Time).FontSize(t.Rem(0.75)).Bold().Width(t.Space(11)).Shrink(0).Padding(t.Space(2), 0, 0, 0)
						rail := ui.Column(c).Width(t.Space(8)).Shrink(0).AlignItems(ui.Center)
						rail.Draw(func(p *ui.Painter, r ui.Rect) {
							if !last {
								x := r.X + r.W/2
								p.Line(x, r.Y+t.Space(8), x, r.Y+r.H, 1, t.Border)
							}
						})
						rail.Children(func() { tile(c, stop.Icon, t.Accent) })
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(1), 0, t.Space(3), 0).Children(func() {
							ui.Text(c, stop.Title).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
							if stop.Place != "" {
								meta(c, "map-pin", stop.Place)
							}
						})
					})
					n++
				}
			})
		}
	})
	return root
}
