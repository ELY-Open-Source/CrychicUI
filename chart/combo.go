package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ComboOptions tune a Combo.
type ComboOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints bars, Line the line; "%g" default.
	Format, Line string
}

// Combo draws bars left, a line right.
func Combo(c *ui.Context, label string, labels []string, bars, line Series, at *int, o ComboOptions) ui.Element {
	values("Combo", label, labels, bars.Values)
	values("Combo", label, labels, line.Values)
	named("Combo", label, []string{bars.Name, line.Name})
	chosen("Combo", label, at, len(labels))
	def(&o.Format, "%g")
	def(&o.Line, "%g")
	t := c.Theme()
	x := mc.NewBand(labels...).Padding(0.3, 0.15)
	yb := mc.NewLinear(extent(bars.Values...)).Nice(4)
	yl := mc.NewLinear(extent(line.Values...)).Nice(4)
	ticks := yb.Ticks(4)
	xs := make([]float32, len(labels))
	for i := range labels {
		xs[i] = x.At(i)
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Right: yl.Ticks(4), Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			area.Draw(grid(t, ticks))
			aim := make([]float32, len(labels))
			for i, v := range bars.Values {
				aim[i] = yl.At(line.Values[i])
				tip := fmt.Sprintf("%s: %s "+o.Format+", %s "+o.Line, labels[i], bars.Name, v, line.Name, line.Values[i])
				c.Key(fmt.Sprint("bar.", i))
				mark(c, root, at, i, x.Band(i), yb.Span(0, v*float64(g)), tone(t, 0), tip, "Combo")
				c.Key(fmt.Sprint("dot.", i))
				dot(c, root, at, i, xs[i], aim[i], t.Space(3)*g, tone(t, 1), tip, "Combo")
			}
			ys := eased(area, "line", aim)
			area.DrawOver(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				p.StrokePath(mc.Line(r, xs, ys, mc.CurveMonotone).Trim(0, g), 2, tone(t, 1))
			}))
		})
		legend(c, []string{bars.Name, line.Name})
	})
}
