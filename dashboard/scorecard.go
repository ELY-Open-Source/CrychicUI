package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScoreGoal is one scorecard goal.
type ScoreGoal struct {
	Label           string
	Current, Target float64
}

// ScorecardOptions are the scorecard labels.
type ScorecardOptions struct{ Met, Behind string }

// Scorecard shows goals with status per row.
func Scorecard(c *ui.Context, goals []ScoreGoal, o ScorecardOptions) ui.Element {
	if len(goals) == 0 {
		panic("dashboard: Scorecard needs goals")
	}
	def(&o.Met, "Met")
	def(&o.Behind, "Behind")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, goal := range goals {
			if goal.Label == "" || goal.Current < 0 || goal.Target <= 0 || !finite(goal.Current) || !finite(goal.Target) {
				panic("dashboard: Scorecard has an invalid goal")
			}
			c.Key(i)
			met := goal.Current >= goal.Target
			text, ink := o.Behind, t.Warning
			if met {
				text, ink = o.Met, t.Success
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, goal.Label).Grow(1)
				kit.Numeric(c, fmt.Sprintf("%.0f/%.0f", goal.Current, goal.Target))
				ui.Text(c, text).TextColor(ink)
			})
		}
	})
	return root
}
