package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Resource is one room and its events.
type Resource struct {
	Name   string
	Events []Event
}

// ResourceRows shows rooms as rows on a time axis.
func ResourceRows(c *ui.Context, resources []Resource) ui.Element {
	if len(resources) == 0 {
		panic("schedule: ResourceRows needs resources")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			ui.Box(c).Width(t.Space(20))
			ui.Row(c).Grow(1).Justify(ui.SpaceBetween).Children(func() {
				for _, hour := range []int{0, 6, 12, 18, 23} {
					ui.Text(c, fmt.Sprintf("%02d", hour)).FontSize(t.Rem(0.65)).TextColor(t.TextMuted)
				}
			})
		})
		for i, resource := range resources {
			if resource.Name == "" {
				panic("schedule: ResourceRows room must not be empty")
			}
			validEvents("ResourceRows", resource.Events)
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, resource.Name).Width(t.Space(20))
				track := ui.Box(c).Grow(1).Height(t.Space(8)).Clip()
				kit.Settle(c, track, []float32{track.LayoutBounds().W}, func() []float32 { return []float32{track.LayoutBounds().W} })
				track.Children(func() {
					for _, event := range resource.Events {
						start := float32(event.Start.Hour()*60+event.Start.Minute()) / (24 * 60)
						duration := float32(event.End.Sub(event.Start)) / float32(24*time.Hour)
						x := track.LayoutBounds().W * start
						w := track.LayoutBounds().W * duration
						c.Key(event.ID)
						ui.Box(c).Absolute().Left(x).Top(0).Width(w).Height(t.Space(8)).Radius(t.Space(0.5)).
							Background(t.Accent.Alpha(0.2)).Children(func() {
							ui.Text(c, event.Title).FontSize(t.Rem(0.65)).Padding(t.Space(0.5)).SingleLine()
						})
					}
				})
			})
		}
	})
	return root
}
