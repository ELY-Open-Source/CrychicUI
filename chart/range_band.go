package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// RangeBandOptions tune a RangeBand.
type RangeBandOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Tip words a point: label, value, low, high.
	Tip string
}

// RangeBand draws a line inside a lo–hi band.
func RangeBand(c *ui.Context, label string, labels []string, mid, lo, hi []float64, at *int, o RangeBandOptions) ui.Element {
	values("RangeBand", label, labels, mid)
	if len(lo) != len(labels) || len(hi) != len(labels) {
		panic(fmt.Sprintf("chart: RangeBand %q needs a low and high per label", label))
	}
	finite("RangeBand", label, lo...)
	finite("RangeBand", label, hi...)
	for i := range mid {
		if lo[i] > mid[i] || mid[i] > hi[i] {
			panic(fmt.Sprintf("chart: RangeBand %q bands must hold the line", label))
		}
	}
	chosen("RangeBand", label, at, len(labels))
	def(&o.Format, "%g")
	def(&o.Tip, "%s: %s (%s–%s)")
	t := c.Theme()
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	bottom, top := mc.Extent(append(append([]float64{}, lo...), hi...)...)
	if bottom == top {
		top = bottom + 1
	}
	x := mc.NewBand(labels...)
	y := mc.NewLinear(bottom, top).Nice(4)
	ticks := y.Ticks(4)
	xs := make([]float32, len(labels))
	for i := range labels {
		xs[i] = x.At(i)
	}
	at32 := func(vs []float64) []float32 {
		out := make([]float32, len(vs))
		for i, v := range vs {
			out[i] = y.At(v)
		}
		return out
	}
	describe := func(i int) string { return fmt.Sprintf(o.Tip, labels[i], f(mid[i]), f(lo[i]), f(hi[i])) }
	root := ui.Column(c).FillWidth()
	hover := ui.Local(root, "hover", func() int { return -1 })
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			aim := at32(mid)
			ym, yl, yh := eased(area, "m", aim), eased(area, "l", at32(lo)), eased(area, "h", at32(hi))
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				p.FillPath(mc.Area(r, xs, yl, yh, mc.CurveMonotone), t.Accent.Alpha(0.18*g))
				p.StrokePath(mc.Line(r, xs, ym, mc.CurveMonotone).Trim(0, g), 2, t.Accent)
			}))
			if i := *at; i >= 0 {
				mc.Point(c, xs[i], aim[i], t.Space(3), t.Space(3)).Radius(999).Background(t.Background).
					Border(2, t.Accent).Transition(spring)
			}
			follow(c, root, area, at, hover, xs, 1, describe, "RangeBand")
		})
	})
}
