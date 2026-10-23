package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TimeOffRequest is one leave request.
type TimeOffRequest struct {
	Name, Range, Status string
}

// TimeOff shows requests with status.
func TimeOff(c *ui.Context, requests []TimeOffRequest) ui.Element {
	if len(requests) == 0 {
		panic("schedule: TimeOff needs requests")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, request := range requests {
			if request.Name == "" || request.Range == "" || request.Status == "" {
				panic("schedule: TimeOff has an invalid request")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, request.Name)
					kit.Muted(c, request.Range).FontSize(t.Rem(0.75))
				})
				ui.Text(c, request.Status).TextColor(t.Accent)
			})
		}
	})
	return root
}
