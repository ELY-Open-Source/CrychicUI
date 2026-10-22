package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// GroupedBarsOptions tune a GroupedBars.
type GroupedBarsOptions = LineOptions

// GroupedBars sets each series' bars side by side per label.
func GroupedBars(c *ui.Context, label string, labels []string, series []Series, at *int, o GroupedBarsOptions) ui.Element {
	vals, names := table("GroupedBars", label, labels, series)
	chosen("GroupedBars", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	var all []float64
	for _, vs := range vals {
		all = append(all, vs...)
	}
	x := mc.NewBand(labels...).Padding(0.25, 0.1)
	y := mc.NewLinear(extent(all...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(grid(t, ticks))
			for i := range labels {
				band := x.Band(i)
				w := (band.Hi - band.Lo) / float32(len(vals))
				for s, vs := range vals {
					c.Key(fmt.Sprint(i, "/", s))
					span := mc.Span{Lo: band.Lo + w*float32(s), Hi: band.Lo + w*float32(s+1)}
					tip := fmt.Sprintf("%s, %s: "+o.Format, labels[i], names[s], vs[i])
					mark(c, root, at, i, span, y.Span(0, vs[i]*g), tone(t, s), tip, "GroupedBars")
				}
			}
		})
		legend(c, names)
	})
}
