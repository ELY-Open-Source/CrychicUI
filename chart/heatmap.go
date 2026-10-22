package chart

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// HeatmapOptions tune a Heatmap.
type HeatmapOptions = LineOptions

// Heatmap shades cells by value; at counts row-major.
func Heatmap(c *ui.Context, label string, rows, cols []string, vs [][]float64, at *int, o HeatmapOptions) ui.Element {
	named("Heatmap", label, rows)
	named("Heatmap", label, cols)
	if len(vs) != len(rows) {
		panic(fmt.Sprintf("chart: Heatmap %q needs a value row per row", label))
	}
	var all []float64
	for _, r := range vs {
		if len(r) != len(cols) {
			panic(fmt.Sprintf("chart: Heatmap %q rows need a value per column", label))
		}
		all = append(all, r...)
	}
	finite("Heatmap", label, all...)
	chosen("Heatmap", label, at, len(rows)*len(cols))
	def(&o.Format, "%g")
	t := c.Theme()
	lo, hi := mc.Extent(all...)
	up := slices.Clone(rows)
	slices.Reverse(up)
	x := mc.NewBand(cols...).Padding(0.06, 0)
	y := mc.NewBand(up...).Padding(0.06, 0)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: y.Ticks(), Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			for r, row := range vs {
				for col, v := range row {
					i := r*len(cols) + col
					c.Key(i)
					k := float32(0.5)
					if hi > lo {
						k = float32((v - lo) / (hi - lo))
					}
					tip := fmt.Sprintf("%s, %s: "+o.Format, rows[r], cols[col], v)
					mark(c, root, at, i, x.Band(col), y.Band(len(rows)-1-r), t.Accent.Alpha((0.08+0.85*k)*g), tip, "Heatmap")
				}
			}
		})
	})
}
