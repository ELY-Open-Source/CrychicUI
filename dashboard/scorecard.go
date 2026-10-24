package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScoreGoal is one scorecard goal.
type ScoreGoal struct {
	Label           string
	Current, Target float64
}

// ScorecardOptions word the card; Count and Of take numbers.
type ScorecardOptions struct{ Title, Met, Behind, Count, Of string }

// Scorecard lists goals with progress and a status per row.
func Scorecard(c *ui.Context, goals []ScoreGoal, o ScorecardOptions) ui.Element {
	if len(goals) == 0 {
		panic("dashboard: Scorecard needs goals")
	}
	names, met := make([]string, len(goals)), 0
	for i, goal := range goals {
		if goal.Current < 0 || goal.Target <= 0 || !finite(goal.Current) || !finite(goal.Target) || !finite(goal.Current/goal.Target) {
			panic(fmt.Sprintf("dashboard: Scorecard goal %q needs current of zero or more and a positive target", goal.Label))
		}
		names[i] = goal.Label
		if goal.Current >= goal.Target {
			met++
		}
	}
	unique("Scorecard", names)
	def(&o.Title, "Scorecard")
	def(&o.Met, "Met")
	def(&o.Behind, "Behind")
	def(&o.Count, "%d of %d met")
	def(&o.Of, "%s of %s")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "clipboard-check", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Count, met, len(goals)), t.Success.Mix(t.Warning, unit(met < len(goals))))
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			for i, goal := range goals {
				c.Key(goal.Label)
				done := goal.Current >= goal.Target
				text, ink, glyph := o.Behind, t.Warning, "circle-alert"
				if done {
					text, ink, glyph = o.Met, t.Success, "circle-check"
				}
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(goals)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(16, 16).Shrink(0).TextColor(ink)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, goal.Label).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							kit.Muted(c, fmt.Sprintf(o.Of, number(c, goal.Current), number(c, goal.Target))).FontSize(t.Rem(0.75)).Shrink(0)
						})
						progress(c, min(1, goal.Current/goal.Target), ink)
					})
					kit.Chip(c, text, ink)
				})
			}
		})
	})
	return root
}
