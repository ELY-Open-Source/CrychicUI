package project

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Sprint is one sprint's completed points.
type Sprint struct {
	Name   string
	Points int
}

// VelocityOptions word the chart; Average takes the mean.
type VelocityOptions struct{ Title, Average string }

// Velocity shows points per sprint bars against their average.
func Velocity(c *ui.Context, sprints []Sprint, o VelocityOptions) ui.Element {
	if len(sprints) == 0 {
		panic("project: Velocity needs sprints")
	}
	def(&o.Title, "Velocity")
	def(&o.Average, "Avg %.0f pts")
	maxPoints, sum := 1, 0.0
	for _, sprint := range sprints {
		if sprint.Name == "" || sprint.Points < 0 {
			panic(fmt.Sprintf("project: Velocity %q has an invalid sprint %q", o.Title, sprint.Name))
		}
		maxPoints = max(maxPoints, sprint.Points)
		sum += float64(sprint.Points)
	}
	mean := sum / float64(len(sprints))
	t := c.Theme()
	barH := t.Space(30)
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "chart-column", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Average, mean), t.Accent) })
		chart := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Padding(t.Space(3))
		chart.DrawOver(func(p *ui.Painter, r ui.Rect) {
			y := r.Y + t.Space(3) + t.Space(5) + barH*(1-float32(mean)/float32(maxPoints))
			for x := r.X + t.Space(3); x < r.X+r.W-t.Space(3); x += t.Space(2) {
				p.Line(x, y, min(x+t.Space(1), r.X+r.W-t.Space(3)), y, 1, t.TextMuted.Alpha(0.6))
			}
		})
		chart.Children(func() {
			for i, sprint := range sprints {
				latest := i == len(sprints)-1
				ink := t.Accent.Alpha(0.45)
				if latest {
					ink = t.Accent
				}
				c.Key(i)
				ui.Column(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center).Gap(t.Space(1)).Children(func() {
					kit.Numeric(c, fmt.Sprint(sprint.Points)).FontSize(t.Rem(0.7)).Bold().Height(t.Space(4))
					well := ui.Box(c).FillWidth().MaxWidth(t.Space(10)).Height(barH)
					delay := time.Duration(i) * 50 * time.Millisecond
					frames := []ui.Keyframe{{Value: 0}}
					if delay > 0 {
						frames = append(frames, ui.Keyframe{At: delay})
					}
					rise := well.Keyframes("rise", append(frames, ui.Keyframe{At: delay + 600*time.Millisecond, Value: 1, Ease: ui.EaseOut}))
					share := well.AnimateTo("share", float32(sprint.Points)/float32(maxPoints), ui.Smooth)
					well.Draw(func(p *ui.Painter, r ui.Rect) {
						h := r.H * share * rise
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, ink, t.Space(1))
					})
					kit.Muted(c, sprint.Name).FontSize(t.Rem(0.65)).SingleLine().MinWidth(0)
				})
			}
		})
	})
	return root
}
