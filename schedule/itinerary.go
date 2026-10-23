package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TripDay is one itinerary day and stops.
type TripDay struct {
	Label string
	Stops []string
}

// Itinerary shows trip days with stops.
func Itinerary(c *ui.Context, days []TripDay) ui.Element {
	if len(days) == 0 {
		panic("schedule: Itinerary needs days")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, day := range days {
			if day.Label == "" || len(day.Stops) == 0 {
				panic("schedule: Itinerary has an invalid day")
			}
			c.Key(i)
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				ui.Text(c, day.Label).Bold()
				for _, stop := range day.Stops {
					if stop == "" {
						panic("schedule: Itinerary stop must not be empty")
					}
					ui.Row(c).Gap(t.Space(1)).Children(func() {
						ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Accent)
						ui.Text(c, stop).FontSize(t.Rem(0.85))
					})
				}
			})
		}
	})
	return root
}
