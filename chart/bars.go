package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// BarsOptions tune a bar chart.
type BarsOptions = LineOptions

// Bars draws a vertical bar per label.
func Bars(c *ui.Context, label string, labels []string, vs []float64, at *int, o BarsOptions) ui.Element {
	values("Bars", label, labels, vs)
	chosen("Bars", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	x := mc.NewBand(labels...).Padding(0.3, 0.15)
	y := mc.NewLinear(extent(vs...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(grid(t, ticks))
			for i, v := range vs {
				c.Key(i)
				mark(c, root, at, i, x.Band(i), y.Span(0, v*g), t.Accent, labels[i]+": "+fmt.Sprintf(o.Format, v), "Bars")
			}
		})
	})
}
