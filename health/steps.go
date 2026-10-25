package health

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StepsOptions word the card; Week and Days add bars.
type StepsOptions struct {
	Title, Of, Left, Done string
	Week                  []int
	Days                  []string
}

// Steps shows a step count toward a daily goal.
func Steps(c *ui.Context, steps, goal int, o StepsOptions) ui.Element {
	def(&o.Title, "Steps")
	if steps < 0 || goal < 1 || len(o.Week) != len(o.Days) || len(o.Week) > 7 {
		panic(fmt.Sprintf("health: Steps %q needs steps of 0 or more, a goal and a label for each of up to 7 days", o.Title))
	}
	top := goal
	for _, n := range o.Week {
		if n < 0 {
			panic(fmt.Sprintf("health: Steps %q has a day of %d steps", o.Title, n))
		}
		top = max(top, n)
	}
	def(&o.Of, "of %s")
	def(&o.Left, "%s to go")
	def(&o.Done, "Goal reached")
	t := c.Theme()
	met := steps >= goal
	ink := t.Accent
	if met {
		ink = t.Success
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			kit.Arrive(tile(c, "footprints", ink, t.Space(11)), ui.Pose{Scale: 0.6, Rotate: -15, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
					kit.Numeric(c, number(c, steps)).FontSize(t.Rem(1.6)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine().
						ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, fmt.Sprintf(o.Of, number(c, goal))).FontSize(t.Rem(0.78)).Grow(1).MinWidth(0).SingleLine()
				})
			})
		})
		meter(c, float32(min(1, float64(steps)/float64(goal))), ink).Height(t.Space(2))
		c.Key(met)
		if met {
			kit.Arrive(kit.Chip(c, o.Done, t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
		} else {
			kit.Muted(c, fmt.Sprintf(o.Left, number(c, goal-steps))).FontSize(t.Rem(0.78)).FillWidth().MinWidth(0).SingleLine()
		}
		if len(o.Week) == 0 {
			return
		}
		week := ui.Row(c).FillWidth().Height(t.Space(24)).Gap(t.Space(1.5)).AlignItems(ui.End)
		week.Children(func() {
			for i, n := range o.Week {
				share := float32(float64(n) / float64(top))
				col := ui.Column(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Gap(t.Space(1)).AlignItems(ui.Center).Justify(ui.End)
				col.Children(func() {
					bar := ui.Box(c).FillWidth().Grow(1).MinHeight(0)
					grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
						ui.KeyframeOptions{Delay: time.Duration(i) * 50 * time.Millisecond})
					day := t.Accent.Alpha(0.35)
					if n >= goal {
						day = t.Success
					}
					bar.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						h := r.H * share * grow
						p.Fill(ui.Rect{X: r.X + r.W*0.2, Y: r.Y + r.H - h, W: r.W * 0.6, H: h}, day, t.Space(1))
					}))
					kit.Muted(c, o.Days[i]).FontSize(t.Rem(0.62)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
			}
		})
	})
	return root
}
