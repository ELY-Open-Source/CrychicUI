package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Goal is one quarter's current and target values.
type Goal struct {
	Quarter string
	Current float64
	Target  float64
}

// GoalProgress shows goal bars by quarter.
func GoalProgress(c *ui.Context, goals []Goal) ui.Element {
	if len(goals) == 0 {
		panic("dashboard: GoalProgress needs goals")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, goal := range goals {
			if goal.Quarter == "" || goal.Current < 0 || goal.Target <= 0 || !finite(goal.Current) || !finite(goal.Target) {
				panic("dashboard: GoalProgress has an invalid goal")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, goal.Quarter).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f/%.0f", goal.Current, goal.Target)).TextColor(t.TextMuted)
				})
				progress(c, min(1, goal.Current/goal.Target), t.Accent)
			})
		}
	})
	return root
}
