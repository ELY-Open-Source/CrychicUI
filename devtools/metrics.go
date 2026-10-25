package devtools

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Metric is one named series; Lower marks lower as better.
type Metric struct {
	Name, Unit string
	Points     []float64
	Lower      bool
}

// MetricsOptions format the change; Change takes a signed percent.
type MetricsOptions struct{ Change string }

// Metrics shows metric charts in a grid.
func Metrics(c *ui.Context, metrics []Metric, o MetricsOptions) ui.Element {
	if len(metrics) == 0 {
		panic("devtools: Metrics needs metrics")
	}
	names := make([]string, len(metrics))
	for i, m := range metrics {
		if len(m.Points) < 2 {
			panic(fmt.Sprintf("devtools: Metrics %q needs at least two points", m.Name))
		}
		bounds("Metrics "+m.Name, m.Points)
		if first := m.Points[0]; first != 0 && !finite((m.Points[len(m.Points)-1]-first)/math.Abs(first)*100) {
			panic(fmt.Sprintf("devtools: Metrics %q change overflows", m.Name))
		}
		names[i] = m.Name
	}
	unique("Metrics", names)
	def(&o.Change, "%+.4g%%")
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, m := range metrics {
			ink := inks[i%len(inks)]
			lo, hi := bounds("Metrics "+m.Name, m.Points)
			first, last := m.Points[0], m.Points[len(m.Points)-1]
			c.Key(m.Name)
			tileBox := ui.Column(c).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(3)).Radius(t.Radius*1.5).
				Border(1, t.Border).Background(t.Background)
			kit.Stagger(c, root, tileBox, i).Children(func() {
				kit.Muted(c, m.Name).FontSize(t.Rem(0.75)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).AlignItems(ui.End).Children(func() {
					kit.Numeric(c, strconv.FormatFloat(last, 'g', 4, 64)).FontSize(t.Rem(1.3)).Bold().Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, m.Unit).FontSize(t.Rem(0.75)).Padding(0, 0, t.Space(0.5), 0).Grow(1).MinWidth(0).SingleLine()
					if first != 0 {
						pct := (last - first) / math.Abs(first) * 100
						good, bad := t.Success, t.Danger
						if m.Lower {
							good, bad = bad, good
						}
						tone := t.TextMuted
						if pct > 0 {
							tone = good
						} else if pct < 0 {
							tone = bad
						}
						kit.Chip(c, fmt.Sprintf(o.Change, pct), tone)
					}
				})
				chart := ui.Box(c).FillWidth().Height(t.Space(12))
				draw := chart.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
				chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					line := spark(r, m.Points, lo, hi)
					area := spark(r, m.Points, lo, hi)
					area.LineTo(r.X+r.W, r.Y+r.H)
					area.LineTo(r.X, r.Y+r.H)
					p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * draw, H: r.H}, 0, func() {
						p.FillPath(area.Close(), ink.Alpha(0.12))
					})
					p.StrokePath(line.Trim(0, draw), t.Space(0.5), ink)
				}))
			})
		}
	})
	return root
}
