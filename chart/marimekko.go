package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// MarimekkoOptions tune a Marimekko.
type MarimekkoOptions = LineOptions

// Marimekko stacks shares in columns as wide as their totals.
func Marimekko(c *ui.Context, label string, cols []string, series []Series, at *int, o MarimekkoOptions) ui.Element {
	vals, names := table("Marimekko", label, cols, series)
	totals := make([]float64, len(cols))
	sum := 0.0
	for _, vs := range vals {
		for i, v := range vs {
			if v < 0 {
				panic(fmt.Sprintf("chart: Marimekko %q values must not be negative", label))
			}
			totals[i] += v
			sum += v
		}
	}
	for _, tot := range totals {
		if tot == 0 {
			panic(fmt.Sprintf("chart: Marimekko %q columns need a value above zero", label))
		}
	}
	chosen("Marimekko", label, at, len(cols))
	def(&o.Format, "%g")
	t := c.Theme()
	y := mc.NewLinear(0, 100)
	left := y.Ticks(4)
	for i := range left {
		left[i].Label += "%"
	}
	const gap = 0.006
	spans := make([]mc.Span, len(cols))
	var bottom []mc.Tick
	from := float32(0)
	for i, tot := range totals {
		w := float32(tot / sum)
		pad := min(gap, w/4)
		spans[i] = mc.Span{Lo: from + pad, Hi: from + w - pad}
		bottom = append(bottom, mc.Tick{At: from + w/2, Label: cols[i]})
		from += w
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: left, Bottom: bottom}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			for i, tot := range totals {
				base := 0.0
				for s, vs := range vals {
					c.Key(fmt.Sprint(i, ".", s))
					share := 100 * vs[i] / tot
					tip := fmt.Sprintf("%s, %s: "+o.Format+" (%.0f%%)", cols[i], names[s], vs[i], share)
					mark(c, root, at, i, spans[i], y.Span(base*g, (base+share)*g), tone(t, s), tip, "Marimekko")
					base += share
				}
			}
		})
		legend(c, names)
	})
}
