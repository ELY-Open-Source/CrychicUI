package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TeamWeek is one person's weekly load fractions.
type TeamWeek struct {
	Person string
	Loads  []float64
}

// TeamLoad shows weeks by person shaded.
func TeamLoad(c *ui.Context, weeks []TeamWeek) ui.Element {
	if len(weeks) == 0 {
		panic("project: TeamLoad needs weeks")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(5).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, week := range weeks {
			if week.Person == "" || len(week.Loads) != 4 {
				panic("project: TeamLoad needs four weeks per person")
			}
			ui.Text(c, week.Person).FontSize(t.Rem(0.75))
			for j, load := range week.Loads {
				if load < 0 || load > 1 || !finite(load) {
					panic("project: TeamLoad has an invalid load")
				}
				c.Key(i*4 + j)
				ui.Text(c, fmt.Sprintf("%.0f%%", load*100)).Center().Padding(t.Space(0.5)).Radius(t.Space(0.5)).
					Background(t.Accent.Alpha(float32(load) * 0.6)).FontSize(t.Rem(0.65))
			}
		}
	})
	return root
}
