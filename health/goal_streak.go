package health

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GoalStreakOptions word the card; Day and Days take the streak.
type GoalStreakOptions struct{ Title, Day, Days, Best, Met, Missed string }

// GoalStreak counts days in a row a goal was met.
func GoalStreak(c *ui.Context, days []bool, o GoalStreakOptions) ui.Element {
	def(&o.Title, "Goal streak")
	if len(days) == 0 || len(days) > 31 {
		panic(fmt.Sprintf("health: GoalStreak %q needs 1 to 31 days", o.Title))
	}
	def(&o.Day, "%d day in a row")
	def(&o.Days, "%d days in a row")
	def(&o.Best, "Best %d")
	def(&o.Met, "Day %d, met")
	def(&o.Missed, "Day %d, missed")
	streak, best, run := 0, 0, 0
	for _, m := range days {
		if m {
			run++
		} else {
			run = 0
		}
		best = max(best, run)
	}
	streak = run
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ink := t.Warning
			if streak == 0 {
				ink = t.TextMuted
			}
			flame := tile(c, "flame", ink, t.Space(12))
			kit.Arrive(flame, ui.Pose{Scale: 0.5, Rotate: -20, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				word := fmt.Sprintf(o.Days, streak)
				if streak == 1 {
					word = fmt.Sprintf(o.Day, streak)
				}
				kit.Numeric(c, word).FontSize(t.Rem(1.15)).Bold().FillWidth().MinWidth(0).SingleLine().
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
			kit.Chip(c, fmt.Sprintf(o.Best, best), t.Accent)
		})
		dots := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5))
		dots.Children(func() {
			for i, m := range days {
				label := fmt.Sprintf(o.Missed, i+1)
				if m {
					label = fmt.Sprintf(o.Met, i+1)
				}
				d := ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Radius(999).Label(label)
				inTail := m && i >= len(days)-streak
				switch {
				case inTail:
					d.Background(t.Warning)
				case m:
					d.Background(t.Warning.Alpha(0.35))
				default:
					d.Border(1, t.Border)
				}
				kit.Arrive(d, ui.Pose{Scale: 0.3, Fade: 1}, time.Duration(i)*25*time.Millisecond)
			}
		})
	})
	return root
}
