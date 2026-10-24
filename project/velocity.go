package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Sprint is one sprint's completed points.
type Sprint struct {
	Name   string
	Points int
}

// Velocity shows points per sprint bars.
func Velocity(c *ui.Context, sprints []Sprint) ui.Element {
	if len(sprints) == 0 {
		panic("project: Velocity needs sprints")
	}
	t := c.Theme()
	maxPoints := 0
	for _, sprint := range sprints {
		if sprint.Name == "" || sprint.Points < 0 {
			panic("project: Velocity has an invalid sprint")
		}
		maxPoints = max(maxPoints, sprint.Points)
	}
	if maxPoints == 0 {
		maxPoints = 1
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, sprint := range sprints {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, sprint.Name).Grow(1)
					kit.Numeric(c, fmt.Sprint(sprint.Points)).TextColor(t.TextMuted)
				})
				progress(c, float64(sprint.Points)/float64(maxPoints), t.Accent)
			})
		}
	})
	return root
}
