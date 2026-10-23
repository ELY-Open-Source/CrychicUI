package shop

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TrackEvent is one carrier tracking event.
type TrackEvent struct {
	Label, Location string
	Time            time.Time
	Done            bool
}

// Tracking shows carrier events on a line.
func Tracking(c *ui.Context, events []TrackEvent) ui.Element {
	if len(events) == 0 {
		panic("shop: Tracking needs events")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, event := range events {
			if event.Label == "" || event.Location == "" {
				panic("shop: Tracking event needs text and location")
			}
			c.Key(i)
			ink := t.TextMuted
			if event.Done {
				ink = t.Success
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, event.Label).FontSize(t.Rem(0.85))
					kit.Muted(c, event.Location).FontSize(t.Rem(0.75))
				})
				kit.Muted(c, event.Time.Format("Jan 02 15:04")).FontSize(t.Rem(0.75))
			})
			kit.Arrive(ui.Box(c).Width(1).Height(t.Space(3)).Margin(0, 0, 0, t.Space(1)).Background(t.Border), ui.Pose{Scale: 0.5, Fade: 1}, 0)
		}
	})
	return root
}
