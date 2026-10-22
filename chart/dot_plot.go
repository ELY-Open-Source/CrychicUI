package chart

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// DotPlotOptions tune a DotPlot.
type DotPlotOptions = LineOptions

// DotPlot sets a dot per category, first on top.
func DotPlot(c *ui.Context, label string, labels []string, vs []float64, at *int, o DotPlotOptions) ui.Element {
	values("DotPlot", label, labels, vs)
	chosen("DotPlot", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	up := slices.Clone(labels)
	slices.Reverse(up)
	y := mc.NewBand(up...)
	lo, hi := mc.Extent(vs...)
	if hi == lo {
		hi = lo + 1
	}
	x := mc.NewLinear(lo, hi).Nice(4)
	ticks := x.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: y.Ticks(), Bottom: ticks}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, tick := range ticks {
					p.Line(mc.X(r, tick.At), r.Y, mc.X(r, tick.At), r.Y+r.H, 1, t.Border)
				}
				for i := range labels {
					yy := mc.Y(r, y.At(i))
					p.Line(r.X, yy, r.X+r.W, yy, 1, t.Border.Alpha(0.6))
				}
			}))
			for i, v := range vs {
				c.Key(i)
				tip := labels[i] + ": " + fmt.Sprintf(o.Format, v)
				dot(c, root, at, i, x.At(x.Lo+(v-x.Lo)*g), y.At(len(vs)-1-i), t.Space(3.5), tone(t, 0), tip, "DotPlot")
			}
		})
	})
}
