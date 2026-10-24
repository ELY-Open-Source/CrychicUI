package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Goal is one quarter's current and target values.
type Goal struct {
	Quarter string
	Current float64
	Target  float64
}

// GoalProgressOptions word the card; Of takes current and target.
type GoalProgressOptions struct{ Title, Of, Met string }

// GoalProgress shows goal bars by quarter; met goals turn green.
func GoalProgress(c *ui.Context, goals []Goal, o GoalProgressOptions) ui.Element {
	if len(goals) == 0 {
		panic("dashboard: GoalProgress needs goals")
	}
	names := make([]string, len(goals))
	for i, goal := range goals {
		if goal.Current < 0 || goal.Target <= 0 || !finite(goal.Current) || !finite(goal.Target) {
			panic(fmt.Sprintf("dashboard: GoalProgress %q needs current of zero or more and a positive target", goal.Quarter))
		}
		if !finite(100 * (goal.Current / goal.Target)) {
			panic(fmt.Sprintf("dashboard: GoalProgress %q percent overflows", goal.Quarter))
		}
		names[i] = goal.Quarter
	}
	unique("GoalProgress", names)
	def(&o.Title, "Goals")
	def(&o.Of, "%s of %s")
	def(&o.Met, "%d of %d met")
	met := 0
	for _, goal := range goals {
		if goal.Current >= goal.Target {
			met++
		}
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "target", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Met, met, len(goals)), t.Success) })
		list := ui.Column(c).FillWidth().Gap(t.Space(3.5)).Padding(t.Space(3))
		list.Children(func() {
			for i, goal := range goals {
				c.Key(goal.Quarter)
				share := goal.Current / goal.Target
				done := share >= 1
				ink := t.Accent.Mix(t.Success, unit(done))
				line := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, list, line, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, goal.Quarter).Bold().FontSize(t.Rem(0.85)).Shrink(0)
						kit.Muted(c, fmt.Sprintf(o.Of, number(c, goal.Current), number(c, goal.Target))).FontSize(t.Rem(0.75)).
							Grow(1).MinWidth(0).SingleLine()
						if done {
							ui.Icon(c, icon.Must("circle-check")).Size(14, 14).Shrink(0).TextColor(t.Success)
						}
						kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*share)).Bold().FontSize(t.Rem(0.85)).TextColor(ink).Shrink(0).
							ContentTransition(ui.ContentTransition{Numeric: true})
					})
					progress(c, min(1, share), ink)
				})
			}
		})
	})
	return root
}
