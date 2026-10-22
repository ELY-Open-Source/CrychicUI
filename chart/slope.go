package chart

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// SlopeOptions tune a Slope.
type SlopeOptions = LineOptions

// Slope joins each name's values across two columns.
func Slope(c *ui.Context, label string, names []string, from, to Series, at *int, o SlopeOptions) ui.Element {
	vals, cols := table("Slope", label, names, []Series{from, to})
	chosen("Slope", label, at, len(names))
	def(&o.Format, "%g")
	t := c.Theme()
	lo, hi := mc.Extent(append(slices.Clone(vals[0]), vals[1]...)...)
	if hi == lo {
		hi = lo + 1
	}
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	xs := [2]float32{0.15, 0.85}
	bottom := []mc.Tick{{At: xs[0], Label: cols[0]}, {At: xs[1], Label: cols[1]}}
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: bottom}, o.Height, func(area ui.Element) {
			g := grow(area)
			aims := [2][]float32{}
			for s := range aims {
				for _, v := range vals[s] {
					aims[s] = append(aims[s], y.At(v))
				}
			}
			a, b := eased(area, "a", aims[0]), eased(area, "b", aims[1])
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for _, x := range xs {
					p.Line(mc.X(r, x), r.Y, mc.X(r, x), r.Y+r.H, 1, t.Border)
				}
				for i := range a {
					ink := t.TextMuted
					switch {
					case vals[1][i] > vals[0][i]:
						ink = t.Success
					case vals[1][i] < vals[0][i]:
						ink = t.Danger
					}
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.3)
					}
					var line ui.Path
					line.MoveTo(mc.X(r, xs[0]), mc.Y(r, a[i])).LineTo(mc.X(r, xs[1]), mc.Y(r, b[i]))
					p.StrokePath(line.Trim(0, g), 2, ink)
				}
			}))
			for i, n := range names {
				tip := fmt.Sprintf("%s: %s %s, %s %s", n, cols[0], f(vals[0][i]), cols[1], f(vals[1][i]))
				for s := range 2 {
					c.Key(fmt.Sprint(s, ".", i))
					dot(c, root, at, i, xs[s], aims[s][i], t.Space(3), tone(t, 0), tip, "Slope")
				}
				c.Key(fmt.Sprint("name.", i))
				mc.Anchor(c, xs[1], aims[1][i], ui.Start, ui.Center, func() {
					kit.Muted(c, n).NoWrap().FontSize(t.Rem(0.75)).Margin(0, 0, 0, t.Space(2.5)).Opacity(g)
				}).Transition(spring)
			}
		})
	})
}
