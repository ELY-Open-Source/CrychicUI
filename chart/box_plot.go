package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Group is a named sample of values.
type Group struct {
	Name   string
	Values []float64
}

// BoxPlotOptions tune a BoxPlot.
type BoxPlotOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Tip words a group: name, median, lower, upper quartile.
	Tip string
}

// BoxPlot draws quartile boxes and whiskers per group.
func BoxPlot(c *ui.Context, label string, groups []Group, at *int, o BoxPlotOptions) ui.Element {
	names, all := sample("BoxPlot", label, groups)
	chosen("BoxPlot", label, at, len(groups))
	def(&o.Format, "%g")
	def(&o.Tip, "%s: median %s, %s–%s")
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	t := c.Theme()
	qs := make([][5]float64, len(groups))
	for i, g := range groups {
		qs[i] = quartiles(g.Values)
	}
	x := mc.NewBand(names...).Padding(0.45, 0.2)
	y := mc.NewLinear(mc.Extent(all...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			shown := make([][]float32, len(qs))
			for i, q := range qs {
				aim := make([]float32, len(q))
				for k, v := range q {
					aim[k] = y.At(v)
				}
				shown[i] = eased(area, fmt.Sprint("q", i), aim)
			}
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for i, q := range shown {
					mid, b := mc.X(r, x.At(i)), x.Band(i)
					w := (mc.X(r, b.Hi) - mc.X(r, b.Lo)) / 4
					for _, end := range []float32{q[0], q[4]} {
						yv := mc.Y(r, q[2]+(end-q[2])*g)
						p.Line(mid, mc.Y(r, q[2]), mid, yv, 1, t.TextMuted)
						p.Line(mid-w, yv, mid+w, yv, 1, t.TextMuted)
					}
				}
			}))
			for i, q := range qs {
				c.Key(i)
				span := y.Span(q[2]-(q[2]-q[1])*float64(g), q[2]+(q[3]-q[2])*float64(g))
				tip := fmt.Sprintf(o.Tip, names[i], f(q[2]), f(q[1]), f(q[3]))
				mark(c, root, at, i, x.Band(i), span, tone(t, i), tip, "BoxPlot").
					DrawOver(func(p *ui.Painter, r ui.Rect) {
						m := r.Y + r.H*float32((q[3]-q[2])/max(q[3]-q[1], 1e-9))
						p.Line(r.X, m, r.X+r.W, m, 2, t.Background)
					})
			}
		})
	})
}
