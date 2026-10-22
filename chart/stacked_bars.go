package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// StackedBarsOptions tune a StackedBars.
type StackedBarsOptions = LineOptions

// StackedBars piles each series' bar per label.
func StackedBars(c *ui.Context, label string, labels []string, series []Series, at *int, o StackedBarsOptions) ui.Element {
	vals, names := table("StackedBars", label, labels, series)
	chosen("StackedBars", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	segs := mc.Stack(vals, mc.StackStandard)
	var ends []float64
	for _, ss := range segs {
		for _, sg := range ss {
			ends = append(ends, sg.Lo, sg.Hi)
		}
	}
	x := mc.NewBand(labels...).Padding(0.3, 0.15)
	y := mc.NewLinear(extent(ends...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(grid(t, ticks))
			for i := range labels {
				for s, ss := range segs {
					c.Key(fmt.Sprint(i, "/", s))
					tip := fmt.Sprintf("%s, %s: "+o.Format, labels[i], names[s], vals[s][i])
					mark(c, root, at, i, x.Band(i), y.Span(ss[i].Lo*g, ss[i].Hi*g), tone(t, s), tip, "StackedBars")
				}
			}
		})
		legend(c, names)
	})
}
