package chart

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// HBarsOptions tune an HBars.
type HBarsOptions = LineOptions

// HBars draws horizontal bars, the first on top.
func HBars(c *ui.Context, label string, labels []string, vs []float64, at *int, o HBarsOptions) ui.Element {
	values("HBars", label, labels, vs)
	chosen("HBars", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	up := slices.Clone(labels)
	slices.Reverse(up)
	y := mc.NewBand(up...).Padding(0.3, 0.15)
	x := mc.NewLinear(extent(vs...)).Nice(4)
	ticks := x.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: y.Ticks(), Bottom: ticks}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, tick := range ticks {
					p.Line(mc.X(r, tick.At), r.Y, mc.X(r, tick.At), r.Y+r.H, 1, t.Border)
				}
			}))
			for i, v := range vs {
				c.Key(i)
				mark(c, root, at, i, x.Span(0, v*g), y.Band(len(vs)-1-i), t.Accent, labels[i]+": "+fmt.Sprintf(o.Format, v), "HBars")
			}
		})
	})
}
