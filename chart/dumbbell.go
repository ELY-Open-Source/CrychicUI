package chart

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// DumbbellOptions tune a Dumbbell.
type DumbbellOptions = LineOptions

// Dumbbell joins two dots per category, the first on top.
func Dumbbell(c *ui.Context, label string, labels []string, from, to Series, at *int, o DumbbellOptions) ui.Element {
	vals, names := table("Dumbbell", label, labels, []Series{from, to})
	chosen("Dumbbell", label, at, len(labels))
	def(&o.Format, "%g")
	t := c.Theme()
	up := slices.Clone(labels)
	slices.Reverse(up)
	y := mc.NewBand(up...)
	lo, hi := mc.Extent(append(slices.Clone(vals[0]), vals[1]...)...)
	if hi == lo {
		hi = lo + 1
	}
	x := mc.NewLinear(lo, hi).Nice(4)
	ticks := x.Ticks(4)
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: y.Ticks(), Bottom: ticks}, o.Height, func(area ui.Element) {
			g := grow(area)
			aims := [2][]float32{}
			for s := range aims {
				for _, v := range vals[s] {
					aims[s] = append(aims[s], x.At(v))
				}
			}
			a, b := eased(area, "a", aims[0]), eased(area, "b", aims[1])
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, tick := range ticks {
					p.Line(mc.X(r, tick.At), r.Y, mc.X(r, tick.At), r.Y+r.H, 1, t.Border)
				}
				for i := range a {
					yy := mc.Y(r, y.At(len(a)-1-i))
					x0, x1 := mc.X(r, a[i]), mc.X(r, b[i])
					p.Line(min(x0, x1), yy, min(x0, x1)+(max(x0, x1)-min(x0, x1))*g, yy, 3, t.Border)
				}
			}))
			for i := range labels {
				tip := fmt.Sprintf("%s: %s %s, %s %s", labels[i], names[0], f(vals[0][i]), names[1], f(vals[1][i]))
				row := y.At(len(labels) - 1 - i)
				for s := range 2 {
					c.Key(fmt.Sprint(s, ".", i))
					dot(c, root, at, i, aims[s][i], row, t.Space(3.5)*g, tone(t, s), tip, "Dumbbell")
				}
			}
		})
		legend(c, names)
	})
}
