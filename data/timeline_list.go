package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Moment is a dated title and text.
type Moment struct{ Date, Title, Text string }

// TimelineList lays dated rows on a line; they stagger in.
func TimelineList(c *ui.Context, label string, moments []Moment) ui.Element {
	if label == "" || len(moments) == 0 {
		panic(fmt.Sprintf("data: TimelineList %q needs moments", label))
	}
	for _, m := range moments {
		if m.Date == "" || m.Title == "" {
			panic(fmt.Sprintf("data: TimelineList %q moments need a date and title", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, m := range moments {
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
			kit.Arrive(row, ui.Pose{X: t.Space(-3), Fade: 1}, time.Duration(i)*60*time.Millisecond).Children(func() {
				kit.Muted(c, m.Date).Width(t.Space(24)).Shrink(0).TextAlign(ui.End).Padding(t.Space(0.5), 0, 0, 0)
				ui.Column(c).Width(t.Space(3)).Shrink(0).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(3), t.Space(3)).Shrink(0).Radius(999).Background(t.Accent).Margin(t.Space(1), 0, 0, 0)
					if i < len(moments)-1 {
						ui.Box(c).Width(1).Grow(1).Background(t.Border)
					}
				})
				ui.Column(c).Grow(1).Gap(t.Space(0.5)).Padding(0, 0, t.Space(4), 0).Children(func() {
					ui.Text(c, m.Title).Bold()
					if m.Text != "" {
						kit.Muted(c, m.Text)
					}
				})
			})
		}
	})
}
