package dashboard

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Anomaly shows a chart with pulsing flagged points.
func Anomaly(c *ui.Context, points []Point, flags []int) ui.Element {
	if len(points) < 2 || len(flags) == 0 {
		panic("dashboard: Anomaly needs points and flags")
	}
	for _, index := range flags {
		if index < 0 || index >= len(points) {
			panic("dashboard: Anomaly flag out of range")
		}
	}
	lo, hi := points[0].Value, points[0].Value
	for _, point := range points {
		if !finite(point.Value) {
			panic("dashboard: Anomaly value must be finite")
		}
		lo, hi = min(lo, point.Value), max(hi, point.Value)
	}
	span := hi - lo
	if !finite(span) {
		panic("dashboard: Anomaly range overflow")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	phase := root.Loop("pulse", time.Second, ui.Linear)
	root.Children(func() {
		box := ui.Box(c).FillWidth().Height(t.Space(16))
		box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			for i, point := range points {
				x := r.X + r.W*float32(i)/float32(len(points)-1)
				y := r.Y + r.H*0.5
				if span > 0 {
					y = r.Y + r.H*(1-float32((point.Value-lo)/span))
				}
				if i == 0 {
					path.MoveTo(x, y)
				} else {
					path.LineTo(x, y)
				}
			}
			p.StrokePath(&path, 2, t.Accent)
			for _, index := range flags {
				x := r.X + r.W*float32(index)/float32(len(points)-1)
				y := r.Y + r.H*0.5
				if span > 0 {
					y = r.Y + r.H*(1-float32((points[index].Value-lo)/span))
				}
				var dot ui.Path
				dot.Circle(x, y, t.Space(1))
				p.FillPath(&dot, t.Danger.Alpha(0.4+0.6*unit(phase < 0.5)))
			}
		}))
	})
	return root
}
