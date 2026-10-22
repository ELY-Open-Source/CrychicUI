package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ErrorBarsOptions tune an ErrorBars.
type ErrorBarsOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Tip words a point: label, value, error.
	Tip string
}

// ErrorBars draws points with whiskers of plus or minus errs.
func ErrorBars(c *ui.Context, label string, labels []string, vs, errs []float64, at *int, o ErrorBarsOptions) ui.Element {
	values("ErrorBars", label, labels, vs)
	if len(errs) != len(labels) {
		panic(fmt.Sprintf("chart: ErrorBars %q needs an error per label", label))
	}
	finite("ErrorBars", label, errs...)
	ends := make([]float64, 0, 2*len(vs))
	for i, e := range errs {
		if e < 0 {
			panic(fmt.Sprintf("chart: ErrorBars %q errors must not be negative", label))
		}
		ends = append(ends, vs[i]-e, vs[i]+e)
	}
	chosen("ErrorBars", label, at, len(labels))
	def(&o.Format, "%g")
	def(&o.Tip, "%s: %s ± %s")
	t := c.Theme()
	lo, hi := mc.Extent(ends...)
	if hi == lo {
		hi = lo + 1
	}
	x := mc.NewBand(labels...).Padding(0.3, 0.15)
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			var at0, hi0, lo0 []float32
			for i, v := range vs {
				at0, hi0, lo0 = append(at0, y.At(v)), append(hi0, y.At(v+errs[i])), append(lo0, y.At(v-errs[i]))
			}
			mids, his, los := eased(area, "v", at0), eased(area, "hi", hi0), eased(area, "lo", lo0)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for i, m := range mids {
					mid, w := mc.X(r, x.At(i)), (mc.X(r, x.Band(i).Hi)-mc.X(r, x.Band(i).Lo))/4
					top, bot := mc.Y(r, m+(his[i]-m)*g), mc.Y(r, m+(los[i]-m)*g)
					p.Line(mid, top, mid, bot, 1, t.TextMuted)
					p.Line(mid-w, top, mid+w, top, 1, t.TextMuted)
					p.Line(mid-w, bot, mid+w, bot, 1, t.TextMuted)
				}
			}))
			for i, v := range vs {
				c.Key(i)
				tip := fmt.Sprintf(o.Tip, labels[i], fmt.Sprintf(o.Format, v), fmt.Sprintf(o.Format, errs[i]))
				dot(c, root, at, i, x.At(i), at0[i], t.Space(3)*g, tone(t, 0), tip, "ErrorBars")
			}
		})
	})
}
