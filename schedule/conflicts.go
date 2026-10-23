package schedule

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// ConflictsOptions are the conflict labels.
type ConflictsOptions struct{ Format string }

// Conflicts flags overlapping events.
func Conflicts(c *ui.Context, events []Event, o ConflictsOptions) ui.Element {
	if len(events) == 0 {
		panic("schedule: Conflicts needs events")
	}
	validEvents("Conflicts", events)
	def(&o.Format, "%d conflicts")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	count := 0
	conflicting := map[string]bool{}
	for i, a := range events {
		for j := i + 1; j < len(events); j++ {
			if a.Start.Before(events[j].End) && events[j].Start.Before(a.End) {
				count++
				conflicting[a.ID] = true
				conflicting[events[j].ID] = true
			}
		}
	}
	root.Children(func() {
		ui.Text(c, fmt.Sprintf(o.Format, count)).Bold().TextColor(t.Danger)
		for _, event := range events {
			c.Key(event.ID)
			border := t.Border
			if conflicting[event.ID] {
				border = t.Danger
			}
			card := ui.Box(c).Border(1, border).Radius(t.Radius)
			card.Children(func() { EventCard(c, event) })
		}
	})
	return root
}
