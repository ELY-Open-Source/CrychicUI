package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// RingsOptions tune a Rings; Format takes a percent.
type RingsOptions = PieOptions

// Rings draws progress rings; vs run 0 to 1.
func Rings(c *ui.Context, label string, names []string, vs []float64, at *int, o RingsOptions) ui.Element {
	values("Rings", label, names, vs)
	for _, v := range vs {
		if v < 0 || v > 1 {
			panic(fmt.Sprintf("chart: Rings %q values must run 0 to 1", label))
		}
	}
	chosen("Rings", label, at, len(names))
	def(&o.Format, "%.0f%%")
	t := c.Theme()
	size := side("Rings", label, o.Size, t.Space(56))
	outer := size/2 - t.Space(1)
	width := outer * 0.7 / float32(len(vs))
	ring := func(i int) (float32, float32) {
		hi := outer - width*float32(i)
		return hi - width*0.78, hi
	}
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		shown := springs(root, "v", vs)
		arc := func(i int, g float32) float64 { return 2 * math.Pi * shown[i] * float64(g) }
		wheel{label: label, comp: "Rings", n: len(vs), size: size,
			hit: func(x, y, g float32) int {
				d := radius(x, y)
				for i := range vs {
					if lo, hi := ring(i); d >= lo && d <= hi && turn(x, y) <= arc(i, g) {
						return i
					}
				}
				return -1
			},
			describe: func(i int) string { return names[i] + ": " + fmt.Sprintf(o.Format, 100*vs[i]) },
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				for i := range vs {
					lo, hi := ring(i)
					ink := tone(t, i)
					p.FillPath(mc.Sector(cx, cy, lo, hi, mc.Slice{End: 2 * math.Pi}), ink.Alpha(0.15))
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.4)
					}
					if i == hot {
						ink = ink.Mix(t.Text, 0.15)
					}
					p.FillPath(mc.Sector(cx, cy, lo, hi, mc.Slice{End: arc(i, g)}), ink)
				}
			},
		}.build(c, root, at, nil)
		legend(c, names)
	})
}
