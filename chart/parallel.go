package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ParallelOptions tune a Parallel.
type ParallelOptions = LineOptions

// Parallel draws rows across their own scaled axes.
func Parallel(c *ui.Context, label string, axes []string, rows []Series, at *int, o ParallelOptions) ui.Element {
	vals, names := table("Parallel", label, axes, rows)
	chosen("Parallel", label, at, len(rows))
	def(&o.Format, "%g")
	t := c.Theme()
	scales := make([]mc.Linear, len(axes))
	for a := range axes {
		col := make([]float64, len(rows))
		for s := range rows {
			col[s] = vals[s][a]
		}
		lo, hi := mc.Extent(col...)
		if hi == lo {
			lo, hi = lo-1, hi+1
		}
		scales[a] = mc.NewLinear(lo, hi)
	}
	// Inset keeps ends clear of the axis labels.
	place := func(a int, v float64) float32 { return 0.08 + 0.84*scales[a].At(v) }
	x := mc.NewBand(axes...)
	xs := make([]float32, len(axes))
	for a := range axes {
		xs[a] = x.At(a)
	}
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		plot(c, label, mc.Axes{Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			aims, ys := make([][]float32, len(rows)), make([][]float32, len(rows))
			for s := range rows {
				for a, v := range vals[s] {
					aims[s] = append(aims[s], place(a, v))
				}
				ys[s] = eased(area, fmt.Sprint("v", s), aims[s])
			}
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, xx := range xs {
					p.Line(mc.X(r, xx), r.Y, mc.X(r, xx), r.Y+r.H, 1, t.Border)
				}
				for s := range ys {
					ink := tone(t, s)
					if *at >= 0 && *at != s {
						ink = ink.Alpha(0.2)
					}
					p.StrokePath(mc.Line(r, xs, ys[s], mc.CurveLinear).Trim(0, g), 2, ink)
				}
			}))
			for a := range axes {
				c.Key(fmt.Sprint("max.", a))
				mc.Anchor(c, xs[a], 1, ui.Center, ui.Start, func() {
					kit.Muted(c, f(scales[a].Hi)).NoWrap().FontSize(t.Rem(0.65))
				})
			}
			for s := range rows {
				for a := range axes {
					c.Key(fmt.Sprint(s, ".", a))
					tip := fmt.Sprintf("%s, %s: %s", names[s], axes[a], f(vals[s][a]))
					dot(c, root, at, s, xs[a], aims[s][a], t.Space(2.5)*g, tone(t, s), tip, "Parallel")
				}
			}
		})
		legend(c, names)
	})
}
