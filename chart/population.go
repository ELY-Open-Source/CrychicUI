package chart

import (
	"fmt"
	"math"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// PopulationOptions tune a Population.
type PopulationOptions = LineOptions

// Population sets two bars back to back per group.
func Population(c *ui.Context, label string, groups []string, left, right Series, at *int, o PopulationOptions) ui.Element {
	vals, names := table("Population", label, groups, []Series{left, right})
	top := 0.0
	for _, vs := range vals {
		for _, v := range vs {
			if v < 0 {
				panic(fmt.Sprintf("chart: Population %q values must not be negative", label))
			}
			top = max(top, v)
		}
	}
	chosen("Population", label, at, len(groups))
	def(&o.Format, "%g")
	t := c.Theme()
	x := mc.NewLinear(-max(top, 1), max(top, 1)).Nice(4)
	ticks := x.Ticks(4)
	for i := range ticks {
		ticks[i].Label = fmt.Sprintf(o.Format, math.Abs(ticks[i].Value))
	}
	up := slices.Clone(groups)
	slices.Reverse(up)
	y := mc.NewBand(up...).Padding(0.15, 0.05)
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: y.Ticks(), Bottom: ticks}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, tick := range ticks {
					p.Line(mc.X(r, tick.At), r.Y, mc.X(r, tick.At), r.Y+r.H, 1, t.Border)
				}
			}))
			for i, gr := range groups {
				tip := fmt.Sprintf("%s: %s %s, %s %s", gr, names[0], f(vals[0][i]), names[1], f(vals[1][i]))
				band := y.Band(len(groups) - 1 - i)
				c.Key(fmt.Sprint("l.", i))
				mark(c, root, at, i, x.Span(-vals[0][i]*g, 0), band, tone(t, 0), tip, "Population")
				c.Key(fmt.Sprint("r.", i))
				mark(c, root, at, i, x.Span(0, vals[1][i]*g), band, tone(t, 1), tip, "Population")
			}
		})
		legend(c, names)
	})
}
