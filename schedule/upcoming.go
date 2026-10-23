package schedule

import "github.com/egoist/mygo/ui"

// Upcoming shows the next events list.
func Upcoming(c *ui.Context, events []Event) ui.Element {
	if len(events) == 0 {
		panic("schedule: Upcoming needs events")
	}
	validEvents("Upcoming", events)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, event := range events {
			c.Key(event.ID)
			EventCard(c, event)
		}
	})
	return root
}
