package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// WaterfallOptions tune a Waterfall.
type WaterfallOptions struct {
	Height float32
	// Format prints a value; Total names the last bar.
	Format, Total string
}

// Waterfall steps up and down to a total.
func Waterfall(c *ui.Context, label string, labels []string, deltas []float64, at *int, o WaterfallOptions) ui.Element {
	values("Waterfall", label, labels, deltas)
	def(&o.Total, "Total")
	def(&o.Format, "%+g")
	all := append(append([]string(nil), labels...), o.Total)
	named("Waterfall", label, all)
	chosen("Waterfall", label, at, len(all))
	t := c.Theme()
	ends := []float64{0}
	run := 0.0
	for _, d := range deltas {
		run += d
		ends = append(ends, run)
	}
	x := mc.NewBand(all...).Padding(0.3, 0.15)
	y := mc.NewLinear(extent(ends...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(grid(t, ticks))
			for i, d := range deltas {
				c.Key(i)
				ink := t.Success
				if d < 0 {
					ink = t.Danger
				}
				from := ends[i]
				tip := labels[i] + ": " + fmt.Sprintf(o.Format, d)
				mark(c, root, at, i, x.Band(i), y.Span(from, from+d*g), ink, tip, "Waterfall")
			}
			c.Key("total")
			mark(c, root, at, len(deltas), x.Band(len(deltas)), y.Span(0, run*g), t.Accent, fmt.Sprintf("%s: "+o.Format, o.Total, run), "Waterfall")
		})
	})
}
