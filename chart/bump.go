package chart

import (
	"fmt"
	"sort"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// BumpOptions tune a Bump.
type BumpOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Tip words a point: series, rank, label.
	Tip string
}

// Bump draws ranks over labels; at picks series.
func Bump(c *ui.Context, label string, labels []string, series []Series, at *int, o BumpOptions) ui.Element {
	vals, names := table("Bump", label, labels, series)
	chosen("Bump", label, at, len(series))
	def(&o.Tip, "%s: #%d in %s")
	t := c.Theme()
	n := len(series)
	ranks := make([][]float64, n)
	for s := range ranks {
		ranks[s] = make([]float64, len(labels))
	}
	for i := range labels {
		order := make([]int, n)
		for s := range order {
			order[s] = s
		}
		sort.SliceStable(order, func(a, b int) bool { return vals[order[a]][i] > vals[order[b]][i] })
		for r, s := range order {
			ranks[s][i] = float64(r + 1)
		}
	}
	place := func(rank float64) float32 {
		if n == 1 {
			return 0.5
		}
		return 1 - float32((rank-1)/float64(n-1))*0.9 - 0.05
	}
	var left []mc.Tick
	for r := 1; r <= n; r++ {
		left = append(left, mc.Tick{At: place(float64(r)), Label: fmt.Sprintf("#%d", r)})
	}
	x := mc.NewBand(labels...)
	xs := make([]float32, len(labels))
	for i := range labels {
		xs[i] = x.At(i)
	}
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: left, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			aims, ys := make([][]float32, n), make([][]float32, n)
			for s := range ranks {
				for _, r := range ranks[s] {
					aims[s] = append(aims[s], place(r))
				}
				ys[s] = eased(area, fmt.Sprint("r", s), aims[s])
			}
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for s := range ys {
					ink := tone(t, s)
					if *at >= 0 && *at != s {
						ink = ink.Alpha(0.25)
					}
					p.StrokePath(mc.Line(r, xs, ys[s], mc.CurveMonotone).Trim(0, g), 3, ink)
				}
			}))
			for s := range series {
				for i := range labels {
					c.Key(fmt.Sprint(s, ".", i))
					tip := fmt.Sprintf(o.Tip, names[s], int(ranks[s][i]), labels[i])
					dot(c, root, at, s, xs[i], aims[s][i], t.Space(3)*g, tone(t, s), tip, "Bump")
				}
			}
		})
		legend(c, names)
	})
}
