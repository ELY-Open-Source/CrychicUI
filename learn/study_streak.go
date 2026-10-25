package learn

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StudyDay is one day's study minutes.
type StudyDay struct {
	Label   string
	Minutes int
	Today   bool
}

// StudyStreakOptions word the week; Goal is daily minutes.
type StudyStreakOptions struct {
	Title, Streak, Week, Minutes string
	Goal                         int
}

// StudyStreak shows a week of study days against a goal.
func StudyStreak(c *ui.Context, days []StudyDay, streak int, o StudyStreakOptions) ui.Element {
	def(&o.Title, "Study streak")
	if len(days) == 0 || len(days) > 7 || streak < 0 || o.Goal < 0 {
		panic(fmt.Sprintf("learn: StudyStreak %q needs 1 to 7 days, a streak and a goal of at least 0", o.Title))
	}
	labels := make([]string, len(days))
	total, today := 0, 0
	for i, d := range days {
		if d.Minutes < 0 || d.Minutes > math.MaxInt-total {
			panic(fmt.Sprintf("learn: StudyStreak day %q has %d minutes", d.Label, d.Minutes))
		}
		total += d.Minutes
		labels[i] = d.Label
		if d.Today {
			today++
		}
	}
	unique("StudyStreak", labels)
	if today > 1 {
		panic(fmt.Sprintf("learn: StudyStreak %q has %d days marked today", o.Title, today))
	}
	def(&o.Streak, "%d-day streak")
	def(&o.Week, "%s this week")
	def(&o.Minutes, "%d min")
	if o.Goal == 0 {
		o.Goal = 20
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			flame := tile(c, "flame", t.Warning, t.Space(11))
			kit.Arrive(flame, ui.Pose{Scale: 0.5, Rotate: -20, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				kit.Numeric(c, fmt.Sprintf(o.Streak, streak)).FontSize(t.Rem(1.1)).Bold().FillWidth().MinWidth(0).SingleLine().
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
			kit.Chip(c, fmt.Sprintf(o.Week, fmt.Sprintf(o.Minutes, total)), t.Accent)
		})
		week := ui.Row(c).FillWidth().Gap(t.Space(1)).Justify(ui.SpaceBetween)
		week.Children(func() {
			size := t.Space(8)
			for i, d := range days {
				share := float32(min(1, float64(d.Minutes)/float64(o.Goal)))
				met := d.Minutes >= o.Goal
				c.Key(i)
				col := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1.5), 0).Radius(t.Radius)
				if d.Today {
					col.Background(t.Accent.Alpha(0.08))
				}
				kit.Stagger(c, week, col, i).Children(func() {
					lab := ui.Text(c, d.Label).FontSize(t.Rem(0.66)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					if d.Today {
						lab.TextColor(t.Accent)
					} else {
						lab.TextColor(t.TextMuted)
					}
					day := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(999)
					grow := day.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
						ui.KeyframeOptions{Delay: time.Duration(i) * 60 * time.Millisecond})
					arc := grow * share
					if met {
						day.Background(t.Warning.Alpha(0.2 * grow)).TextColor(t.Warning)
					}
					day.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						w := t.Space(0.75)
						cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
						p.StrokePath(ring(cx, cy, rr), w, t.SurfaceHover)
						if arc > 0 {
							p.StrokePath(ring(cx, cy, rr).Trim(0, arc), w, t.Warning)
						}
					}))
					day.Children(func() {
						if met {
							glyph(c, "flame", 16)
						}
					})
					kit.Muted(c, fmt.Sprintf(o.Minutes, d.Minutes)).FontSize(t.Rem(0.6)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
			}
		})
	})
	return root
}
