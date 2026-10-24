package dashboard

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AnomalyOptions word the card; Flagged takes the count.
type AnomalyOptions struct{ Title, Flagged string }

// Anomaly charts a series with flagged points pulsing and listed.
func Anomaly(c *ui.Context, points []Point, flags []int, o AnomalyOptions) ui.Element {
	if len(points) < 2 || len(flags) == 0 {
		panic("dashboard: Anomaly needs two points and a flag")
	}
	for i, index := range flags {
		if index < 0 || index >= len(points) || slices.Index(flags, index) != i {
			panic(fmt.Sprintf("dashboard: Anomaly flag %d must be a unique point index", index))
		}
	}
	lo, hi := points[0].Value, points[0].Value
	for _, point := range points {
		if !finite(point.Value) {
			panic(fmt.Sprintf("dashboard: Anomaly value at %q must be finite", point.Label))
		}
		lo, hi = min(lo, point.Value), max(hi, point.Value)
	}
	span := hi - lo
	if !finite(span) {
		panic("dashboard: Anomaly range overflows")
	}
	def(&o.Title, "Anomalies")
	def(&o.Flagged, "%d flagged")
	t := c.Theme()
	root := card(c)
	phase := root.Loop("pulse", 1400*time.Millisecond, ui.Linear)
	root.Children(func() {
		kit.CardHead(c, "radar", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Flagged, len(flags)), t.Danger) })
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			chart := ui.Box(c).FillWidth().Height(t.Space(28))
			draw := chart.Keyframes("crychic.draw", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				pad := t.Space(3)
				at := func(i int) (float32, float32) {
					x := r.X + pad + (r.W-2*pad)*float32(i)/float32(len(points)-1)
					y := r.Y + r.H/2
					if span > 0 {
						y = r.Y + pad + (r.H-2*pad)*(1-float32((points[i].Value-lo)/span))
					}
					return x, y
				}
				var line, area ui.Path
				x0, _ := at(0)
				xn, _ := at(len(points) - 1)
				area.MoveTo(x0, r.Y+r.H)
				for i := range points {
					x, y := at(i)
					if i == 0 {
						line.MoveTo(x, y)
					} else {
						line.LineTo(x, y)
					}
					area.LineTo(x, y)
				}
				area.LineTo(xn, r.Y+r.H).Close()
				p.FillPath(&area, t.Accent.Alpha(0.1*draw))
				p.StrokePath(line.Trim(0, draw), t.Space(0.5), t.Accent)
				for _, index := range flags {
					x, y := at(index)
					halo := t.Space(1.5) + t.Space(2)*phase
					p.Fill(ui.Rect{X: x - halo, Y: y - halo, W: 2 * halo, H: 2 * halo}, t.Danger.Alpha(0.3*(1-phase)*draw), 999)
					d := t.Space(2.5)
					p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Danger.Alpha(draw), 999)
				}
			}))
			list := ui.Column(c).FillWidth().Gap(t.Space(0.5))
			list.Children(func() {
				for i, index := range flags {
					c.Key(index)
					point := points[index]
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
						Background(t.Danger.Alpha(0.06))
					kit.Stagger(c, list, row, i).Children(func() {
						ui.Icon(c, icon.Must("triangle-alert")).Size(14, 14).Shrink(0).TextColor(t.Danger)
						ui.Text(c, point.Label).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, number(c, point.Value)).FontSize(t.Rem(0.8)).Shrink(0)
					})
				}
			})
		})
	})
	return root
}
