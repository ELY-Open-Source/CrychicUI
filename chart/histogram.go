package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// HistogramOptions tune a Histogram.
type HistogramOptions struct {
	// Bins is the bar count; zero gives 10.
	Bins   int
	Height float32
	// Format prints a bin edge; "%g" when empty.
	Format string
}

// Histogram counts values into even bins.
func Histogram(c *ui.Context, label string, vs []float64, at *int, o HistogramOptions) ui.Element {
	if label == "" || len(vs) == 0 || o.Bins < 0 {
		panic(fmt.Sprintf("chart: Histogram %q needs values and bins", label))
	}
	finite("Histogram", label, vs...)
	if o.Bins == 0 {
		o.Bins = 10
	}
	chosen("Histogram", label, at, o.Bins)
	def(&o.Format, "%g")
	t := c.Theme()
	lo, hi := mc.Extent(vs...)
	if hi == lo {
		hi = lo + 1
	}
	step := (hi - lo) / float64(o.Bins)
	counts := make([]float64, o.Bins)
	for _, v := range vs {
		counts[min(int((v-lo)/step), o.Bins-1)]++
	}
	x := mc.NewLinear(lo, hi)
	y := mc.NewLinear(extent(counts...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks(5)}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(grid(t, ticks))
			for i, n := range counts {
				c.Key(i)
				a, b := lo+step*float64(i), lo+step*float64(i+1)
				span := x.Span(a, b)
				span.Lo, span.Hi = span.Lo+0.002, span.Hi-0.002
				tip := fmt.Sprintf(o.Format+"–"+o.Format+": %d", a, b, int(n))
				mark(c, root, at, i, span, y.Span(0, n*g), t.Accent, tip, "Histogram")
			}
		})
	})
}
