package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MilestoneOptions are the milestone labels.
type MilestoneOptions struct{ Date string }

// Milestone shows a diamond with date and progress.
func Milestone(c *ui.Context, title, date string, value float64, o MilestoneOptions) ui.Element {
	if title == "" || date == "" || value < 0 || value > 1 || !finite(value) {
		panic("project: Milestone has invalid data")
	}
	def(&o.Date, "Due %s")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Box(c).Size(t.Space(5), t.Space(5)).Rotate(45).Background(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, title).Bold()
			kit.Muted(c, fmt.Sprintf(o.Date, date)).FontSize(t.Rem(0.75))
		})
		progress(c, value, t.Success)
	})
	return root
}
