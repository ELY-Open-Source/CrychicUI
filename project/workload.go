package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WorkloadPerson is one person's load fraction.
type WorkloadPerson struct {
	Name string
	Load float64
}

// Workload shows people with load bars.
func Workload(c *ui.Context, people []WorkloadPerson) ui.Element {
	if len(people) == 0 {
		panic("project: Workload needs people")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, person := range people {
			if person.Name == "" || person.Load < 0 || person.Load > 1 || !finite(person.Load) {
				panic("project: Workload has an invalid person")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, person.Name).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f%%", person.Load*100)).TextColor(t.TextMuted)
				})
				progress(c, person.Load, t.Accent)
			})
		}
	})
	return root
}
