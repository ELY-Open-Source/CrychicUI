package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// LollipopOptions tune a Lollipop.
type LollipopOptions = LineOptions

// Lollipop draws a stem up to a dot per label.
func Lollipop(c *ui.Context, label string, labels []string, vs []float64, at *int, o LollipopOptions) ui.Element {
	values("Lollipop", label, labels, vs)
	chosen("Lollipop", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	x := mc.NewBand(labels...)
	y := mc.NewLinear(extent(vs...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			base := y.At(0)
			aims := make([]float32, len(vs))
			for i, v := range vs {
				aims[i] = y.At(v)
			}
			shown := eased(area, "v", aims)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for i, f := range shown {
					mid := mc.X(r, x.At(i))
					p.Line(mid, mc.Y(r, base), mid, mc.Y(r, base+(f-base)*float32(g)), 2, tone(t, 0).Alpha(0.5))
				}
			}))
			for i, v := range vs {
				c.Key(i)
				tip := labels[i] + ": " + fmt.Sprintf(o.Format, v)
				dot(c, root, at, i, x.At(i), y.At(v*g), t.Space(3.5), tone(t, 0), tip, "Lollipop")
			}
		})
	})
}
